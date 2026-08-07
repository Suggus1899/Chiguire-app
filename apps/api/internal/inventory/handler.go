package inventory

import (
	"log"
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

// ---------- Branches ----------

type Branch struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

const branchCols = `id, name, COALESCE(address,''), COALESCE(phone,''), created_at`

func HandleCreateBranch(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name    string `json:"name"`
			Address string `json:"address"`
			Phone   string `json:"phone"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		var b Branch
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO branches (tenant_id, name, address, phone)
				 VALUES ($1,$2,$3,$4)
				 RETURNING `+branchCols,
				mw.TenantIDFrom(r.Context()), body.Name, body.Address, body.Phone,
			).Scan(&b.ID, &b.Name, &b.Address, &b.Phone, &b.CreatedAt)
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(b)
	}
}

func HandleListBranches(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var branches []Branch
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+branchCols+` FROM branches ORDER BY name`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var b Branch
				if err := rows.Scan(&b.ID, &b.Name, &b.Address, &b.Phone, &b.CreatedAt); err != nil {
					return err
				}
				branches = append(branches, b)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if branches == nil {
			branches = []Branch{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(branches)
	}
}

func HandleDeleteBranch(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(), `DELETE FROM branches WHERE id=$1`, id)
			if err != nil {
				return err
			}
			if ct.RowsAffected() == 0 {
				return errNotFound
			}
			return nil
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// ---------- Warehouses ----------

type Warehouse struct {
	ID        string    `json:"id"`
	BranchID  string    `json:"branch_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

const warehouseCols = `id, COALESCE(branch_id::text,''), name, created_at`

func HandleCreateWarehouse(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			BranchID string `json:"branch_id"`
			Name     string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		var wh Warehouse
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO warehouses (tenant_id, branch_id, name)
				 VALUES ($1,NULLIF($2,'')::uuid,$3)
				 RETURNING `+warehouseCols,
				mw.TenantIDFrom(r.Context()), body.BranchID, body.Name,
			).Scan(&wh.ID, &wh.BranchID, &wh.Name, &wh.CreatedAt)
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(wh)
	}
}

func HandleListWarehouses(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var warehouses []Warehouse
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+warehouseCols+` FROM warehouses ORDER BY name`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var wh Warehouse
				if err := rows.Scan(&wh.ID, &wh.BranchID, &wh.Name, &wh.CreatedAt); err != nil {
					return err
				}
				warehouses = append(warehouses, wh)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if warehouses == nil {
			warehouses = []Warehouse{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(warehouses)
	}
}

func HandleDeleteWarehouse(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(), `DELETE FROM warehouses WHERE id=$1`, id)
			if err != nil {
				return err
			}
			if ct.RowsAffected() == 0 {
				return errNotFound
			}
			return nil
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// ---------- Stock movements (append-only) ----------

type Movement struct {
	ID            string    `json:"id"`
	WarehouseID   string    `json:"warehouse_id"`
	ProductID     string    `json:"product_id"`
	Type          string    `json:"type"`
	Qty           string    `json:"qty"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   string    `json:"reference_id"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

const movementCols = `id, warehouse_id::text, product_id::text, type, qty::text, COALESCE(reference_type,''), COALESCE(reference_id::text,''), COALESCE(notes,''), created_at`

func HandleCreateMovement(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			WarehouseID   string `json:"warehouse_id"`
			ProductID     string `json:"product_id"`
			Type          string `json:"type"`
			Qty           string `json:"qty"`
			ReferenceType string `json:"reference_type"`
			ReferenceID   string `json:"reference_id"`
			Notes         string `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.WarehouseID == "" || body.ProductID == "" || body.Type == "" || body.Qty == "" {
			http.Error(w, "warehouse_id, product_id, type and qty are required", http.StatusBadRequest)
			return
		}

		var m Movement
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO stock_movements (tenant_id, warehouse_id, product_id, type, qty, reference_type, reference_id, notes)
				 VALUES ($1,$2,$3,$4,$5,$6,NULLIF($7,'')::uuid,$8)
				 RETURNING `+movementCols,
				mw.TenantIDFrom(r.Context()), body.WarehouseID, body.ProductID, body.Type, body.Qty, body.ReferenceType, body.ReferenceID, body.Notes,
			).Scan(&m.ID, &m.WarehouseID, &m.ProductID, &m.Type, &m.Qty, &m.ReferenceType, &m.ReferenceID, &m.Notes, &m.CreatedAt)
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(m)
	}
}

func HandleListMovements(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouseID := r.URL.Query().Get("warehouse_id")
		productID := r.URL.Query().Get("product_id")
		var movements []Movement
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+movementCols+` FROM stock_movements
				 WHERE ($1='' OR warehouse_id=$1::uuid) AND ($2='' OR product_id=$2::uuid)
				 ORDER BY created_at DESC`,
				warehouseID, productID)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var m Movement
				if err := rows.Scan(&m.ID, &m.WarehouseID, &m.ProductID, &m.Type, &m.Qty, &m.ReferenceType, &m.ReferenceID, &m.Notes, &m.CreatedAt); err != nil {
					return err
				}
				movements = append(movements, m)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if movements == nil {
			movements = []Movement{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movements)
	}
}

type StockLevel struct {
	WarehouseID string `json:"warehouse_id"`
	ProductID   string `json:"product_id"`
	QtyOnHand   string `json:"qty_on_hand"`
}

func HandleGetStock(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouseID := r.URL.Query().Get("warehouse_id")
		productID := r.URL.Query().Get("product_id")
		var levels []StockLevel
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT warehouse_id::text, product_id::text, COALESCE(SUM(qty),0)::text AS qty_on_hand
				 FROM stock_movements
				 WHERE ($1='' OR warehouse_id=$1::uuid) AND ($2='' OR product_id=$2::uuid)
				 GROUP BY warehouse_id, product_id
				 ORDER BY warehouse_id, product_id`,
				warehouseID, productID)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var s StockLevel
				if err := rows.Scan(&s.WarehouseID, &s.ProductID, &s.QtyOnHand); err != nil {
					return err
				}
				levels = append(levels, s)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if levels == nil {
			levels = []StockLevel{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(levels)
	}
}
