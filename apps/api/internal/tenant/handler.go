package tenant

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

type Tenant struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Plan        string     `json:"plan"`
	TrialEndsAt *time.Time `json:"trial_ends_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func HandleCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mw.UserIDFrom(r.Context())

		var body struct {
			Name string `json:"name"`
			Slug string `json:"slug"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.Slug == "" {
			http.Error(w, "name and slug are required", http.StatusBadRequest)
			return
		}
		body.Slug = strings.ToLower(strings.ReplaceAll(body.Slug, " ", "-"))

		tx, err := pool.Begin(r.Context())
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(r.Context())

		trialEnd := time.Now().Add(7 * 24 * time.Hour)
		var t Tenant
		err = tx.QueryRow(r.Context(),
			`INSERT INTO tenants (name, slug, plan, trial_ends_at)
			 VALUES ($1,$2,'trial',$3)
			 RETURNING id, name, slug, plan, trial_ends_at, created_at`,
			body.Name, body.Slug, trialEnd,
		).Scan(&t.ID, &t.Name, &t.Slug, &t.Plan, &t.TrialEndsAt, &t.CreatedAt)
		if err != nil {
			http.Error(w, "slug already taken", http.StatusConflict)
			return
		}

		_, err = tx.Exec(r.Context(),
			`INSERT INTO user_tenants (user_id, tenant_id, role) VALUES ($1,$2,'owner')`,
			userID, t.ID,
		)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(r.Context()); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	}
}

func HandleList(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mw.UserIDFrom(r.Context())

		rows, err := pool.Query(r.Context(),
			`SELECT t.id, t.name, t.slug, t.plan, t.trial_ends_at, t.created_at
			 FROM tenants t
			 JOIN user_tenants ut ON ut.tenant_id = t.id
			 WHERE ut.user_id = $1
			 ORDER BY t.created_at DESC`,
			userID,
		)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tenants []Tenant
		for rows.Next() {
			var t Tenant
			if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.Plan, &t.TrialEndsAt, &t.CreatedAt); err != nil {
				http.Error(w, "scan error", http.StatusInternalServerError)
				return
			}
			tenants = append(tenants, t)
		}
		if tenants == nil {
			tenants = []Tenant{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tenants)
	}
}
