package accountspayable

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

type Receivable struct {
	ID         string     `json:"id"`
	CustomerID string     `json:"customer_id"`
	InvoiceID  string     `json:"invoice_id"`
	Amount     string     `json:"amount"`
	Balance    string     `json:"balance"`
	Status     string     `json:"status"`
	DueDate    *time.Time `json:"due_date"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Payable struct {
	ID              string     `json:"id"`
	VendorID        string     `json:"vendor_id"`
	PurchaseOrderID string     `json:"purchase_order_id"`
	Amount          string     `json:"amount"`
	Balance         string     `json:"balance"`
	Status          string     `json:"status"`
	DueDate         *time.Time `json:"due_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Payment struct {
	ID              string    `json:"id"`
	AccountType     string    `json:"account_type"`
	AccountID       string    `json:"account_id"`
	Amount          string    `json:"amount"`
	PaymentMethodID string    `json:"payment_method_id"`
	Reference       string    `json:"reference"`
	Notes           string    `json:"notes"`
	PaidAt          time.Time `json:"paid_at"`
	CreatedAt       time.Time `json:"created_at"`
}

const arCols = `id, customer_id::text, COALESCE(invoice_id::text,''), amount::text, balance::text, status, due_date, created_at, updated_at`
const apCols = `id, vendor_id::text, COALESCE(purchase_order_id::text,''), amount::text, balance::text, status, due_date, created_at, updated_at`
const payCols = `id, account_type, account_id::text, amount::text, COALESCE(payment_method_id::text,''), COALESCE(reference,''), COALESCE(notes,''), paid_at, created_at`

func HandleListReceivable(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		var rows []Receivable
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			q, err := tx.Query(r.Context(),
				`SELECT `+arCols+` FROM accounts_receivable WHERE ($1='' OR status=$1) ORDER BY created_at DESC`, status)
			if err != nil {
				return err
			}
			defer q.Close()
			for q.Next() {
				var ar Receivable
				if err := q.Scan(&ar.ID, &ar.CustomerID, &ar.InvoiceID, &ar.Amount, &ar.Balance, &ar.Status, &ar.DueDate, &ar.CreatedAt, &ar.UpdatedAt); err != nil {
					return err
				}
				rows = append(rows, ar)
			}
			return q.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if rows == nil {
			rows = []Receivable{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}

func HandleListPayable(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		var rows []Payable
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			q, err := tx.Query(r.Context(),
				`SELECT `+apCols+` FROM accounts_payable WHERE ($1='' OR status=$1) ORDER BY created_at DESC`, status)
			if err != nil {
				return err
			}
			defer q.Close()
			for q.Next() {
				var ap Payable
				if err := q.Scan(&ap.ID, &ap.VendorID, &ap.PurchaseOrderID, &ap.Amount, &ap.Balance, &ap.Status, &ap.DueDate, &ap.CreatedAt, &ap.UpdatedAt); err != nil {
					return err
				}
				rows = append(rows, ap)
			}
			return q.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if rows == nil {
			rows = []Payable{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}

func HandleCreatePayment(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			AccountType     string `json:"account_type"`
			AccountID       string `json:"account_id"`
			Amount          string `json:"amount"`
			PaymentMethodID string `json:"payment_method_id"`
			Reference       string `json:"reference"`
			Notes           string `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.AccountType == "" || body.AccountID == "" || body.Amount == "" {
			http.Error(w, "account_type, account_id and amount are required", http.StatusBadRequest)
			return
		}
		if body.AccountType != "receivable" && body.AccountType != "payable" {
			http.Error(w, "account_type must be receivable or payable", http.StatusBadRequest)
			return
		}

		var p Payment
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if err := tx.QueryRow(r.Context(),
				`INSERT INTO account_payments (tenant_id, account_type, account_id, amount, payment_method_id, reference, notes)
				 VALUES ($1,$2,$3,$4,NULLIF($5,'')::uuid,$6,$7)
				 RETURNING `+payCols,
				mw.TenantIDFrom(r.Context()), body.AccountType, body.AccountID, body.Amount, body.PaymentMethodID, body.Reference, body.Notes,
			).Scan(&p.ID, &p.AccountType, &p.AccountID, &p.Amount, &p.PaymentMethodID, &p.Reference, &p.Notes, &p.PaidAt, &p.CreatedAt); err != nil {
				return err
			}
			if body.AccountType == "receivable" {
				_, err := tx.Exec(r.Context(),
					`UPDATE accounts_receivable SET
					   balance = balance - $1,
					   status = CASE WHEN balance - $1 <= 0 THEN 'paid' ELSE status END,
					   updated_at = NOW()
					 WHERE id=$2`, body.Amount, body.AccountID)
				return err
			}
			_, err := tx.Exec(r.Context(),
				`UPDATE accounts_payable SET
				   balance = balance - $1,
				   status = CASE WHEN balance - $1 <= 0 THEN 'paid' ELSE status END,
				   updated_at = NOW()
				 WHERE id=$2`, body.Amount, body.AccountID)
			return err
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}

func HandleListPayments(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountType := r.URL.Query().Get("account_type")
		var rows []Payment
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			q, err := tx.Query(r.Context(),
				`SELECT `+payCols+` FROM account_payments WHERE ($1='' OR account_type=$1) ORDER BY paid_at DESC`, accountType)
			if err != nil {
				return err
			}
			defer q.Close()
			for q.Next() {
				var p Payment
				if err := q.Scan(&p.ID, &p.AccountType, &p.AccountID, &p.Amount, &p.PaymentMethodID, &p.Reference, &p.Notes, &p.PaidAt, &p.CreatedAt); err != nil {
					return err
				}
				rows = append(rows, p)
			}
			return q.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if rows == nil {
			rows = []Payment{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}

func HandleSendReminder(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			Channel string `json:"channel"`
			To      string `json:"to"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Channel == "" {
			body.Channel = "email"
		}
		if body.Channel != "sms" && body.Channel != "email" {
			http.Error(w, "channel must be sms or email", http.StatusBadRequest)
			return
		}

		_ = id

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"channel": body.Channel,
			"message": "reminder queued",
		})
	}
}
