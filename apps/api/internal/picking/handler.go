package picking

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

var errNotFound = errors.New("not found")

type PickingList struct {
	ID          string    `json:"id"`
	WarehouseID string    `json:"warehouse_id"`
	RouteID     string    `json:"route_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Items       []Item    `json:"items,omitempty"`
}

type Item struct {
	ID          string `json:"id"`
	ProductID   string `json:"product_id"`
	QtyRequired string `json:"qty_required"`
	QtyPicked   string `json:"qty_picked"`
	IsVerified  bool   `json:"is_verified"`
}

const plCols = `id, warehouse_id::text, COALESCE(route_id::text,''), status, created_at, updated_at`
const itemCols = `id, product_id::text, qty_required::text, qty_picked::text, is_verified`

type itemInput struct {
	ProductID   string `json:"product_id"`
	QtyRequired string `json:"qty_required"`
}

func HandleCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			WarehouseID string      `json:"warehouse_id"`
			RouteID     string      `json:"route_id"`
			Items       []itemInput `json:"items"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.WarehouseID == "" {
			http.Error(w, "warehouse_id is required", http.StatusBadRequest)
			return
		}
		if len(body.Items) == 0 {
			http.Error(w, "items required", http.StatusBadRequest)
			return
		}

		var pl PickingList
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`INSERT INTO picking_lists (tenant_id, warehouse_id, route_id, status)
				 VALUES ($1,$2,NULLIF($3,'')::uuid,'pending')
				 RETURNING `+plCols,
				mw.TenantIDFrom(r.Context()), body.WarehouseID, body.RouteID,
			).Scan(&pl.ID, &pl.WarehouseID, &pl.RouteID, &pl.Status, &pl.CreatedAt, &pl.UpdatedAt); err != nil {
				return err
			}
			for _, it := range body.Items {
				if it.QtyRequired == "" {
					it.QtyRequired = "0"
				}
				var item Item
				if err := tx.QueryRow(r.Context(),
					`INSERT INTO picking_items (picking_list_id, product_id, qty_required, qty_picked, is_verified)
					 VALUES ($1,$2,$3,0,false)
					 RETURNING `+itemCols,
					pl.ID, it.ProductID, it.QtyRequired,
				).Scan(&item.ID, &item.ProductID, &item.QtyRequired, &item.QtyPicked, &item.IsVerified); err != nil {
					return err
				}
				pl.Items = append(pl.Items, item)
			}
			return nil
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(pl)
	}
}

func HandleList(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		var pls []PickingList
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+plCols+` FROM picking_lists WHERE ($1='' OR status=$1) ORDER BY created_at DESC`, status)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var pl PickingList
				if err := rows.Scan(&pl.ID, &pl.WarehouseID, &pl.RouteID, &pl.Status, &pl.CreatedAt, &pl.UpdatedAt); err != nil {
					return err
				}
				pls = append(pls, pl)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if pls == nil {
			pls = []PickingList{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pls)
	}
}

func HandleGet(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var pl PickingList
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`SELECT `+plCols+` FROM picking_lists WHERE id=$1`, id,
			).Scan(&pl.ID, &pl.WarehouseID, &pl.RouteID, &pl.Status, &pl.CreatedAt, &pl.UpdatedAt); err != nil {
				return err
			}
			rows, err := tx.Query(r.Context(),
				`SELECT `+itemCols+` FROM picking_items WHERE picking_list_id=$1`, id)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var it Item
				if err := rows.Scan(&it.ID, &it.ProductID, &it.QtyRequired, &it.QtyPicked, &it.IsVerified); err != nil {
					return err
				}
				pl.Items = append(pl.Items, it)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if pl.Items == nil {
			pl.Items = []Item{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pl)
	}
}

func HandleVerifyItem(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID := chi.URLParam(r, "itemId")
		var body struct {
			QtyPicked string `json:"qty_picked"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.QtyPicked == "" {
			body.QtyPicked = "0"
		}

		var it Item
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`UPDATE picking_items SET is_verified=true, qty_picked=$1
				 WHERE id=$2
				 RETURNING `+itemCols,
				body.QtyPicked, itemID,
			).Scan(&it.ID, &it.ProductID, &it.QtyRequired, &it.QtyPicked, &it.IsVerified)
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(it)
	}
}

func HandleComplete(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			var status string
			if err := tx.QueryRow(r.Context(),
				`SELECT status FROM picking_lists WHERE id=$1`, id,
			).Scan(&status); err != nil {
				return err
			}
			if status != "pending" {
				return errors.New("picking list not pending")
			}
			var unverified int
			if err := tx.QueryRow(r.Context(),
				`SELECT COUNT(*) FROM picking_items WHERE picking_list_id=$1 AND is_verified=false`, id,
			).Scan(&unverified); err != nil {
				return err
			}
			if unverified > 0 {
				return errors.New("not all items verified")
			}
			_, err := tx.Exec(r.Context(),
				`UPDATE picking_lists SET status='completed', updated_at=NOW() WHERE id=$1`, id)
			return err
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
