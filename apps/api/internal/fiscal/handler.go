package fiscal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaxCategory struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	IVARate     float64 `json:"iva_rate"`
	Description string  `json:"description,omitempty"`
}

func HandleListTaxCategories(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var items []TaxCategory
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT id, name, iva_rate, description FROM tax_categories ORDER BY iva_rate`,
			)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var c TaxCategory
				if err := rows.Scan(&c.ID, &c.Name, &c.IVARate, &c.Description); err != nil {
					return err
				}
				items = append(items, c)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if items == nil {
			items = []TaxCategory{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func HandleCreateTaxCategory(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name        string  `json:"name"`
			IVARate     float64 `json:"iva_rate"`
			Description string  `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, "name required", http.StatusBadRequest)
			return
		}
		var c TaxCategory
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO tax_categories (name, iva_rate, description)
				 VALUES ($1,$2,$3) RETURNING id, name, iva_rate, description`,
				body.Name, body.IVARate, body.Description,
			).Scan(&c.ID, &c.Name, &c.IVARate, &c.Description)
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(c)
	}
}

type Withholding struct {
	ID             string  `json:"id"`
	InvoiceID      *string `json:"invoice_id,omitempty"`
	VendorID       *string `json:"vendor_id,omitempty"`
	CustomerID     *string `json:"customer_id,omitempty"`
	Type           string  `json:"type"`
	Base           float64 `json:"base"`
	Rate           float64 `json:"rate"`
	Amount         float64 `json:"amount"`
	DocumentNumber string  `json:"document_number"`
	Period         string  `json:"period"`
}

func HandleCreateWithholding(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			InvoiceID      *string `json:"invoice_id"`
			VendorID       *string `json:"vendor_id"`
			CustomerID     *string `json:"customer_id"`
			Type           string  `json:"type"`
			Base           float64 `json:"base"`
			Rate           float64 `json:"rate"`
			Amount         float64 `json:"amount"`
			DocumentNumber string  `json:"document_number"`
			Period         string  `json:"period"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Type == "" || body.Period == "" {
			http.Error(w, "type and period required", http.StatusBadRequest)
			return
		}
		var wh Withholding
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO tax_withholdings (invoice_id, vendor_id, customer_id, type, base, rate, amount, document_number, period)
				 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
				 RETURNING id, invoice_id, vendor_id, customer_id, type, base, rate, amount, document_number, period`,
				body.InvoiceID, body.VendorID, body.CustomerID, body.Type, body.Base, body.Rate, body.Amount, body.DocumentNumber, body.Period,
			).Scan(&wh.ID, &wh.InvoiceID, &wh.VendorID, &wh.CustomerID, &wh.Type, &wh.Base, &wh.Rate, &wh.Amount, &wh.DocumentNumber, &wh.Period)
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(wh)
	}
}

func HandleListWithholdings(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		period := r.URL.Query().Get("period")
		q := `SELECT id, invoice_id, vendor_id, customer_id, type, base, rate, amount, document_number, period
		      FROM tax_withholdings`
		args := []interface{}{}
		if period != "" {
			q += ` WHERE period = $1`
			args = append(args, period)
		}
		q += ` ORDER BY period DESC, created_at DESC`
		var items []Withholding
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(), q, args...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var wh Withholding
				if err := rows.Scan(&wh.ID, &wh.InvoiceID, &wh.VendorID, &wh.CustomerID, &wh.Type, &wh.Base, &wh.Rate, &wh.Amount, &wh.DocumentNumber, &wh.Period); err != nil {
					return err
				}
				items = append(items, wh)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if items == nil {
			items = []Withholding{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

type FiscalPeriod struct {
	Period   string     `json:"period"`
	IsClosed bool       `json:"is_closed"`
	ClosedAt *time.Time `json:"closed_at,omitempty"`
}

func HandleListFiscalPeriods(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var items []FiscalPeriod
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT period, is_closed, closed_at FROM fiscal_periods ORDER BY period DESC`,
			)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var p FiscalPeriod
				if err := rows.Scan(&p.Period, &p.IsClosed, &p.ClosedAt); err != nil {
					return err
				}
				items = append(items, p)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if items == nil {
			items = []FiscalPeriod{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func HandleClosePeriod(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		period := chi.URLParam(r, "period")
		if period == "" {
			http.Error(w, "period required", http.StatusBadRequest)
			return
		}
		var p FiscalPeriod
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`UPDATE fiscal_periods SET is_closed = true, closed_at = NOW()
				 WHERE period = $1 AND is_closed = false
				 RETURNING period, is_closed, closed_at`,
				period,
			).Scan(&p.Period, &p.IsClosed, &p.ClosedAt)
		})
		if err != nil {
			http.Error(w, "period not found or already closed", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

type FiscalBook struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Period      string `json:"period"`
	EntriesJSON string `json:"entries_json"`
}

func HandleGenerateFiscalBook(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Type   string `json:"type"`
			Period string `json:"period"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Type == "" || body.Period == "" {
			http.Error(w, "type and period required", http.StatusBadRequest)
			return
		}
		var book FiscalBook
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			// Build entries from invoices for the period
			rows, err := tx.Query(r.Context(),
				`SELECT json_agg(row_to_json(t)) FROM (
					SELECT id, number, date, customer_id, vendor_id, subtotal, tax, total, currency
					FROM invoices WHERE period = $1 AND type = $2 ORDER BY date
				 ) t`,
				body.Period, body.Type,
			)
			if err != nil {
				return err
			}
			defer rows.Close()
			var entriesJSON []byte
			if rows.Next() {
				if err := rows.Scan(&entriesJSON); err != nil {
					return err
				}
			}
			if entriesJSON == nil {
				entriesJSON = []byte("[]")
			}
			return tx.QueryRow(r.Context(),
				`INSERT INTO fiscal_books (type, period, entries_json)
				 VALUES ($1,$2,$3) RETURNING id, type, period, entries_json`,
				body.Type, body.Period, string(entriesJSON),
			).Scan(&book.ID, &book.Type, &book.Period, &book.EntriesJSON)
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(book)
	}
}

func HandleGetFiscalBook(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}
		var book FiscalBook
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`SELECT id, type, period, entries_json FROM fiscal_books WHERE id = $1`,
				id,
			).Scan(&book.ID, &book.Type, &book.Period, &book.EntriesJSON)
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	}
}

type ExchangeRate struct {
	Currency    string    `json:"currency"`
	RateToVES   float64   `json:"rate_to_ves"`
	Source      string    `json:"source"`
	EffectiveAt time.Time `json:"effective_at"`
}

func HandleGetExchangeRate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currency := r.URL.Query().Get("currency")
		if currency == "" {
			currency = "USD"
		}
		var rate ExchangeRate
		err := pool.QueryRow(r.Context(),
			`SELECT currency, rate_to_ves, source, effective_at
			 FROM exchange_rates WHERE currency = $1 ORDER BY effective_at DESC LIMIT 1`,
			currency,
		).Scan(&rate.Currency, &rate.RateToVES, &rate.Source, &rate.EffectiveAt)
		if err != nil {
			http.Error(w, "no rate found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rate)
	}
}

func HandleSetExchangeRate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Currency  string  `json:"currency"`
			RateToVES float64 `json:"rate_to_ves"`
			Source    string  `json:"source"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Currency == "" || body.RateToVES <= 0 {
			http.Error(w, "currency and positive rate required", http.StatusBadRequest)
			return
		}
		if body.Source == "" {
			body.Source = "manual"
		}
		var rate ExchangeRate
		err := pool.QueryRow(r.Context(),
			`INSERT INTO exchange_rates (currency, rate_to_ves, source, effective_at)
			 VALUES ($1,$2,$3,NOW()) RETURNING currency, rate_to_ves, source, effective_at`,
			body.Currency, body.RateToVES, body.Source,
		).Scan(&rate.Currency, &rate.RateToVES, &rate.Source, &rate.EffectiveAt)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(rate)
	}
}

func HandleListExchangeRates(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currency := r.URL.Query().Get("currency")
		q := `SELECT currency, rate_to_ves, source, effective_at FROM exchange_rates`
		args := []interface{}{}
		if currency != "" {
			q += ` WHERE currency = $1`
			args = append(args, currency)
		}
		q += ` ORDER BY effective_at DESC LIMIT 100`
		rows, err := pool.Query(r.Context(), q, args...)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var items []ExchangeRate
		for rows.Next() {
			var rate ExchangeRate
			if err := rows.Scan(&rate.Currency, &rate.RateToVES, &rate.Source, &rate.EffectiveAt); err != nil {
				http.Error(w, "scan error", http.StatusInternalServerError)
				return
			}
			items = append(items, rate)
		}
		if items == nil {
			items = []ExchangeRate{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}
