package middleware

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RequireTenant(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := UserIDFrom(r.Context())
			tenantID := TenantIDFrom(r.Context())

			if userID == "" || tenantID == "" {
				http.Error(w, "tenant context required", http.StatusForbidden)
				return
			}

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

func WithTenantTx(ctx context.Context, pool *pgxpool.Pool, fn func(pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", TenantIDFrom(ctx)); err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
