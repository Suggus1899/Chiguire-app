package seller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

var errNotFound = errors.New("not found")

type Seller struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	CommissionPct string    `json:"commission_pct"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

const sellerCols = `id, COALESCE(user_id::text,''), name, COALESCE(email,''), COALESCE(phone,''), commission_pct::text, is_active, created_at, updated_at`

func HandleCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			UserID        string `json:"user_id"`
			Name          string `json:"name"`
			Email         string `json:"email"`
			Phone         string `json:"phone"`
			CommissionPct string `json:"commission_pct"`
			IsActive      *bool  `json:"is_active"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		if body.CommissionPct == "" {
			body.CommissionPct = "0"
		}
		isActive := true
		if body.IsActive != nil {
			isActive = *body.IsActive
		}

		var s Seller
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO sellers (tenant_id, user_id, name, email, phone, commission_pct, is_active)
				 VALUES ($1,NULLIF($2,'')::uuid,$3,$4,$5,$6,$7)
				 RETURNING `+sellerCols,
				mw.TenantIDFrom(r.Context()), body.UserID, body.Name, body.Email, body.Phone, body.CommissionPct, isActive,
			).Scan(&s.ID, &s.UserID, &s.Name, &s.Email, &s.Phone, &s.CommissionPct, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(s)
	}
}

func HandleList(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		var sellers []Seller
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+sellerCols+` FROM sellers ORDER BY created_at DESC LIMIT $1 OFFSET $2`, perPage, offset)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var s Seller
				if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.Email, &s.Phone, &s.CommissionPct, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
					return err
				}
				sellers = append(sellers, s)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if sellers == nil {
			sellers = []Seller{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sellers)
	}
}

func HandleGet(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var s Seller
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`SELECT `+sellerCols+` FROM sellers WHERE id=$1`, id,
			).Scan(&s.ID, &s.UserID, &s.Name, &s.Email, &s.Phone, &s.CommissionPct, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)
	}
}

func HandleUpdate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			UserID        string `json:"user_id"`
			Name          string `json:"name"`
			Email         string `json:"email"`
			Phone         string `json:"phone"`
			CommissionPct string `json:"commission_pct"`
			IsActive      *bool  `json:"is_active"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.CommissionPct == "" {
			body.CommissionPct = "0"
		}

		var s Seller
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			if body.IsActive != nil {
				return tx.QueryRow(r.Context(),
					`UPDATE sellers SET user_id=NULLIF($1,'')::uuid, name=$2, email=$3, phone=$4, commission_pct=$5, is_active=$6, updated_at=NOW()
					 WHERE id=$7
					 RETURNING `+sellerCols,
					body.UserID, body.Name, body.Email, body.Phone, body.CommissionPct, *body.IsActive, id,
				).Scan(&s.ID, &s.UserID, &s.Name, &s.Email, &s.Phone, &s.CommissionPct, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
			}
			return tx.QueryRow(r.Context(),
				`UPDATE sellers SET user_id=NULLIF($1,'')::uuid, name=$2, email=$3, phone=$4, commission_pct=$5, updated_at=NOW()
				 WHERE id=$6
				 RETURNING `+sellerCols,
				body.UserID, body.Name, body.Email, body.Phone, body.CommissionPct, id,
			).Scan(&s.ID, &s.UserID, &s.Name, &s.Email, &s.Phone, &s.CommissionPct, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)
	}
}

func HandleDelete(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(), `DELETE FROM sellers WHERE id=$1`, id)
			if err != nil {
				return err
			}
			if ct.RowsAffected() == 0 {
				return errNotFound
			}
			return nil
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleListCommissions(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sellerID := chi.URLParam(r, "id")
		period := r.URL.Query().Get("period")
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
		type commissionRow struct {
			ID               string `json:"id"`
			VendorUserID     string `json:"vendor_user_id"`
			InvoiceID        string `json:"invoice_id"`
			RatePct          string `json:"rate_pct"`
			BaseAmount       string `json:"base_amount"`
			CommissionAmount string `json:"commission_amount"`
			Period           string `json:"period"`
			Status           string `json:"status"`
		}
		var rows []commissionRow
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			q, err := tx.Query(r.Context(),
				`SELECT sc.id, sc.vendor_user_id::text, COALESCE(sc.invoice_id::text,''), sc.rate_pct::text,
				        sc.base_amount::text, sc.commission_amount::text, sc.period, COALESCE(sc.status,'pending')
				 FROM sales_commissions sc
				 JOIN sellers s ON s.user_id = sc.vendor_user_id
				 WHERE s.id=$1 AND ($2='' OR sc.period=$2)
				 ORDER BY sc.period DESC
				 LIMIT $3 OFFSET $4`, sellerID, period, perPage, offset)
			if err != nil {
				return err
			}
			defer q.Close()
			for q.Next() {
				var c commissionRow
				if err := q.Scan(&c.ID, &c.VendorUserID, &c.InvoiceID, &c.RatePct, &c.BaseAmount, &c.CommissionAmount, &c.Period, &c.Status); err != nil {
					return err
				}
				rows = append(rows, c)
			}
			return q.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if rows == nil {
			rows = []commissionRow{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}

func HandleMarkCommissionsPaid(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			CommissionIDs []string `json:"commission_ids"`
			Reference     string   `json:"reference"`
			Notes         string   `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if len(body.CommissionIDs) == 0 {
			http.Error(w, "commission_ids required", http.StatusBadRequest)
			return
		}

		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			// Fetch all seller IDs for the given commission IDs in a single query.
			rows, err := tx.Query(r.Context(),
				`SELECT s.id FROM sales_commissions sc
				 JOIN sellers s ON s.user_id = sc.vendor_user_id
				 WHERE sc.id = ANY($1)`, body.CommissionIDs)
			if err != nil {
				return err
			}
			defer rows.Close()
			sellerSet := make(map[string]struct{})
			for rows.Next() {
				var sid string
				if err := rows.Scan(&sid); err != nil {
					return err
				}
				sellerSet[sid] = struct{}{}
			}
			if err := rows.Err(); err != nil {
				return err
			}
			if len(sellerSet) == 0 {
				return errors.New("commission not found")
			}
			if len(sellerSet) > 1 {
				return errors.New("all commissions must belong to the same seller")
			}
			// Single UPDATE for all matching commission IDs.
			_, err = tx.Exec(r.Context(),
				`UPDATE sales_commissions SET status='paid' WHERE id = ANY($1)`, body.CommissionIDs)
			return err
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"status": "paid", "count": len(body.CommissionIDs)})
	}
}

func HandleListCommissionLimits(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sellerID := chi.URLParam(r, "id")
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
		type limitRow struct {
			ProductID   string `json:"product_id"`
			ProductName string `json:"product_name"`
			RatePct     string `json:"rate_pct"`
		}
		var rows []limitRow
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			q, err := tx.Query(r.Context(),
				`SELECT p.id::text, p.name, COALESCE(s.commission_pct,0)::text
				 FROM products p
				 CROSS JOIN sellers s
				 WHERE s.id=$1
				 ORDER BY p.name
				 LIMIT $2 OFFSET $3`, sellerID, perPage, offset)
			if err != nil {
				return err
			}
			defer q.Close()
			for q.Next() {
				var l limitRow
				if err := q.Scan(&l.ProductID, &l.ProductName, &l.RatePct); err != nil {
					return err
				}
				rows = append(rows, l)
			}
			return q.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if rows == nil {
			rows = []limitRow{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}
