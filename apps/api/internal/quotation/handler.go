package quotation

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

type Quotation struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id,omitempty"`
	Number            string     `json:"number"`
	CustomerID        string     `json:"customer_id"`
	Status            string     `json:"status"`
	Currency          string     `json:"currency"`
	Subtotal          float64    `json:"subtotal"`
	DiscountTotal     float64    `json:"discount_total"`
	TaxTotal          float64    `json:"tax_total"`
	Total             float64    `json:"total"`
	ValidDays         int        `json:"valid_days"`
	ValidUntil        *time.Time `json:"valid_until,omitempty"`
	ConvertedInvoiceID *string   `json:"converted_invoice_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

type QuotationItem struct {
	ID          string  `json:"id"`
	QuotationID string  `json:"quotation_id"`
	ProductID   string  `json:"product_id"`
	Description string  `json:"description"`
	Qty         float64 `json:"qty"`
	UnitPrice   float64 `json:"unit_price"`
	DiscountPct float64 `json:"discount_pct"`
	LineTotal   float64 `json:"line_total"`
}

func HandleCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Number        string  `json:"number"`
			CustomerID    string  `json:"customer_id"`
			Currency      string  `json:"currency"`
			Subtotal      float64 `json:"subtotal"`
			DiscountTotal float64 `json:"discount_total"`
			TaxTotal      float64 `json:"tax_total"`
			Total         float64 `json:"total"`
			ValidDays     int     `json:"valid_days"`
			Items         []struct {
				ProductID   string  `json:"product_id"`
				Description string  `json:"description"`
				Qty         float64 `json:"qty"`
				UnitPrice   float64 `json:"unit_price"`
				DiscountPct float64 `json:"discount_pct"`
				LineTotal   float64 `json:"line_total"`
			} `json:"items"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.CustomerID == "" || body.Number == "" {
			http.Error(w, "customer_id and number required", http.StatusBadRequest)
			return
		}
		if body.Currency == "" {
			body.Currency = "VES"
		}
		if body.ValidDays == 0 {
			body.ValidDays = 15
		}
		validUntil := time.Now().AddDate(0, 0, body.ValidDays)

		tid := mw.TenantIDFrom(r.Context())
		var q Quotation
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`INSERT INTO quotations (tenant_id, number, customer_id, status, currency, subtotal, discount_total, tax_total, total, valid_days, valid_until)
				 VALUES ($1,$2,$3,'draft',$4,$5,$6,$7,$8,$9,$10)
				 RETURNING id, tenant_id, number, customer_id, status, currency, subtotal, discount_total, tax_total, total, valid_days, valid_until, created_at`,
				tid, body.Number, body.CustomerID, body.Currency, body.Subtotal, body.DiscountTotal,
				body.TaxTotal, body.Total, body.ValidDays, validUntil,
			).Scan(&q.ID, &q.TenantID, &q.Number, &q.CustomerID, &q.Status, &q.Currency,
				&q.Subtotal, &q.DiscountTotal, &q.TaxTotal, &q.Total, &q.ValidDays, &q.ValidUntil, &q.CreatedAt); err != nil {
				return err
			}
			for _, it := range body.Items {
				if _, err := tx.Exec(r.Context(),
					`INSERT INTO quotation_items (tenant_id, quotation_id, product_id, description, qty, unit_price, discount_pct, line_total)
					 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
					tid, q.ID, it.ProductID, it.Description, it.Qty, it.UnitPrice, it.DiscountPct, it.LineTotal,
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
		json.NewEncoder(w).Encode(q)
	}
}

func HandleList(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		var qs []Quotation
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT id, tenant_id, number, customer_id, status, currency, subtotal, discount_total, tax_total, total, valid_days, valid_until, converted_invoice_id, created_at
				 FROM quotations WHERE ($1 = '' OR status = $1) ORDER BY created_at DESC`, status)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var q Quotation
				if err := rows.Scan(&q.ID, &q.TenantID, &q.Number, &q.CustomerID, &q.Status, &q.Currency,
					&q.Subtotal, &q.DiscountTotal, &q.TaxTotal, &q.Total, &q.ValidDays, &q.ValidUntil,
					&q.ConvertedInvoiceID, &q.CreatedAt); err != nil {
					return err
				}
				qs = append(qs, q)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if qs == nil {
			qs = []Quotation{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(qs)
	}
}

func HandleGet(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var q Quotation
		var items []QuotationItem
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`SELECT id, tenant_id, number, customer_id, status, currency, subtotal, discount_total, tax_total, total, valid_days, valid_until, converted_invoice_id, created_at
				 FROM quotations WHERE id=$1`, id,
			).Scan(&q.ID, &q.TenantID, &q.Number, &q.CustomerID, &q.Status, &q.Currency,
				&q.Subtotal, &q.DiscountTotal, &q.TaxTotal, &q.Total, &q.ValidDays, &q.ValidUntil,
				&q.ConvertedInvoiceID, &q.CreatedAt); err != nil {
				return err
			}
			rows, err := tx.Query(r.Context(),
				`SELECT id, quotation_id, product_id, description, qty, unit_price, discount_pct, line_total
				 FROM quotation_items WHERE quotation_id=$1`, id)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var it QuotationItem
				if err := rows.Scan(&it.ID, &it.QuotationID, &it.ProductID, &it.Description,
					&it.Qty, &it.UnitPrice, &it.DiscountPct, &it.LineTotal); err != nil {
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
			items = []QuotationItem{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"quotation": q, "items": items})
	}
}

func HandleUpdate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			Currency      string  `json:"currency"`
			Subtotal      float64 `json:"subtotal"`
			DiscountTotal float64 `json:"discount_total"`
			TaxTotal      float64 `json:"tax_total"`
			Total         float64 `json:"total"`
			ValidDays     int     `json:"valid_days"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		var validUntil *time.Time
		if body.ValidDays > 0 {
			vu := time.Now().AddDate(0, 0, body.ValidDays)
			validUntil = &vu
		}
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			_, err := tx.Exec(r.Context(),
				`UPDATE quotations SET currency=$2, subtotal=$3, discount_total=$4, tax_total=$5, total=$6, valid_days=$7, valid_until=$8
				 WHERE id=$1`,
				id, body.Currency, body.Subtotal, body.DiscountTotal, body.TaxTotal, body.Total, body.ValidDays, validUntil)
			return err
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleAddItem(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			ProductID   string  `json:"product_id"`
			Description string  `json:"description"`
			Qty         float64 `json:"qty"`
			UnitPrice   float64 `json:"unit_price"`
			DiscountPct float64 `json:"discount_pct"`
			LineTotal   float64 `json:"line_total"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		tid := mw.TenantIDFrom(r.Context())
		var item QuotationItem
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO quotation_items (tenant_id, quotation_id, product_id, description, qty, unit_price, discount_pct, line_total)
				 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
				 RETURNING id, quotation_id, product_id, description, qty, unit_price, discount_pct, line_total`,
				tid, id, body.ProductID, body.Description, body.Qty, body.UnitPrice, body.DiscountPct, body.LineTotal,
			).Scan(&item.ID, &item.QuotationID, &item.ProductID, &item.Description,
				&item.Qty, &item.UnitPrice, &item.DiscountPct, &item.LineTotal)
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

func HandleRemoveItem(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID := chi.URLParam(r, "itemId")
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			_, err := tx.Exec(r.Context(), `DELETE FROM quotation_items WHERE id=$1`, itemID)
			return err
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleConvertToInvoice(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		tid := mw.TenantIDFrom(r.Context())
		var result struct {
			InvoiceID string `json:"invoice_id"`
		}
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			var q Quotation
			if err := tx.QueryRow(r.Context(),
				`SELECT id, number, customer_id, currency, subtotal, discount_total, tax_total, total
				 FROM quotations WHERE id=$1`, id,
			).Scan(&q.ID, &q.Number, &q.CustomerID, &q.Currency, &q.Subtotal, &q.DiscountTotal, &q.TaxTotal, &q.Total); err != nil {
				return err
			}
			if err := tx.QueryRow(r.Context(),
				`INSERT INTO invoices (tenant_id, number, customer_id, status, currency, subtotal, discount_total, tax_total, total)
				 VALUES ($1,$2,$3,'issued',$4,$5,$6,$7,$8)
				 RETURNING id`,
				tid, q.Number+"-INV", q.CustomerID, q.Currency, q.Subtotal, q.DiscountTotal, q.TaxTotal, q.Total,
			).Scan(&result.InvoiceID); err != nil {
				return err
			}
			rows, err := tx.Query(r.Context(),
				`SELECT product_id, description, qty, unit_price, discount_pct, line_total
				 FROM quotation_items WHERE quotation_id=$1`, id)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var productID, description string
				var qty, unitPrice, discountPct, lineTotal float64
				if err := rows.Scan(&productID, &description, &qty, &unitPrice, &discountPct, &lineTotal); err != nil {
					return err
				}
				if _, err := tx.Exec(r.Context(),
					`INSERT INTO invoice_items (tenant_id, invoice_id, product_id, description, qty, unit_price, discount_pct, line_total)
					 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
					tid, result.InvoiceID, productID, description, qty, unitPrice, discountPct, lineTotal,
				); err != nil {
					return err
				}
			}
			if _, err := tx.Exec(r.Context(),
				`UPDATE quotations SET status='converted', converted_invoice_id=$2 WHERE id=$1`,
				id, result.InvoiceID); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(result)
	}
}

func HandleSendQuote(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			_, err := tx.Exec(r.Context(), `UPDATE quotations SET status='sent' WHERE id=$1`, id)
			return err
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
	}
}
