package purchase

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

type PurchaseOrder struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id,omitempty"`
	VendorID   string    `json:"vendor_id"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Currency   string    `json:"currency"`
	Subtotal   float64   `json:"subtotal"`
	TaxTotal   float64   `json:"tax_total"`
	Total      float64   `json:"total"`
	CreatedAt  time.Time `json:"created_at"`
}

type POItem struct {
	ID             string  `json:"id"`
	PurchaseOrderID string `json:"purchase_order_id"`
	ProductID      string  `json:"product_id"`
	Description    string  `json:"description"`
	QtyOrdered     float64 `json:"qty_ordered"`
	QtyReceived    float64 `json:"qty_received"`
	UnitCost       float64 `json:"unit_cost"`
	LineTotal      float64 `json:"line_total"`
}

func HandleCreatePO(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			VendorID string  `json:"vendor_id"`
			Number   string  `json:"number"`
			Currency string  `json:"currency"`
			Subtotal float64 `json:"subtotal"`
			TaxTotal float64 `json:"tax_total"`
			Total    float64 `json:"total"`
			Items    []struct {
				ProductID   string  `json:"product_id"`
				Description string  `json:"description"`
				QtyOrdered  float64 `json:"qty_ordered"`
				UnitCost    float64 `json:"unit_cost"`
				LineTotal   float64 `json:"line_total"`
			} `json:"items"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.VendorID == "" || body.Number == "" {
			http.Error(w, "vendor_id and number required", http.StatusBadRequest)
			return
		}
		if body.Currency == "" {
			body.Currency = "VES"
		}

		tid := mw.TenantIDFrom(r.Context())
		var po PurchaseOrder
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`INSERT INTO purchase_orders (tenant_id, vendor_id, number, status, currency, subtotal, tax_total, total)
				 VALUES ($1,$2,$3,'draft',$4,$5,$6,$7)
				 RETURNING id, tenant_id, vendor_id, number, status, currency, subtotal, tax_total, total, created_at`,
				tid, body.VendorID, body.Number, body.Currency, body.Subtotal, body.TaxTotal, body.Total,
			).Scan(&po.ID, &po.TenantID, &po.VendorID, &po.Number, &po.Status, &po.Currency,
				&po.Subtotal, &po.TaxTotal, &po.Total, &po.CreatedAt); err != nil {
				return err
			}
			for _, it := range body.Items {
				if _, err := tx.Exec(r.Context(),
					`INSERT INTO purchase_order_items (tenant_id, purchase_order_id, product_id, description, qty_ordered, qty_received, unit_cost, line_total)
					 VALUES ($1,$2,$3,$4,$5,0,$6,$7)`,
					tid, po.ID, it.ProductID, it.Description, it.QtyOrdered, it.UnitCost, it.LineTotal,
				); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(po)
	}
}

func HandleListPOs(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		var pos []PurchaseOrder
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT id, tenant_id, vendor_id, number, status, currency, subtotal, tax_total, total, created_at
				 FROM purchase_orders WHERE ($1 = '' OR status = $1) ORDER BY created_at DESC`, status)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var po PurchaseOrder
				if err := rows.Scan(&po.ID, &po.TenantID, &po.VendorID, &po.Number, &po.Status, &po.Currency,
					&po.Subtotal, &po.TaxTotal, &po.Total, &po.CreatedAt); err != nil {
					return err
				}
				pos = append(pos, po)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if pos == nil {
			pos = []PurchaseOrder{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pos)
	}
}

func HandleGetPO(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var po PurchaseOrder
		var items []POItem
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`SELECT id, tenant_id, vendor_id, number, status, currency, subtotal, tax_total, total, created_at
				 FROM purchase_orders WHERE id=$1`, id,
			).Scan(&po.ID, &po.TenantID, &po.VendorID, &po.Number, &po.Status, &po.Currency,
				&po.Subtotal, &po.TaxTotal, &po.Total, &po.CreatedAt); err != nil {
				return err
			}
			rows, err := tx.Query(r.Context(),
				`SELECT id, purchase_order_id, product_id, description, qty_ordered, qty_received, unit_cost, line_total
				 FROM purchase_order_items WHERE purchase_order_id=$1`, id)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var it POItem
				if err := rows.Scan(&it.ID, &it.PurchaseOrderID, &it.ProductID, &it.Description,
					&it.QtyOrdered, &it.QtyReceived, &it.UnitCost, &it.LineTotal); err != nil {
					return err
				}
				items = append(items, it)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if items == nil {
			items = []POItem{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"po": po, "items": items})
	}
}

func HandleUpdatePO(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			VendorID string  `json:"vendor_id"`
			Currency string  `json:"currency"`
			Subtotal float64 `json:"subtotal"`
			TaxTotal float64 `json:"tax_total"`
			Total    float64 `json:"total"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			_, err := tx.Exec(r.Context(),
				`UPDATE purchase_orders SET vendor_id=$2, currency=$3, subtotal=$4, tax_total=$5, total=$6
				 WHERE id=$1`,
				id, body.VendorID, body.Currency, body.Subtotal, body.TaxTotal, body.Total)
			return err
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleApprovePO(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			_, err := tx.Exec(r.Context(),
				`UPDATE purchase_orders SET status='approved' WHERE id=$1`, id)
			return err
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleAddPOItem(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			ProductID   string  `json:"product_id"`
			Description string  `json:"description"`
			QtyOrdered  float64 `json:"qty_ordered"`
			UnitCost    float64 `json:"unit_cost"`
			LineTotal   float64 `json:"line_total"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		tid := mw.TenantIDFrom(r.Context())
		var item POItem
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO purchase_order_items (tenant_id, purchase_order_id, product_id, description, qty_ordered, qty_received, unit_cost, line_total)
				 VALUES ($1,$2,$3,$4,$5,0,$6,$7)
				 RETURNING id, purchase_order_id, product_id, description, qty_ordered, qty_received, unit_cost, line_total`,
				tid, id, body.ProductID, body.Description, body.QtyOrdered, body.UnitCost, body.LineTotal,
			).Scan(&item.ID, &item.PurchaseOrderID, &item.ProductID, &item.Description,
				&item.QtyOrdered, &item.QtyReceived, &item.UnitCost, &item.LineTotal)
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(item)
	}
}

func HandleReceivePOItem(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		poID := chi.URLParam(r, "id")
		itemID := chi.URLParam(r, "itemId")
		var body struct {
			WarehouseID string  `json:"warehouse_id"`
			QtyReceived float64 `json:"qty_received"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		tid := mw.TenantIDFrom(r.Context())
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			var productID string
			if err := tx.QueryRow(r.Context(),
				`UPDATE purchase_order_items SET qty_received = qty_received + $2
				 WHERE id=$1 AND purchase_order_id=$3 RETURNING product_id`,
				itemID, body.QtyReceived, poID,
			).Scan(&productID); err != nil {
				return err
			}
			if _, err := tx.Exec(r.Context(),
				`INSERT INTO purchase_order_receipts (tenant_id, purchase_order_id, po_item_id, warehouse_id, qty_received)
				 VALUES ($1,$2,$3,$4,$5)`,
				tid, poID, itemID, body.WarehouseID, body.QtyReceived,
			); err != nil {
				return err
			}
			if _, err := tx.Exec(r.Context(),
				`INSERT INTO stock_movements (tenant_id, product_id, warehouse_id, movement_type, qty, ref_type, ref_id)
				 VALUES ($1,$2,$3,'in',$4,'purchase_order',$5)`,
				tid, productID, body.WarehouseID, body.QtyReceived, poID,
			); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
