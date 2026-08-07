package vendor

import (
	"log"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

type Vendor struct {
	ID                string    `json:"id"`
	Rif               string    `json:"rif"`
	Name              string    `json:"name"`
	Address           string    `json:"address"`
	Phone             string    `json:"phone"`
	Email             string    `json:"email"`
	IsSpecialTaxpayer bool      `json:"is_special_taxpayer"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

var errNotFound = errors.New("not found")

const vendorCols = `id, rif, name, COALESCE(address,''), COALESCE(phone,''), COALESCE(email,''), is_special_taxpayer, created_at, updated_at`

func HandleCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Rif               string `json:"rif"`
			Name              string `json:"name"`
			Address           string `json:"address"`
			Phone             string `json:"phone"`
			Email             string `json:"email"`
			IsSpecialTaxpayer bool   `json:"is_special_taxpayer"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Rif == "" || body.Name == "" {
			http.Error(w, "rif and name are required", http.StatusBadRequest)
			return
		}

		var v Vendor
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO vendors (tenant_id, rif, name, address, phone, email, is_special_taxpayer)
				 VALUES ($1,$2,$3,$4,$5,$6,$7)
				 RETURNING `+vendorCols,
				mw.TenantIDFrom(r.Context()), body.Rif, body.Name, body.Address, body.Phone, body.Email, body.IsSpecialTaxpayer,
			).Scan(&v.ID, &v.Rif, &v.Name, &v.Address, &v.Phone, &v.Email, &v.IsSpecialTaxpayer, &v.CreatedAt, &v.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "rif already exists", http.StatusConflict)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(v)
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
		var vendors []Vendor
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+vendorCols+` FROM vendors ORDER BY created_at DESC LIMIT $1 OFFSET $2`, perPage, offset)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var v Vendor
				if err := rows.Scan(&v.ID, &v.Rif, &v.Name, &v.Address, &v.Phone, &v.Email, &v.IsSpecialTaxpayer, &v.CreatedAt, &v.UpdatedAt); err != nil {
					return err
				}
				vendors = append(vendors, v)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if vendors == nil {
			vendors = []Vendor{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vendors)
	}
}

func HandleGet(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var v Vendor
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`SELECT `+vendorCols+` FROM vendors WHERE id=$1`, id,
			).Scan(&v.ID, &v.Rif, &v.Name, &v.Address, &v.Phone, &v.Email, &v.IsSpecialTaxpayer, &v.CreatedAt, &v.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(v)
	}
}

func HandleUpdate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var body struct {
			Rif               string `json:"rif"`
			Name              string `json:"name"`
			Address           string `json:"address"`
			Phone             string `json:"phone"`
			Email             string `json:"email"`
			IsSpecialTaxpayer bool   `json:"is_special_taxpayer"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		var v Vendor
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`UPDATE vendors SET rif=$1, name=$2, address=$3, phone=$4, email=$5, is_special_taxpayer=$6, updated_at=NOW()
				 WHERE id=$7
				 RETURNING `+vendorCols,
				body.Rif, body.Name, body.Address, body.Phone, body.Email, body.IsSpecialTaxpayer, id,
			).Scan(&v.ID, &v.Rif, &v.Name, &v.Address, &v.Phone, &v.Email, &v.IsSpecialTaxpayer, &v.CreatedAt, &v.UpdatedAt)
		})
		if err != nil {
			http.Error(w, "not found or conflict", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(v)
	}
}

func HandleDelete(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(), `DELETE FROM vendors WHERE id=$1`, id)
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
