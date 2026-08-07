package payment

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

type PaymentLinkRow struct {
	ID          string     `json:"id"`
	InvoiceID   string     `json:"invoice_id"`
	Provider    string     `json:"provider"`
	ShortCode   string     `json:"short_code"`
	AmountUSD   float64    `json:"amount_usd"`
	Status      string     `json:"status"`
	ProviderRef string     `json:"provider_ref,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func newShortCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func HandleCreatePaymentLink(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			InvoiceID string  `json:"invoice_id"`
			Provider  string  `json:"provider"`
			AmountUSD float64 `json:"amount_usd"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.InvoiceID == "" || body.Provider == "" || body.AmountUSD <= 0 {
			http.Error(w, "invoice_id, provider and positive amount required", http.StatusBadRequest)
			return
		}
		provider := GetProvider(body.Provider)
		if provider == nil {
			http.Error(w, "unknown provider", http.StatusBadRequest)
			return
		}
		link, err := provider.CreatePaymentLink(body.InvoiceID, body.AmountUSD)
		if err != nil {
			http.Error(w, "provider error", http.StatusBadGateway)
			return
		}
		shortCode := newShortCode()
		expiresAt := time.Now().Add(24 * time.Hour)
		var row PaymentLinkRow
		err = mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO payment_links (invoice_id, provider, short_code, amount_usd, status, provider_ref, expires_at)
				 VALUES ($1,$2,$3,$4,$5,$6,$7)
				 RETURNING id, invoice_id, provider, short_code, amount_usd, status, provider_ref, expires_at, paid_at, created_at`,
				body.InvoiceID, body.Provider, shortCode, body.AmountUSD, link.Status, link.ProviderRef, expiresAt,
			).Scan(&row.ID, &row.InvoiceID, &row.Provider, &row.ShortCode, &row.AmountUSD, &row.Status, &row.ProviderRef, &row.ExpiresAt, &row.PaidAt, &row.CreatedAt)
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(row)
	}
}

func HandleGetPaymentLink(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		if code == "" {
			http.Error(w, "code required", http.StatusBadRequest)
			return
		}
		var row PaymentLinkRow
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`SELECT id, invoice_id, provider, short_code, amount_usd, status, provider_ref, expires_at, paid_at, created_at
				 FROM payment_links WHERE short_code = $1`,
				code,
			).Scan(&row.ID, &row.InvoiceID, &row.Provider, &row.ShortCode, &row.AmountUSD, &row.Status, &row.ProviderRef, &row.ExpiresAt, &row.PaidAt, &row.CreatedAt)
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(row)
	}
}

func HandleListPaymentLinks(pool *pgxpool.Pool) http.HandlerFunc {
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
		q := `SELECT id, invoice_id, provider, short_code, amount_usd, status, provider_ref, expires_at, paid_at, created_at
		      FROM payment_links`
		args := []interface{}{}
		if status != "" {
			q += ` WHERE status = $1`
			args = append(args, status)
		}
		q += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
		args = append(args, perPage, offset)
		var items []PaymentLinkRow
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(), q, args...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var row PaymentLinkRow
				if err := rows.Scan(&row.ID, &row.InvoiceID, &row.Provider, &row.ShortCode, &row.AmountUSD, &row.Status, &row.ProviderRef, &row.ExpiresAt, &row.PaidAt, &row.CreatedAt); err != nil {
					return err
				}
				items = append(items, row)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if items == nil {
			items = []PaymentLinkRow{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func HandleWebhook(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerName := chi.URLParam(r, "provider")
		provider := GetProvider(providerName)
		if provider == nil {
			http.Error(w, "unknown provider", http.StatusBadRequest)
			return
		}
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read error", http.StatusBadRequest)
			return
		}
		signature := r.Header.Get("X-Signature")
		event, err := provider.HandleWebhook(payload, signature)
		if err != nil {
			http.Error(w, "webhook error", http.StatusBadGateway)
			return
		}
		// Resolve tenant_id from payment_links by provider_ref (webhooks lack JWT context).
		var tenantID string
		err = pool.QueryRow(r.Context(),
			`SELECT tenant_id::text FROM payment_links WHERE provider_ref = $1 LIMIT 1`,
			event.ProviderRef,
		).Scan(&tenantID)
		if err != nil && err != pgx.ErrNoRows {
			log.Printf("db error resolving tenant for webhook: %v", err)
		}
		if tenantID == "" {
			log.Printf("warning: no tenant found for payment webhook provider_ref=%s", event.ProviderRef)
		}

		tx, err := pool.Begin(r.Context())
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(r.Context())

		// Set app.tenant_id so RLS applies for this transaction (same pattern as WithTenantTx).
		if tenantID != "" {
			if _, err := tx.Exec(r.Context(), "SELECT set_config('app.tenant_id', $1, true)", tenantID); err != nil {
				log.Printf("db error: %v", err)
				http.Error(w, "db error", http.StatusInternalServerError)
				return
			}
		}

		_, err = tx.Exec(r.Context(),
			`INSERT INTO payment_webhook_events (provider, event_type, payload, signature, processed)
			 VALUES ($1,$2,$3,$4,true)`,
			providerName, event.Status, payload, signature,
		)
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		// update matching payment link
		if event.Status == "paid" {
			_, err = tx.Exec(r.Context(),
				`UPDATE payment_links SET status = 'paid', paid_at = NOW()
				 WHERE provider_ref = $1`,
				event.ProviderRef,
			)
		} else {
			_, err = tx.Exec(r.Context(),
				`UPDATE payment_links SET status = $1 WHERE provider_ref = $2`,
				event.Status, event.ProviderRef,
			)
		}
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if err := tx.Commit(r.Context()); err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
