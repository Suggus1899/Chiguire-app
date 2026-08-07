package invoice

import (
	"log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

var errNotFound = errors.New("not found")

type Invoice struct {
	ID              string     `json:"id"`
	Number          string     `json:"number"`
	CustomerID      string     `json:"customer_id"`
	VendorUserID    string     `json:"vendor_user_id"`
	BranchID        string     `json:"branch_id"`
	Status          string     `json:"status"`
	Currency        string     `json:"currency"`
	Subtotal        string     `json:"subtotal"`
	DiscountTotal   string     `json:"discount_total"`
	TaxTotal        string     `json:"tax_total"`
	Total           string     `json:"total"`
	ExchangeRateID  string     `json:"exchange_rate_id"`
	IssuedAt        *time.Time `json:"issued_at"`
	VoidedAt        *time.Time `json:"voided_at"`
	Notes           string     `json:"notes"`
	PendingEmission bool       `json:"pending_emission"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Items           []Item     `json:"items,omitempty"`
	Payments        []Payment  `json:"payments,omitempty"`
}

type Item struct {
	ID             string `json:"id"`
	ProductID      string `json:"product_id"`
	Description    string `json:"description"`
	Qty            string `json:"qty"`
	UnitPrice      string `json:"unit_price"`
	DiscountPct    string `json:"discount_pct"`
	DiscountAmount string `json:"discount_amount"`
	TaxRate        string `json:"tax_rate"`
	TaxAmount      string `json:"tax_amount"`
	LineTotal      string `json:"line_total"`
}

type Payment struct {
	ID            string    `json:"id"`
	Method        string    `json:"method"`
	Provider      string    `json:"provider"`
	AmountUSD     string    `json:"amount_usd"`
	AmountVES     string    `json:"amount_ves"`
	RateAtPayment string    `json:"rate_at_payment"`
	IgtAmount     string    `json:"igt_amount"`
	Reference     string    `json:"reference"`
	Status        string    `json:"status"`
	PaidAt        time.Time `json:"paid_at"`
	CreatedAt     time.Time `json:"created_at"`
}

const invoiceCols = `id, COALESCE(number,''), COALESCE(customer_id::text,''), COALESCE(vendor_user_id::text,''), COALESCE(branch_id::text,''), status, currency, subtotal::text, discount_total::text, tax_total::text, total::text, COALESCE(exchange_rate_id::text,''), issued_at, voided_at, COALESCE(notes,''), pending_emission, created_at, updated_at`

const itemCols = `id, COALESCE(product_id::text,''), description, qty::text, unit_price::text, discount_pct::text, discount_amount::text, tax_rate::text, tax_amount::text, line_total::text`

const paymentCols = `id, method, COALESCE(provider,''), amount_usd::text, amount_ves::text, rate_at_payment::text, igt_amount::text, COALESCE(reference,''), status, paid_at, created_at`

func HandleCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			CustomerID string `json:"customer_id"`
			BranchID   string `json:"branch_id"`
			Currency   string `json:"currency"`
			Notes      string `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Currency == "" {
			body.Currency = "USD"
		}

		var inv Invoice
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO invoices (tenant_id, customer_id, vendor_user_id, branch_id, status, currency, notes)
				 VALUES ($1,NULLIF($2,'')::uuid,NULLIF($3,'')::uuid,NULLIF($4,'')::uuid,'draft',$5,$6)
				 RETURNING `+invoiceCols,
				mw.TenantIDFrom(r.Context()), body.CustomerID, mw.UserIDFrom(r.Context()), body.BranchID, body.Currency, body.Notes,
			).Scan(&inv.ID, &inv.Number, &inv.CustomerID, &inv.VendorUserID, &inv.BranchID, &inv.Status, &inv.Currency, &inv.Subtotal, &inv.DiscountTotal, &inv.TaxTotal, &inv.Total, &inv.ExchangeRateID, &inv.IssuedAt, &inv.VoidedAt, &inv.Notes, &inv.PendingEmission, &inv.CreatedAt, &inv.UpdatedAt)
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(inv)
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
		var invoices []Invoice
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+invoiceCols+` FROM invoices WHERE ($1='' OR status=$1) ORDER BY created_at DESC LIMIT $2 OFFSET $3`, status, perPage, offset)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var inv Invoice
				if err := rows.Scan(&inv.ID, &inv.Number, &inv.CustomerID, &inv.VendorUserID, &inv.BranchID, &inv.Status, &inv.Currency, &inv.Subtotal, &inv.DiscountTotal, &inv.TaxTotal, &inv.Total, &inv.ExchangeRateID, &inv.IssuedAt, &inv.VoidedAt, &inv.Notes, &inv.PendingEmission, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
					return err
				}
				invoices = append(invoices, inv)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if invoices == nil {
			invoices = []Invoice{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invoices)
	}
}

func HandleGet(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var inv Invoice
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`SELECT `+invoiceCols+` FROM invoices WHERE id=$1`, id,
			).Scan(&inv.ID, &inv.Number, &inv.CustomerID, &inv.VendorUserID, &inv.BranchID, &inv.Status, &inv.Currency, &inv.Subtotal, &inv.DiscountTotal, &inv.TaxTotal, &inv.Total, &inv.ExchangeRateID, &inv.IssuedAt, &inv.VoidedAt, &inv.Notes, &inv.PendingEmission, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
				return err
			}
			rows, err := tx.Query(r.Context(),
				`SELECT `+itemCols+` FROM invoice_items WHERE invoice_id=$1 ORDER BY created_at`, id)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var it Item
				if err := rows.Scan(&it.ID, &it.ProductID, &it.Description, &it.Qty, &it.UnitPrice, &it.DiscountPct, &it.DiscountAmount, &it.TaxRate, &it.TaxAmount, &it.LineTotal); err != nil {
					return err
				}
				inv.Items = append(inv.Items, it)
			}
			if err := rows.Err(); err != nil {
				return err
			}
			prows, err := tx.Query(r.Context(),
				`SELECT `+paymentCols+` FROM invoice_payments WHERE invoice_id=$1 ORDER BY paid_at DESC`, id)
			if err != nil {
				return err
			}
			defer prows.Close()
			for prows.Next() {
				var p Payment
				if err := prows.Scan(&p.ID, &p.Method, &p.Provider, &p.AmountUSD, &p.AmountVES, &p.RateAtPayment, &p.IgtAmount, &p.Reference, &p.Status, &p.PaidAt, &p.CreatedAt); err != nil {
					return err
				}
				inv.Payments = append(inv.Payments, p)
			}
			return prows.Err()
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if inv.Items == nil {
			inv.Items = []Item{}
		}
		if inv.Payments == nil {
			inv.Payments = []Payment{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(inv)
	}
}

func HandleUpdate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			CustomerID string `json:"customer_id"`
			BranchID   string `json:"branch_id"`
			Currency   string `json:"currency"`
			Notes      string `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Currency == "" {
			body.Currency = "USD"
		}

		var inv Invoice
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`UPDATE invoices SET customer_id=NULLIF($1,'')::uuid, branch_id=NULLIF($2,'')::uuid, currency=$3, notes=$4, updated_at=NOW()
				 WHERE id=$5 AND status='draft'
				 RETURNING `+invoiceCols,
				body.CustomerID, body.BranchID, body.Currency, body.Notes, id,
			).Scan(&inv.ID, &inv.Number, &inv.CustomerID, &inv.VendorUserID, &inv.BranchID, &inv.Status, &inv.Currency, &inv.Subtotal, &inv.DiscountTotal, &inv.TaxTotal, &inv.Total, &inv.ExchangeRateID, &inv.IssuedAt, &inv.VoidedAt, &inv.Notes, &inv.PendingEmission, &inv.CreatedAt, &inv.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "not found or not draft", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(inv)
	}
}

func HandleVoid(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(),
				`UPDATE invoices SET status='void', voided_at=NOW(), updated_at=NOW()
				 WHERE id=$1 AND status IN ('draft','issued','partial')`, id)
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

func HandleAddItem(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			ProductID   string `json:"product_id"`
			Description string `json:"description"`
			Qty         string `json:"qty"`
			UnitPrice   string `json:"unit_price"`
			DiscountPct string `json:"discount_pct"`
			TaxRate     string `json:"tax_rate"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Description == "" || body.Qty == "" || body.UnitPrice == "" {
			http.Error(w, "description, qty and unit_price are required", http.StatusBadRequest)
			return
		}
		if body.DiscountPct == "" {
			body.DiscountPct = "0"
		}
		if body.TaxRate == "" {
			body.TaxRate = "16"
		}

		var it Item
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			var status string
			if err := tx.QueryRow(r.Context(), `SELECT status FROM invoices WHERE id=$1`, id).Scan(&status); err != nil {
				return err
			}
			if status != "draft" {
				return errors.New("not draft")
			}
			if err := tx.QueryRow(r.Context(),
				`INSERT INTO invoice_items (tenant_id, invoice_id, product_id, description, qty, unit_price, discount_pct, discount_amount, tax_rate, tax_amount, line_total)
				 VALUES ($1,$2,NULLIF($3,'')::uuid,$4,$5,$6,$7,
				 $5*$6*$7/100,
				 $8,
				 ($5*$6 - $5*$6*$7/100)*$8/100,
				 ($5*$6 - $5*$6*$7/100) + ($5*$6 - $5*$6*$7/100)*$8/100)
				 RETURNING `+itemCols,
				mw.TenantIDFrom(r.Context()), id, body.ProductID, body.Description, body.Qty, body.UnitPrice, body.DiscountPct, body.TaxRate,
			).Scan(&it.ID, &it.ProductID, &it.Description, &it.Qty, &it.UnitPrice, &it.DiscountPct, &it.DiscountAmount, &it.TaxRate, &it.TaxAmount, &it.LineTotal); err != nil {
				return err
			}
			_, err := tx.Exec(r.Context(),
				`UPDATE invoices SET
				   subtotal=(SELECT COALESCE(SUM(qty*unit_price),0) FROM invoice_items WHERE invoice_id=$1),
				   discount_total=(SELECT COALESCE(SUM(discount_amount),0) FROM invoice_items WHERE invoice_id=$1),
				   tax_total=(SELECT COALESCE(SUM(tax_amount),0) FROM invoice_items WHERE invoice_id=$1),
				   total=(SELECT COALESCE(SUM(line_total),0) FROM invoice_items WHERE invoice_id=$1),
				   updated_at=NOW()
				 WHERE id=$1`, id)
			return err
		})
		if err != nil {
			http.Error(w, "not found or not draft", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(it)
	}
}

func HandleRemoveItem(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		itemID := chi.URLParam(r, "itemId")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			var status string
			if err := tx.QueryRow(r.Context(), `SELECT status FROM invoices WHERE id=$1`, id).Scan(&status); err != nil {
				return err
			}
			if status != "draft" {
				return errors.New("not draft")
			}
			ct, err := tx.Exec(r.Context(), `DELETE FROM invoice_items WHERE id=$1 AND invoice_id=$2`, itemID, id)
			if err != nil {
				return err
			}
			if ct.RowsAffected() == 0 {
				return errNotFound
			}
			_, err = tx.Exec(r.Context(),
				`UPDATE invoices SET
				   subtotal=(SELECT COALESCE(SUM(qty*unit_price),0) FROM invoice_items WHERE invoice_id=$1),
				   discount_total=(SELECT COALESCE(SUM(discount_amount),0) FROM invoice_items WHERE invoice_id=$1),
				   tax_total=(SELECT COALESCE(SUM(tax_amount),0) FROM invoice_items WHERE invoice_id=$1),
				   total=(SELECT COALESCE(SUM(line_total),0) FROM invoice_items WHERE invoice_id=$1),
				   updated_at=NOW()
				 WHERE id=$1`, id)
			return err
		})
		if err != nil {
			http.Error(w, "not found or not draft", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleEmit(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			ExchangeRateID string `json:"exchange_rate_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		var inv Invoice
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			var status, branchID string
			if err := tx.QueryRow(r.Context(), `SELECT status, COALESCE(branch_id::text,'') FROM invoices WHERE id=$1`, id).Scan(&status, &branchID); err != nil {
				return err
			}
			if status != "draft" {
				return errors.New("not draft")
			}
			prefix := "01"
			today := time.Now().UTC().Format("2006-01-02")
			var seq int
			if err := tx.QueryRow(r.Context(),
				`INSERT INTO invoice_sequences (tenant_id, branch_id, prefix, last_seq, last_date)
				 VALUES ($1,NULLIF($2,'')::uuid,$3,1,$4::date)
				 ON CONFLICT (tenant_id, branch_id, prefix, last_date)
				 DO UPDATE SET last_seq = invoice_sequences.last_seq + 1
				 RETURNING last_seq`,
				mw.TenantIDFrom(r.Context()), branchID, prefix, today,
			).Scan(&seq); err != nil {
				return err
			}
			datePart := time.Now().UTC().Format("20060102")
			number := formatInvoiceNumber(prefix, datePart, seq)

			return tx.QueryRow(r.Context(),
				`UPDATE invoices SET number=$1, status='issued', issued_at=NOW(), exchange_rate_id=NULLIF($2,'')::uuid, updated_at=NOW()
				 WHERE id=$3 AND status='draft'
				 RETURNING `+invoiceCols,
				number, body.ExchangeRateID, id,
			).Scan(&inv.ID, &inv.Number, &inv.CustomerID, &inv.VendorUserID, &inv.BranchID, &inv.Status, &inv.Currency, &inv.Subtotal, &inv.DiscountTotal, &inv.TaxTotal, &inv.Total, &inv.ExchangeRateID, &inv.IssuedAt, &inv.VoidedAt, &inv.Notes, &inv.PendingEmission, &inv.CreatedAt, &inv.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "not found or not draft", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(inv)
	}
}

func HandleAddPayment(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			Method        string `json:"method"`
			Provider      string `json:"provider"`
			AmountUSD     string `json:"amount_usd"`
			AmountVES     string `json:"amount_ves"`
			RateAtPayment string `json:"rate_at_payment"`
			IgtAmount     string `json:"igt_amount"`
			Reference     string `json:"reference"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Method == "" || body.AmountUSD == "" {
			http.Error(w, "method and amount_usd are required", http.StatusBadRequest)
			return
		}
		if body.AmountVES == "" {
			body.AmountVES = "0"
		}
		if body.RateAtPayment == "" {
			body.RateAtPayment = "1"
		}
		if body.IgtAmount == "" {
			body.IgtAmount = "0"
		}

		var p Payment
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			var status string
			if err := tx.QueryRow(r.Context(), `SELECT status FROM invoices WHERE id=$1`, id).Scan(&status); err != nil {
				return err
			}
			if status == "void" {
				return errors.New("void invoice")
			}
			if err := tx.QueryRow(r.Context(),
				`INSERT INTO invoice_payments (tenant_id, invoice_id, method, provider, amount_usd, amount_ves, rate_at_payment, igt_amount, reference)
				 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
				 RETURNING `+paymentCols,
				mw.TenantIDFrom(r.Context()), id, body.Method, body.Provider, body.AmountUSD, body.AmountVES, body.RateAtPayment, body.IgtAmount, body.Reference,
			).Scan(&p.ID, &p.Method, &p.Provider, &p.AmountUSD, &p.AmountVES, &p.RateAtPayment, &p.IgtAmount, &p.Reference, &p.Status, &p.PaidAt, &p.CreatedAt); err != nil {
				return err
			}
			_, err := tx.Exec(r.Context(),
				`UPDATE invoices SET
				   status = CASE WHEN (SELECT COALESCE(SUM(amount_usd),0) FROM invoice_payments WHERE invoice_id=$1) >= total THEN 'paid' ELSE 'partial' END,
				   updated_at=NOW()
				 WHERE id=$1`, id)
			return err
		})
		if err != nil {
			http.Error(w, "not found or void", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}

func formatInvoiceNumber(prefix, datePart string, seq int) string {
	return fmt.Sprintf("%s-%s-%05d", prefix, datePart, seq)
}
