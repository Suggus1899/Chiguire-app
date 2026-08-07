package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

type Webhook struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id,omitempty"`
	URL      string `json:"url"`
	Events   string `json:"events"`
	Secret   string `json:"-"`
	IsActive bool   `json:"is_active"`
}

type Delivery struct {
	ID           string  `json:"id"`
	TenantID     string  `json:"tenant_id,omitempty"`
	WebhookID    string  `json:"webhook_id"`
	EventType    string  `json:"event_type"`
	Payload      []byte  `json:"payload"`
	Status       string  `json:"status"`
	Attempts     int     `json:"attempts"`
	LastResponse *string `json:"last_response,omitempty"`
}

func HandleCreateWebhook(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			URL    string `json:"url"`
			Events string `json:"events"`
			Secret string `json:"secret"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.URL == "" {
			http.Error(w, "url required", http.StatusBadRequest)
			return
		}
		var wh Webhook
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO outgoing_webhooks (tenant_id, url, events, secret, is_active)
				 VALUES ($1,$2,$3,$4,true)
				 RETURNING id, tenant_id, url, events, is_active`,
				mw.TenantIDFrom(r.Context()), body.URL, body.Events, body.Secret,
			).Scan(&wh.ID, &wh.TenantID, &wh.URL, &wh.Events, &wh.IsActive)
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(wh)
	}
}

func HandleListWebhooks(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var whs []Webhook
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT id, tenant_id, url, events, is_active FROM outgoing_webhooks ORDER BY id`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var wh Webhook
				if err := rows.Scan(&wh.ID, &wh.TenantID, &wh.URL, &wh.Events, &wh.IsActive); err != nil {
					return err
				}
				whs = append(whs, wh)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if whs == nil {
			whs = []Webhook{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(whs)
	}
}

func HandleDeleteWebhook(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			_, err := tx.Exec(r.Context(), `DELETE FROM outgoing_webhooks WHERE id=$1`, id)
			return err
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleTestWebhook(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var wh Webhook
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`SELECT id, tenant_id, url, events, secret, is_active FROM outgoing_webhooks WHERE id=$1`, id,
			).Scan(&wh.ID, &wh.TenantID, &wh.URL, &wh.Events, &wh.Secret, &wh.IsActive)
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		payload, _ := json.Marshal(map[string]any{"event": "test", "webhook_id": id})
		resp, err := deliver(wh.URL, wh.Secret, "test", payload)
		status := "success"
		if err != nil {
			status = "failed"
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": status, "response": resp})
	}
}

func HandleListDeliveries(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ds []Delivery
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT id, tenant_id, webhook_id, event_type, payload, status, attempts, last_response
				 FROM webhook_deliveries ORDER BY id DESC LIMIT 100`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var d Delivery
				if err := rows.Scan(&d.ID, &d.TenantID, &d.WebhookID, &d.EventType, &d.Payload,
					&d.Status, &d.Attempts, &d.LastResponse); err != nil {
					return err
				}
				ds = append(ds, d)
			}
			return rows.Err()
		})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if ds == nil {
			ds = []Delivery{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ds)
	}
}
