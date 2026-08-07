package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/Suggus1899/chiguire/api/internal/auth"
	"github.com/Suggus1899/chiguire/api/internal/commission"
	"github.com/Suggus1899/chiguire/api/internal/customer"
	"github.com/Suggus1899/chiguire/api/internal/delivery"
	"github.com/Suggus1899/chiguire/api/internal/fiscal"
	"github.com/Suggus1899/chiguire/api/internal/importexport"
	chimiddleware "github.com/Suggus1899/chiguire/api/internal/middleware"
	"github.com/Suggus1899/chiguire/api/internal/inventory"
	"github.com/Suggus1899/chiguire/api/internal/invoice"
	"github.com/Suggus1899/chiguire/api/internal/payment"
	"github.com/Suggus1899/chiguire/api/internal/powersync"
	"github.com/Suggus1899/chiguire/api/internal/product"
	"github.com/Suggus1899/chiguire/api/internal/purchase"
	"github.com/Suggus1899/chiguire/api/internal/quotation"
	"github.com/Suggus1899/chiguire/api/internal/saas"
	"github.com/Suggus1899/chiguire/api/internal/tenant"
	"github.com/Suggus1899/chiguire/api/internal/vendor"
	"github.com/Suggus1899/chiguire/api/internal/webhook"
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

	// Background workers
	go webhook.StartWebhookWorker(pool)
	go fiscal.StartBCVCron(context.Background(), pool, 6*time.Hour)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Public
	r.Post("/auth/register", authSvc.HandleRegister)
	r.Post("/auth/login", authSvc.HandleLogin)
	r.Post("/auth/refresh", authSvc.HandleRefresh)

	// Payment webhooks (public — signed by provider)
	r.Post("/webhooks/payments/{provider}", payment.HandleWebhook(pool))
	r.Post("/webhooks/stripe", saas.HandleStripeWebhook(pool))

	// Authenticated
	r.Group(func(r chi.Router) {
		r.Use(chimiddleware.RequireAuth(jwtSecret))
		r.Post("/auth/logout", authSvc.HandleLogout)
		r.Get("/powersync/token", psSvc.HandleToken)

		// Tenant management
		r.Route("/tenants", func(r chi.Router) {
			r.Post("/", tenant.HandleCreate(pool))
			r.Get("/", tenant.HandleList(pool))
		})

		// Tenant-scoped routes (require active tenant in JWT)
		r.Group(func(r chi.Router) {
			r.Use(chimiddleware.RequireTenant(pool))

			// Phase 0 test endpoint
			r.Route("/test-items", func(r chi.Router) {
				r.Get("/", handleListTestItems(pool))
				r.Post("/", handleCreateTestItem(pool))
			})

			// Phase 1: Customers
			r.Route("/customers", func(r chi.Router) {
				r.Post("/", customer.HandleCreate(pool))
				r.Get("/", customer.HandleList(pool))
				r.Get("/{id}", customer.HandleGet(pool))
				r.Put("/{id}", customer.HandleUpdate(pool))
				r.Delete("/{id}", customer.HandleDelete(pool))
			})

			// Phase 1: Vendors
			r.Route("/vendors", func(r chi.Router) {
				r.Post("/", vendor.HandleCreate(pool))
				r.Get("/", vendor.HandleList(pool))
				r.Get("/{id}", vendor.HandleGet(pool))
				r.Put("/{id}", vendor.HandleUpdate(pool))
				r.Delete("/{id}", vendor.HandleDelete(pool))
			})

			// Phase 1: Products, categories, units, prices
			r.Route("/products", func(r chi.Router) {
				r.Post("/", product.HandleCreateProduct(pool))
				r.Get("/", product.HandleListProducts(pool))
				r.Get("/{id}", product.HandleGetProduct(pool))
				r.Put("/{id}", product.HandleUpdateProduct(pool))
				r.Delete("/{id}", product.HandleDeleteProduct(pool))
				r.Get("/{id}/prices", product.HandleListPrices(pool))
				r.Post("/{id}/prices", product.HandleCreatePrice(pool))
			})
			r.Route("/categories", func(r chi.Router) {
				r.Post("/", product.HandleCreateCategory(pool))
				r.Get("/", product.HandleListCategories(pool))
				r.Delete("/{id}", product.HandleDeleteCategory(pool))
			})
			r.Route("/units", func(r chi.Router) {
				r.Post("/", product.HandleCreateUnit(pool))
				r.Get("/", product.HandleListUnits(pool))
				r.Delete("/{id}", product.HandleDeleteUnit(pool))
			})

			// Phase 1: Inventory
			r.Route("/branches", func(r chi.Router) {
				r.Post("/", inventory.HandleCreateBranch(pool))
				r.Get("/", inventory.HandleListBranches(pool))
				r.Delete("/{id}", inventory.HandleDeleteBranch(pool))
			})
			r.Route("/warehouses", func(r chi.Router) {
				r.Post("/", inventory.HandleCreateWarehouse(pool))
				r.Get("/", inventory.HandleListWarehouses(pool))
				r.Delete("/{id}", inventory.HandleDeleteWarehouse(pool))
			})
			r.Route("/stock", func(r chi.Router) {
				r.Post("/movements", inventory.HandleCreateMovement(pool))
				r.Get("/movements", inventory.HandleListMovements(pool))
				r.Get("/", inventory.HandleGetStock(pool))
			})

			// Phase 1: Invoices
			r.Route("/invoices", func(r chi.Router) {
				r.Post("/", invoice.HandleCreate(pool))
				r.Get("/", invoice.HandleList(pool))
				r.Get("/{id}", invoice.HandleGet(pool))
				r.Put("/{id}", invoice.HandleUpdate(pool))
				r.Delete("/{id}", invoice.HandleVoid(pool))
				r.Post("/{id}/items", invoice.HandleAddItem(pool))
				r.Delete("/{id}/items/{itemId}", invoice.HandleRemoveItem(pool))
				r.Post("/{id}/emit", invoice.HandleEmit(pool))
				r.Post("/{id}/payments", invoice.HandleAddPayment(pool))
			})

			// Phase 2: Fiscal
			r.Route("/fiscal", func(r chi.Router) {
				r.Get("/tax-categories", fiscal.HandleListTaxCategories(pool))
				r.Post("/tax-categories", fiscal.HandleCreateTaxCategory(pool))
				r.Post("/withholdings", fiscal.HandleCreateWithholding(pool))
				r.Get("/withholdings", fiscal.HandleListWithholdings(pool))
				r.Get("/periods", fiscal.HandleListFiscalPeriods(pool))
				r.Post("/periods/{period}/close", fiscal.HandleClosePeriod(pool))
				r.Post("/books/generate", fiscal.HandleGenerateFiscalBook(pool))
				r.Get("/books/{id}", fiscal.HandleGetFiscalBook(pool))
				r.Get("/exchange-rate", fiscal.HandleGetExchangeRate(pool))
				r.Post("/exchange-rate", fiscal.HandleSetExchangeRate(pool))
				r.Get("/exchange-rates", fiscal.HandleListExchangeRates(pool))
			})

			// Phase 4: Payments
			r.Route("/payments", func(r chi.Router) {
				r.Post("/links", payment.HandleCreatePaymentLink(pool))
				r.Get("/links", payment.HandleListPaymentLinks(pool))
				r.Get("/links/{code}", payment.HandleGetPaymentLink(pool))
			})

			// Phase 5: Purchases
			r.Route("/purchases", func(r chi.Router) {
				r.Post("/", purchase.HandleCreatePO(pool))
				r.Get("/", purchase.HandleListPOs(pool))
				r.Get("/{id}", purchase.HandleGetPO(pool))
				r.Put("/{id}", purchase.HandleUpdatePO(pool))
				r.Post("/{id}/approve", purchase.HandleApprovePO(pool))
				r.Post("/{id}/items", purchase.HandleAddPOItem(pool))
				r.Post("/{id}/items/{itemId}/receive", purchase.HandleReceivePOItem(pool))
			})

			// Phase 5: Quotations
			r.Route("/quotations", func(r chi.Router) {
				r.Post("/", quotation.HandleCreate(pool))
				r.Get("/", quotation.HandleList(pool))
				r.Get("/{id}", quotation.HandleGet(pool))
				r.Put("/{id}", quotation.HandleUpdate(pool))
				r.Post("/{id}/items", quotation.HandleAddItem(pool))
				r.Delete("/{id}/items/{itemId}", quotation.HandleRemoveItem(pool))
				r.Post("/{id}/convert", quotation.HandleConvertToInvoice(pool))
				r.Post("/{id}/send", quotation.HandleSendQuote(pool))
			})

			// Phase 5: Commissions
			r.Route("/commissions", func(r chi.Router) {
				r.Get("/", commission.HandleListCommissions(pool))
				r.Post("/calculate", commission.HandleCalculateCommissions(pool))
			})

			// Phase 5: Delivery routes
			r.Route("/delivery", func(r chi.Router) {
				r.Post("/routes", delivery.HandleCreateRoute(pool))
				r.Get("/routes", delivery.HandleListRoutes(pool))
				r.Get("/routes/{id}", delivery.HandleGetRoute(pool))
				r.Post("/routes/{id}/stops", delivery.HandleAddStop(pool))
				r.Put("/stops/{stopId}", delivery.HandleUpdateStopStatus(pool))
				r.Post("/stops/{stopId}/signature", delivery.HandleRecordSignature(pool))
			})

			// Phase 6: SaaS
			r.Route("/saas", func(r chi.Router) {
				r.Get("/subscription", saas.HandleGetSubscription(pool))
				r.Put("/subscription", saas.HandleUpdatePlan(pool))
				r.Get("/plan-limits", saas.HandleGetPlanLimits(pool))
				r.Post("/checkout", saas.HandleCreateStripeCheckoutSession(pool))
			})

			// Phase 6: Webhooks
			r.Route("/webhooks", func(r chi.Router) {
				r.Post("/", webhook.HandleCreateWebhook(pool))
				r.Get("/", webhook.HandleListWebhooks(pool))
				r.Delete("/{id}", webhook.HandleDeleteWebhook(pool))
				r.Post("/{id}/test", webhook.HandleTestWebhook(pool))
				r.Get("/deliveries", webhook.HandleListDeliveries(pool))
			})

			// Phase 6: Import/Export
			r.Route("/import", func(r chi.Router) {
				r.Post("/customers", importexport.HandleImportCustomers(pool))
				r.Post("/products", importexport.HandleImportProducts(pool))
			})
			r.Route("/export", func(r chi.Router) {
				r.Get("/invoices", importexport.HandleExportInvoices(pool))
				r.Get("/customers", importexport.HandleExportCustomers(pool))
			})
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
