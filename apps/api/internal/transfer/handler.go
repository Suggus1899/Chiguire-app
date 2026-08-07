package transfer

import (
	"log"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

var errNotFound = errors.New("not found")

type Transfer struct {
	ID              string    `json:"id"`
	FromWarehouseID string    `json:"from_warehouse_id"`
	ToWarehouseID   string    `json:"to_warehouse_id"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Items           []Item    `json:"items,omitempty"`
}

type Item struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Qty       string `json:"qty"`
}

const trCols = `id, from_warehouse_id::text, to_warehouse_id::text, status, COALESCE(notes,''), created_at, updated_at`
const itemCols = `id, product_id::text, qty::text`

type itemInput struct {
	ProductID string `json:"product_id"`
	Qty       string `json:"qty"`
}

func HandleCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			FromWarehouseID string      `json:"from_warehouse_id"`
			ToWarehouseID   string      `json:"to_warehouse_id"`
			Notes           string      `json:"notes"`
			Items           []itemInput `json:"items"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.FromWarehouseID == "" || body.ToWarehouseID == "" {
			http.Error(w, "from_warehouse_id and to_warehouse_id are required", http.StatusBadRequest)
			return
		}
		if body.FromWarehouseID == body.ToWarehouseID {
			http.Error(w, "from and to warehouses must differ", http.StatusBadRequest)
			return
		}
		if len(body.Items) == 0 {
			http.Error(w, "items required", http.StatusBadRequest)
			return
		}

		var tr Transfer
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`INSERT INTO inventory_transfers (tenant_id, from_warehouse_id, to_warehouse_id, status, notes)
				 VALUES ($1,$2,$3,'pending',$4)
				 RETURNING `+trCols,
				mw.TenantIDFrom(r.Context()), body.FromWarehouseID, body.ToWarehouseID, body.Notes,
			).Scan(&tr.ID, &tr.FromWarehouseID, &tr.ToWarehouseID, &tr.Status, &tr.Notes, &tr.CreatedAt, &tr.UpdatedAt); err != nil {
				return err
			}
			for _, it := range body.Items {
				if it.Qty == "" {
					it.Qty = "0"
				}
				var item Item
				if err := tx.QueryRow(r.Context(),
					`INSERT INTO inventory_transfer_items (transfer_id, product_id, qty)
					 VALUES ($1,$2,$3)
					 RETURNING `+itemCols,
					tr.ID, it.ProductID, it.Qty,
				).Scan(&item.ID, &item.ProductID, &item.Qty); err != nil {
					return err
				}
				tr.Items = append(tr.Items, item)
			}
			return nil
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tr)
	}
}

func HandleList(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		if perPage < 1 {
			perPage = 50
		}
		if perPage > 200 {
			perPage = 200
		}
		offset := (page - 1) * perPage
		var trs []Transfer
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+trCols+` FROM inventory_transfers WHERE ($1='' OR status=$1) ORDER BY created_at DESC LIMIT $2 OFFSET $3`, status, perPage, offset)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var tr Transfer
				if err := rows.Scan(&tr.ID, &tr.FromWarehouseID, &tr.ToWarehouseID, &tr.Status, &tr.Notes, &tr.CreatedAt, &tr.UpdatedAt); err != nil {
					return err
				}
				trs = append(trs, tr)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if trs == nil {
			trs = []Transfer{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trs)
	}
}

func HandleGet(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var tr Transfer
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`SELECT `+trCols+` FROM inventory_transfers WHERE id=$1`, id,
			).Scan(&tr.ID, &tr.FromWarehouseID, &tr.ToWarehouseID, &tr.Status, &tr.Notes, &tr.CreatedAt, &tr.UpdatedAt); err != nil {
				return err
			}
			rows, err := tx.Query(r.Context(),
				`SELECT `+itemCols+` FROM inventory_transfer_items WHERE transfer_id=$1`, id)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var it Item
				if err := rows.Scan(&it.ID, &it.ProductID, &it.Qty); err != nil {
					return err
				}
				tr.Items = append(tr.Items, it)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if tr.Items == nil {
			tr.Items = []Item{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tr)
	}
}

func HandleShip(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			var status, fromWh string
			if err := tx.QueryRow(r.Context(),
				`SELECT status, from_warehouse_id::text FROM inventory_transfers WHERE id=$1`, id,
			).Scan(&status, &fromWh); err != nil {
				return err
			}
			if status != "pending" {
				return errors.New("transfer not pending")
			}
			rows, err := tx.Query(r.Context(),
				`SELECT product_id::text, qty::text FROM inventory_transfer_items WHERE transfer_id=$1`, id)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var pid, qty string
				if err := rows.Scan(&pid, &qty); err != nil {
					return err
				}
				if _, err := tx.Exec(r.Context(),
					`INSERT INTO stock_movements (tenant_id, warehouse_id, product_id, type, qty, reference_type, reference_id)
					 VALUES ($1,$2,$3,'out',-$4,'transfer',$5::uuid)`,
					mw.TenantIDFrom(r.Context()), fromWh, pid, qty, id); err != nil {
					return err
				}
			}
			if err := rows.Err(); err != nil {
				return err
			}
			_, err = tx.Exec(r.Context(),
				`UPDATE inventory_transfers SET status='shipped', updated_at=NOW() WHERE id=$1`, id)
			return err
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleReceive(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			var status, toWh string
			if err := tx.QueryRow(r.Context(),
				`SELECT status, to_warehouse_id::text FROM inventory_transfers WHERE id=$1`, id,
			).Scan(&status, &toWh); err != nil {
				return err
			}
			if status != "shipped" {
				return errors.New("transfer not shipped")
			}
			rows, err := tx.Query(r.Context(),
				`SELECT product_id::text, qty::text FROM inventory_transfer_items WHERE transfer_id=$1`, id)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var pid, qty string
				if err := rows.Scan(&pid, &qty); err != nil {
					return err
				}
				if _, err := tx.Exec(r.Context(),
					`INSERT INTO stock_movements (tenant_id, warehouse_id, product_id, type, qty, reference_type, reference_id)
					 VALUES ($1,$2,$3,'in',$4,'transfer',$5::uuid)`,
					mw.TenantIDFrom(r.Context()), toWh, pid, qty, id); err != nil {
					return err
				}
			}
			if err := rows.Err(); err != nil {
				return err
			}
			_, err = tx.Exec(r.Context(),
				`UPDATE inventory_transfers SET status='received', updated_at=NOW() WHERE id=$1`, id)
			return err
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
