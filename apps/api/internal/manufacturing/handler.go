package manufacturing

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

type ManufacturingOrder struct {
	ID           string     `json:"id"`
	ProductID    string     `json:"product_id"`
	WarehouseID  string     `json:"warehouse_id"`
	QtyToProduce string     `json:"qty_to_produce"`
	Status       string     `json:"status"`
	StartedAt    *time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	Notes        string     `json:"notes"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	BOM          []BOMItem  `json:"bom,omitempty"`
}

type BOMItem struct {
	ID                 string `json:"id"`
	ComponentProductID string `json:"component_product_id"`
	QtyPerUnit         string `json:"qty_per_unit"`
}

const moCols = `id, product_id::text, warehouse_id::text, qty_to_produce::text, status, started_at, completed_at, COALESCE(notes,''), created_at, updated_at`
const bomCols = `id, component_product_id::text, qty_per_unit::text`

type bomInput struct {
	ComponentProductID string `json:"component_product_id"`
	QtyPerUnit         string `json:"qty_per_unit"`
}

func HandleCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			ProductID    string      `json:"product_id"`
			WarehouseID  string      `json:"warehouse_id"`
			QtyToProduce string      `json:"qty_to_produce"`
			Notes        string      `json:"notes"`
			BOM          []bomInput  `json:"bom"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.ProductID == "" || body.WarehouseID == "" {
			http.Error(w, "product_id and warehouse_id are required", http.StatusBadRequest)
			return
		}
		if body.QtyToProduce == "" {
			body.QtyToProduce = "1"
		}
		if len(body.BOM) == 0 {
			http.Error(w, "bom required", http.StatusBadRequest)
			return
		}

		var mo ManufacturingOrder
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`INSERT INTO manufacturing_orders (tenant_id, product_id, warehouse_id, qty_to_produce, status, notes)
				 VALUES ($1,$2,$3,$4,'pending',$5)
				 RETURNING `+moCols,
				mw.TenantIDFrom(r.Context()), body.ProductID, body.WarehouseID, body.QtyToProduce, body.Notes,
			).Scan(&mo.ID, &mo.ProductID, &mo.WarehouseID, &mo.QtyToProduce, &mo.Status, &mo.StartedAt, &mo.CompletedAt, &mo.Notes, &mo.CreatedAt, &mo.UpdatedAt); err != nil {
				return err
			}
			for _, b := range body.BOM {
				if b.QtyPerUnit == "" {
					b.QtyPerUnit = "1"
				}
				var bom BOMItem
				if err := tx.QueryRow(r.Context(),
					`INSERT INTO manufacturing_bom (order_id, component_product_id, qty_per_unit)
					 VALUES ($1,$2,$3)
					 RETURNING `+bomCols,
					mo.ID, b.ComponentProductID, b.QtyPerUnit,
				).Scan(&bom.ID, &bom.ComponentProductID, &bom.QtyPerUnit); err != nil {
					return err
				}
				mo.BOM = append(mo.BOM, bom)
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
		json.NewEncoder(w).Encode(mo)
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
		var mos []ManufacturingOrder
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+moCols+` FROM manufacturing_orders WHERE ($1='' OR status=$1) ORDER BY created_at DESC LIMIT $2 OFFSET $3`, status, perPage, offset)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var mo ManufacturingOrder
				if err := rows.Scan(&mo.ID, &mo.ProductID, &mo.WarehouseID, &mo.QtyToProduce, &mo.Status, &mo.StartedAt, &mo.CompletedAt, &mo.Notes, &mo.CreatedAt, &mo.UpdatedAt); err != nil {
					return err
				}
				mos = append(mos, mo)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if mos == nil {
			mos = []ManufacturingOrder{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mos)
	}
}

func HandleGet(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var mo ManufacturingOrder
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`SELECT `+moCols+` FROM manufacturing_orders WHERE id=$1`, id,
			).Scan(&mo.ID, &mo.ProductID, &mo.WarehouseID, &mo.QtyToProduce, &mo.Status, &mo.StartedAt, &mo.CompletedAt, &mo.Notes, &mo.CreatedAt, &mo.UpdatedAt); err != nil {
				return err
			}
			rows, err := tx.Query(r.Context(),
				`SELECT `+bomCols+` FROM manufacturing_bom WHERE order_id=$1`, id)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var b BOMItem
				if err := rows.Scan(&b.ID, &b.ComponentProductID, &b.QtyPerUnit); err != nil {
					return err
				}
				mo.BOM = append(mo.BOM, b)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if mo.BOM == nil {
			mo.BOM = []BOMItem{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mo)
	}
}

func HandleStart(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(),
				`UPDATE manufacturing_orders SET status='in_progress', started_at=NOW(), updated_at=NOW()
				 WHERE id=$1 AND status='pending'`, id)
			if err != nil {
				return err
			}
			if ct.RowsAffected() == 0 {
				return errNotFound
			}
			return nil
		})
		if err != nil {
			http.Error(w, "not found or not pending", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleComplete(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			var status, productID, warehouseID string
			var qtyToProduce string
			if err := tx.QueryRow(r.Context(),
				`SELECT status, product_id::text, warehouse_id::text, qty_to_produce::text FROM manufacturing_orders WHERE id=$1`, id,
			).Scan(&status, &productID, &warehouseID, &qtyToProduce); err != nil {
				return err
			}
			if status != "in_progress" {
				return errors.New("order not in progress")
			}
			if _, err := tx.Exec(r.Context(),
				`INSERT INTO stock_movements (tenant_id, warehouse_id, product_id, type, qty, reference_type, reference_id)
				 VALUES ($1,$2,$3,'in',$4,'manufacturing',$5::uuid)`,
				mw.TenantIDFrom(r.Context()), warehouseID, productID, qtyToProduce, id); err != nil {
				return err
			}
			rows, err := tx.Query(r.Context(),
				`SELECT component_product_id::text, qty_per_unit::text FROM manufacturing_bom WHERE order_id=$1`, id)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var compID, qtyPerUnit string
				if err := rows.Scan(&compID, &qtyPerUnit); err != nil {
					return err
				}
				totalQty := qtyPerUnit + " * " + qtyToProduce
				if _, err := tx.Exec(r.Context(),
					`INSERT INTO stock_movements (tenant_id, warehouse_id, product_id, type, qty, reference_type, reference_id)
					 VALUES ($1,$2,$3,'out',-($4),$5,'manufacturing',$6::uuid)`,
					mw.TenantIDFrom(r.Context()), warehouseID, compID, totalQty, "manufacturing", id); err != nil {
					return err
				}
			}
			if err := rows.Err(); err != nil {
				return err
			}
			_, err = tx.Exec(r.Context(),
				`UPDATE manufacturing_orders SET status='completed', completed_at=NOW(), updated_at=NOW() WHERE id=$1`, id)
			return err
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
