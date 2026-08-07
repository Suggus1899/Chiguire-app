package middleware

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RequireTenant ensures the JWT has a tenant_id AND that the user belongs to it.
func RequireTenant(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := UserIDFrom(r.Context())
			tenantID := TenantIDFrom(r.Context())

			if userID == "" || tenantID == "" {
				http.Error(w, "tenant context required", http.StatusForbidden)
				return
			}

			// Verify membership
			var exists bool
			err := pool.QueryRow(r.Context(),
				`SELECT EXISTS(SELECT 1 FROM user_tenants WHERE user_id=$1 AND tenant_id=$2)`,
				userID, tenantID,
			).Scan(&exists)
			if err != nil || !exists {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
