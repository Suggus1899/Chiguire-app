package paymentmethod

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

type PaymentMethod struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Type                 string    `json:"type"`
	Currency             string    `json:"currency"`
	AvailableForInvoices bool      `json:"available_for_invoices"`
	AvailableForChange   bool      `json:"available_for_change"`
	AvailableForRefunds  bool      `json:"available_for_refunds"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

const pmCols = `id, name, type, currency, available_for_invoices, available_for_change, available_for_refunds, is_active, created_at, updated_at`

func HandleCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name                 string `json:"name"`
			Type                 string `json:"type"`
			Currency             string `json:"currency"`
			AvailableForInvoices *bool  `json:"available_for_invoices"`
			AvailableForChange   *bool  `json:"available_for_change"`
			AvailableForRefunds  *bool  `json:"available_for_refunds"`
			IsActive             *bool  `json:"is_active"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		if body.Type == "" {
			body.Type = "cash"
		}
		if body.Currency == "" {
			body.Currency = "USD"
		}
		availInv := true
		availChange := false
		availRefund := false
		isActive := true
		if body.AvailableForInvoices != nil {
			availInv = *body.AvailableForInvoices
		}
		if body.AvailableForChange != nil {
			availChange = *body.AvailableForChange
		}
		if body.AvailableForRefunds != nil {
			availRefund = *body.AvailableForRefunds
		}
		if body.IsActive != nil {
			isActive = *body.IsActive
		}

		var pm PaymentMethod
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO payment_methods (tenant_id, name, type, currency, available_for_invoices, available_for_change, available_for_refunds, is_active)
				 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
				 RETURNING `+pmCols,
				mw.TenantIDFrom(r.Context()), body.Name, body.Type, body.Currency, availInv, availChange, availRefund, isActive,
			).Scan(&pm.ID, &pm.Name, &pm.Type, &pm.Currency, &pm.AvailableForInvoices, &pm.AvailableForChange, &pm.AvailableForRefunds, &pm.IsActive, &pm.CreatedAt, &pm.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(pm)
	}
}

func HandleList(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pms []PaymentMethod
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+pmCols+` FROM payment_methods ORDER BY created_at DESC`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var pm PaymentMethod
				if err := rows.Scan(&pm.ID, &pm.Name, &pm.Type, &pm.Currency, &pm.AvailableForInvoices, &pm.AvailableForChange, &pm.AvailableForRefunds, &pm.IsActive, &pm.CreatedAt, &pm.UpdatedAt); err != nil {
					return err
				}
				pms = append(pms, pm)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if pms == nil {
			pms = []PaymentMethod{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pms)
	}
}

func HandleUpdate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			Name                 string `json:"name"`
			Type                 string `json:"type"`
			Currency             string `json:"currency"`
			AvailableForInvoices *bool  `json:"available_for_invoices"`
			AvailableForChange   *bool  `json:"available_for_change"`
			AvailableForRefunds  *bool  `json:"available_for_refunds"`
			IsActive             *bool  `json:"is_active"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Type == "" {
			body.Type = "cash"
		}
		if body.Currency == "" {
			body.Currency = "USD"
		}
		availInv := true
		availChange := false
		availRefund := false
		isActive := true
		if body.AvailableForInvoices != nil {
			availInv = *body.AvailableForInvoices
		}
		if body.AvailableForChange != nil {
			availChange = *body.AvailableForChange
		}
		if body.AvailableForRefunds != nil {
			availRefund = *body.AvailableForRefunds
		}
		if body.IsActive != nil {
			isActive = *body.IsActive
		}

		var pm PaymentMethod
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`UPDATE payment_methods SET name=$1, type=$2, currency=$3, available_for_invoices=$4, available_for_change=$5, available_for_refunds=$6, is_active=$7, updated_at=NOW()
				 WHERE id=$8
				 RETURNING `+pmCols,
				body.Name, body.Type, body.Currency, availInv, availChange, availRefund, isActive, id,
			).Scan(&pm.ID, &pm.Name, &pm.Type, &pm.Currency, &pm.AvailableForInvoices, &pm.AvailableForChange, &pm.AvailableForRefunds, &pm.IsActive, &pm.CreatedAt, &pm.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pm)
	}
}

func HandleDelete(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(), `DELETE FROM payment_methods WHERE id=$1`, id)
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
