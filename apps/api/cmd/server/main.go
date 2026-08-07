package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/Suggus1899/chiguire/api/internal/auth"
	chimiddleware "github.com/Suggus1899/chiguire/api/internal/middleware"
	"github.com/Suggus1899/chiguire/api/internal/powersync"
	"github.com/Suggus1899/chiguire/api/internal/tenant"
)

func main() {
	_ = godotenv.Load()

	dbURL := mustEnv("DATABASE_URL")
	jwtSecret := mustEnv("JWT_SECRET")
	psJWTSecret := mustEnv("POWERSYNC_JWT_SECRET")
	port := getEnv("PORT", "3001")

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	authSvc := auth.NewService(pool, jwtSecret)
	psSvc := powersync.NewService(psJWTSecret)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Public
	r.Post("/auth/register", authSvc.HandleRegister)
	r.Post("/auth/login", authSvc.HandleLogin)
	r.Post("/auth/refresh", authSvc.HandleRefresh)

	// Authenticated
	r.Group(func(r chi.Router) {
		r.Use(chimiddleware.RequireAuth(jwtSecret))
		r.Post("/auth/logout", authSvc.HandleLogout)

		// PowerSync: return a short-lived token scoped to the user's current tenant
		r.Get("/powersync/token", psSvc.HandleToken)

		// Tenant management
		r.Route("/tenants", func(r chi.Router) {
			r.Post("/", tenant.HandleCreate(pool))
			r.Get("/", tenant.HandleList(pool))
		})

		// Phase 0 test endpoint
		r.Route("/test-items", func(r chi.Router) {
			r.Use(chimiddleware.RequireTenant(pool))
			r.Get("/", handleListTestItems(pool))
			r.Post("/", handleCreateTestItem(pool))
		})
	})

	log.Printf("API listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
