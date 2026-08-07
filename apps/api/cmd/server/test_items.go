package main

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	chimiddleware "github.com/Suggus1899/chiguire/api/internal/middleware"
)

func handleListTestItems(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := chimiddleware.TenantIDFrom(r.Context())
		rows, err := pool.Query(r.Context(),
			`SELECT id, label, created_at FROM test_items WHERE tenant_id = $1 ORDER BY created_at DESC`,
			tenantID,
		)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type item struct {
			ID        string `json:"id"`
			Label     string `json:"label"`
			CreatedAt string `json:"created_at"`
		}
		var items []item
		for rows.Next() {
			var it item
			if err := rows.Scan(&it.ID, &it.Label, &it.CreatedAt); err != nil {
				http.Error(w, "scan error", http.StatusInternalServerError)
				return
			}
			items = append(items, it)
		}
		if items == nil {
			items = []item{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func handleCreateTestItem(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := chimiddleware.TenantIDFrom(r.Context())
		var body struct {
			Label string `json:"label"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Label == "" {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		var id, createdAt string
		err := pool.QueryRow(r.Context(),
			`INSERT INTO test_items (tenant_id, label) VALUES ($1, $2) RETURNING id, created_at`,
			tenantID, body.Label,
		).Scan(&id, &createdAt)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id, "label": body.Label, "created_at": createdAt})
	}
}
