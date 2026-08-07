package report

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

type SalesBookEntry struct {
	ID         string     `json:"id"`
	Number     string     `json:"number"`
	CustomerID string     `json:"customer_id"`
	Status     string     `json:"status"`
	Currency   string     `json:"currency"`
	Subtotal   string     `json:"subtotal"`
	TaxTotal   string     `json:"tax_total"`
	Total      string     `json:"total"`
	IssuedAt   *time.Time `json:"issued_at"`
}

func HandleSalesBook(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		var rows []SalesBookEntry
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			q, err := tx.Query(r.Context(),
				`SELECT id, COALESCE(number,''), COALESCE(customer_id::text,''), status, currency,
				        subtotal::text, tax_total::text, total::text, issued_at
				 FROM invoices
				 WHERE status IN ('issued','paid','partial')
				   AND ($1='' OR issued_at >= $1::timestamptz)
				   AND ($2='' OR issued_at <= $2::timestamptz)
				 ORDER BY issued_at DESC`, from, to)
			if err != nil {
				return err
			}
			defer q.Close()
			for q.Next() {
				var e SalesBookEntry
				if err := q.Scan(&e.ID, &e.Number, &e.CustomerID, &e.Status, &e.Currency, &e.Subtotal, &e.TaxTotal, &e.Total, &e.IssuedAt); err != nil {
					return err
				}
				rows = append(rows, e)
			}
			return q.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if rows == nil {
			rows = []SalesBookEntry{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}

type PurchasesBookEntry struct {
	ID        string    `json:"id"`
	VendorID  string    `json:"vendor_id"`
	Status    string    `json:"status"`
	Total     string    `json:"total"`
	CreatedAt time.Time `json:"created_at"`
}

func HandlePurchasesBook(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		var rows []PurchasesBookEntry
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			q, err := tx.Query(r.Context(),
				`SELECT id, COALESCE(vendor_id::text,''), status, total::text, created_at
				 FROM purchase_orders
				 WHERE status IN ('approved','received')
				   AND ($1='' OR created_at >= $1::timestamptz)
				   AND ($2='' OR created_at <= $2::timestamptz)
				 ORDER BY created_at DESC`, from, to)
			if err != nil {
				return err
			}
			defer q.Close()
			for q.Next() {
				var e PurchasesBookEntry
				if err := q.Scan(&e.ID, &e.VendorID, &e.Status, &e.Total, &e.CreatedAt); err != nil {
					return err
				}
				rows = append(rows, e)
			}
			return q.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if rows == nil {
			rows = []PurchasesBookEntry{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}

type InventoryCurrentRow struct {
	WarehouseID string `json:"warehouse_id"`
	ProductID   string `json:"product_id"`
	QtyOnHand   string `json:"qty_on_hand"`
}

func HandleInventoryCurrent(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouseID := r.URL.Query().Get("warehouse_id")
		productID := r.URL.Query().Get("product_id")
		var rows []InventoryCurrentRow
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			q, err := tx.Query(r.Context(),
				`SELECT warehouse_id::text, product_id::text, COALESCE(SUM(qty),0)::text AS qty_on_hand
				 FROM stock_movements
				 WHERE ($1='' OR warehouse_id=$1::uuid) AND ($2='' OR product_id=$2::uuid)
				 GROUP BY warehouse_id, product_id
				 ORDER BY warehouse_id, product_id`, warehouseID, productID)
			if err != nil {
				return err
			}
			defer q.Close()
			for q.Next() {
				var row InventoryCurrentRow
				if err := q.Scan(&row.WarehouseID, &row.ProductID, &row.QtyOnHand); err != nil {
					return err
				}
				rows = append(rows, row)
			}
			return q.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if rows == nil {
			rows = []InventoryCurrentRow{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}

type InventoryValuedRow struct {
	WarehouseID string `json:"warehouse_id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	QtyOnHand   string `json:"qty_on_hand"`
	Cost        string `json:"cost"`
	SalePrice   string `json:"sale_price"`
	TotalCost   string `json:"total_cost"`
	TotalSale   string `json:"total_sale"`
}

func HandleInventoryValued(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rows []InventoryValuedRow
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			q, err := tx.Query(r.Context(),
				`SELECT sm.warehouse_id::text, sm.product_id::text, p.name,
				        COALESCE(SUM(sm.qty),0)::text AS qty_on_hand,
				        COALESCE(p.avg_cost_usd,0)::text AS cost,
				        COALESCE(pp.amount,0)::text AS sale_price,
				        COALESCE(SUM(sm.qty)*p.avg_cost_usd,0)::text AS total_cost,
				        COALESCE(SUM(sm.qty)*pp.amount,0)::text AS total_sale
				 FROM stock_movements sm
				 JOIN products p ON p.id = sm.product_id
				 LEFT JOIN LATERAL (
				   SELECT amount FROM product_prices pp2
				   WHERE pp2.product_id = sm.product_id
				   ORDER BY valid_from DESC LIMIT 1
				 ) pp ON true
				 GROUP BY sm.warehouse_id, sm.product_id, p.name, p.avg_cost_usd, pp.amount
				 ORDER BY sm.warehouse_id, p.name`)
			if err != nil {
				return err
			}
			defer q.Close()
			for q.Next() {
				var row InventoryValuedRow
				if err := q.Scan(&row.WarehouseID, &row.ProductID, &row.ProductName, &row.QtyOnHand, &row.Cost, &row.SalePrice, &row.TotalCost, &row.TotalSale); err != nil {
					return err
				}
				rows = append(rows, row)
			}
			return q.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if rows == nil {
			rows = []InventoryValuedRow{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}

type KardexRow struct {
	ID          string    `json:"id"`
	WarehouseID string    `json:"warehouse_id"`
	ProductID   string    `json:"product_id"`
	Type        string    `json:"type"`
	Qty         string    `json:"qty"`
	ReferenceType string  `json:"reference_type"`
	CreatedAt   time.Time `json:"created_at"`
}

func HandleKardex(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		productID := r.URL.Query().Get("product_id")
		var rows []KardexRow
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			q, err := tx.Query(r.Context(),
				`SELECT id, warehouse_id::text, product_id::text, type, qty::text, COALESCE(reference_type,''), created_at
				 FROM stock_movements
				 WHERE ($1='' OR created_at >= $1::timestamptz)
				   AND ($2='' OR created_at <= $2::timestamptz)
				   AND ($3='' OR product_id=$3::uuid)
				 ORDER BY created_at DESC`, from, to, productID)
			if err != nil {
				return err
			}
			defer q.Close()
			for q.Next() {
				var row KardexRow
				if err := q.Scan(&row.ID, &row.WarehouseID, &row.ProductID, &row.Type, &row.Qty, &row.ReferenceType, &row.CreatedAt); err != nil {
					return err
				}
				rows = append(rows, row)
			}
			return q.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if rows == nil {
			rows = []KardexRow{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}

type Art177Row struct {
	WarehouseID string `json:"warehouse_id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	QtyOnHand   string `json:"qty_on_hand"`
	HistoricalCost string `json:"historical_cost_bs"`
	TotalBs     string `json:"total_bs"`
}

func HandleArt177(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		var rows []Art177Row
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			q, err := tx.Query(r.Context(),
				`SELECT sm.warehouse_id::text, sm.product_id::text, p.name,
				        COALESCE(SUM(sm.qty),0)::text AS qty_on_hand,
				        COALESCE(p.avg_cost_usd,0)::text AS historical_cost_bs,
				        COALESCE(SUM(sm.qty)*p.avg_cost_usd,0)::text AS total_bs
				 FROM stock_movements sm
				 JOIN products p ON p.id = sm.product_id
				 WHERE ($1='' OR sm.created_at >= $1::timestamptz)
				   AND ($2='' OR sm.created_at <= $2::timestamptz)
				 GROUP BY sm.warehouse_id, sm.product_id, p.name, p.avg_cost_usd
				 ORDER BY sm.warehouse_id, p.name`, from, to)
			if err != nil {
				return err
			}
			defer q.Close()
			for q.Next() {
				var row Art177Row
				if err := q.Scan(&row.WarehouseID, &row.ProductID, &row.ProductName, &row.QtyOnHand, &row.HistoricalCost, &row.TotalBs); err != nil {
					return err
				}
				rows = append(rows, row)
			}
			return q.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if rows == nil {
			rows = []Art177Row{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}

type IGTFRow struct {
	ID            string    `json:"id"`
	InvoiceID     string    `json:"invoice_id"`
	Method        string    `json:"method"`
	AmountUSD     string    `json:"amount_usd"`
	IgtAmount     string    `json:"igt_amount"`
	PaidAt        time.Time `json:"paid_at"`
}

func HandleIGTFReport(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		var rows []IGTFRow
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			q, err := tx.Query(r.Context(),
				`SELECT id, COALESCE(invoice_id::text,''), method, amount_usd::text, igt_amount::text, paid_at
				 FROM invoice_payments
				 WHERE igt_amount > 0
				   AND ($1='' OR paid_at >= $1::timestamptz)
				   AND ($2='' OR paid_at <= $2::timestamptz)
				 ORDER BY paid_at DESC`, from, to)
			if err != nil {
				return err
			}
			defer q.Close()
			for q.Next() {
				var row IGTFRow
				if err := q.Scan(&row.ID, &row.InvoiceID, &row.Method, &row.AmountUSD, &row.IgtAmount, &row.PaidAt); err != nil {
					return err
				}
				rows = append(rows, row)
			}
			return q.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if rows == nil {
			rows = []IGTFRow{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}
