package saas

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Subscription struct {
	ID                    string     `json:"id"`
	TenantID              string     `json:"tenant_id,omitempty"`
	Plan                  string     `json:"plan"`
	Status                string     `json:"status"`
	StripeCustomerID      *string    `json:"stripe_customer_id,omitempty"`
	StripeSubscriptionID  *string    `json:"stripe_subscription_id,omitempty"`
	CurrentPeriodStart    *time.Time `json:"current_period_start,omitempty"`
	CurrentPeriodEnd      *time.Time `json:"current_period_end,omitempty"`
	CancelAtPeriodEnd     bool       `json:"cancel_at_period_end"`
}

type PlanLimit struct {
	Plan             string `json:"plan"`
	MaxInvoicesMonth int    `json:"max_invoices_month"`
	MaxBranches      int    `json:"max_branches"`
	MaxUsers         int    `json:"max_users"`
	Features         []byte `json:"features"`
}

func HandleGetSubscription(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sub Subscription
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`SELECT id, tenant_id, plan, status, stripe_customer_id, stripe_subscription_id,
				 	current_period_start, current_period_end, cancel_at_period_end
				 FROM subscriptions LIMIT 1`,
			).Scan(&sub.ID, &sub.TenantID, &sub.Plan, &sub.Status, &sub.StripeCustomerID,
				&sub.StripeSubscriptionID, &sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelAtPeriodEnd)
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sub)
	}
}

func HandleUpdatePlan(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Plan string `json:"plan"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Plan == "" {
			http.Error(w, "plan required", http.StatusBadRequest)
			return
		}
		err := withTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			_, err := tx.Exec(r.Context(), `UPDATE subscriptions SET plan=$1`, body.Plan)
			return err
		})
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleGetPlanLimits(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		plan := r.URL.Query().Get("plan")
		var pl PlanLimit
		err := pool.QueryRow(r.Context(),
			`SELECT plan, max_invoices_month, max_branches, max_users, features
			 FROM plan_limits WHERE plan=$1`, plan,
		).Scan(&pl.Plan, &pl.MaxInvoicesMonth, &pl.MaxBranches, &pl.MaxUsers, &pl.Features)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pl)
	}
}

func HandleCreateStripeCheckoutSession(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Plan string `json:"plan"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"checkout_url": "https://checkout.stripe.com/placeholder",
			"session_id":   "cs_placeholder_" + body.Plan,
		})
	}
}

func HandleStripeWebhook(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event map[string]any
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		eventType, _ := event["type"].(string)
		switch eventType {
		case "customer.subscription.created", "customer.subscription.updated":
			if data, ok := event["data"].(map[string]any); ok {
				if obj, ok := data["object"].(map[string]any); ok {
					custID, _ := obj["customer"].(string)
					subID, _ := obj["id"].(string)
					status, _ := obj["status"].(string)
					_, _ = pool.Exec(r.Context(),
						`UPDATE subscriptions SET stripe_subscription_id=$2, status=$3 WHERE stripe_customer_id=$1`,
						custID, subID, status)
				}
			}
		case "customer.subscription.deleted":
			if data, ok := event["data"].(map[string]any); ok {
				if obj, ok := data["object"].(map[string]any); ok {
					subID, _ := obj["id"].(string)
					_, _ = pool.Exec(r.Context(),
						`UPDATE subscriptions SET status='canceled' WHERE stripe_subscription_id=$1`, subID)
				}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"received": "true"})
	}
}
