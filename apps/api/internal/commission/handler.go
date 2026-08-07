package commission

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

type Commission struct {
	ID               string  `json:"id"`
	TenantID         string  `json:"tenant_id,omitempty"`
	VendorUserID     string  `json:"vendor_user_id"`
	InvoiceID        string  `json:"invoice_id"`
	CategoryID       *string `json:"category_id,omitempty"`
	RatePct          float64 `json:"rate_pct"`
	BaseAmount       float64 `json:"base_amount"`
	CommissionAmount float64 `json:"commission_amount"`
	Period           string  `json:"period"`
}

func HandleListCommissions(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vendorID := r.URL.Query().Get("vendor_user_id")
		period := r.URL.Query().Get("period")
		var cs []Commission
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT id, tenant_id, vendor_user_id, invoice_id, category_id, rate_pct, base_amount, commission_amount, period
				 FROM sales_commissions WHERE ($1 = '' OR vendor_user_id = $1) AND ($2 = '' OR period = $2)
				 ORDER BY period DESC`, vendorID, period)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var c Commission
				if err := rows.Scan(&c.ID, &c.TenantID, &c.VendorUserID, &c.InvoiceID, &c.CategoryID,
					&c.RatePct, &c.BaseAmount, &c.CommissionAmount, &c.Period); err != nil {
					return err
				}
				cs = append(cs, c)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if cs == nil {
			cs = []Commission{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cs)
	}
}

func HandleCalculateCommissions(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Period string `json:"period"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Period == "" {
			http.Error(w, "period required", http.StatusBadRequest)
			return
		}
		tid := mw.TenantIDFrom(r.Context())
		var created []Commission
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT i.id, i.vendor_user_id, i.subtotal, COALESCE(vc.rate_pct, 0), COALESCE(vc.category_id, NULL)
				 FROM invoices i
				 LEFT JOIN vendor_commission_rates vc ON vc.vendor_user_id = i.vendor_user_id
				 WHERE i.status = 'paid' AND to_char(i.paid_at, 'YYYY-MM') = $1`,
				body.Period)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var invoiceID, vendorUserID string
				var baseAmount, ratePct float64
				var categoryID *string
				if err := rows.Scan(&invoiceID, &vendorUserID, &baseAmount, &ratePct, &categoryID); err != nil {
					return err
				}
				commAmount := baseAmount * ratePct / 100
				var c Commission
				if err := tx.QueryRow(r.Context(),
					`INSERT INTO sales_commissions (tenant_id, vendor_user_id, invoice_id, category_id, rate_pct, base_amount, commission_amount, period)
					 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
					 ON CONFLICT (tenant_id, invoice_id) DO UPDATE SET base_amount=EXCLUDED.base_amount, commission_amount=EXCLUDED.commission_amount
					 RETURNING id, tenant_id, vendor_user_id, invoice_id, category_id, rate_pct, base_amount, commission_amount, period`,
					tid, vendorUserID, invoiceID, categoryID, ratePct, baseAmount, commAmount, body.Period,
				).Scan(&c.ID, &c.TenantID, &c.VendorUserID, &c.InvoiceID, &c.CategoryID,
					&c.RatePct, &c.BaseAmount, &c.CommissionAmount, &c.Period); err != nil {
					return err
				}
				created = append(created, c)
			}
			return nil
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		if created == nil {
			created = []Commission{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(created)
	}
}
