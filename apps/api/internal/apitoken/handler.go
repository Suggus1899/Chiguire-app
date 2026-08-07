package apitoken

import (
	"log"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

var errNotFound = errors.New("not found")

type APIToken struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	IsActive   bool       `json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

const tokenCols = `id, name, is_active, last_used_at, expires_at, created_at`

func generateToken() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	token := "chig_" + hex.EncodeToString(b)
	h := sha256.Sum256([]byte(token))
	return token, hex.EncodeToString(h[:]), nil
}

func HandleCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name      string `json:"name"`
			ExpiresAt string `json:"expires_at"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		plaintext, hash, err := generateToken()
		if err != nil {
			http.Error(w, "token generation failed", http.StatusInternalServerError)
			return
		}

		var t APIToken
		err = mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			return tx.QueryRow(r.Context(),
				`INSERT INTO api_tokens (tenant_id, user_id, name, token_hash, expires_at)
				 VALUES ($1,$2,$3,$4,NULLIF($5,'')::timestamptz)
				 RETURNING `+tokenCols,
				mw.TenantIDFrom(r.Context()), mw.UserIDFrom(r.Context()), body.Name, hash, body.ExpiresAt,
			).Scan(&t.ID, &t.Name, &t.IsActive, &t.LastUsedAt, &t.ExpiresAt, &t.CreatedAt)
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"token":      plaintext,
			"id":         t.ID,
			"name":       t.Name,
			"is_active":  t.IsActive,
			"expires_at": t.ExpiresAt,
			"created_at": t.CreatedAt,
		})
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
		var tokens []APIToken
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			rows, err := tx.Query(r.Context(),
				`SELECT `+tokenCols+` FROM api_tokens ORDER BY created_at DESC LIMIT $1 OFFSET $2`, perPage, offset)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var t APIToken
				if err := rows.Scan(&t.ID, &t.Name, &t.IsActive, &t.LastUsedAt, &t.ExpiresAt, &t.CreatedAt); err != nil {
					return err
				}
				tokens = append(tokens, t)
			}
			return rows.Err()
		})
		if err != nil {
			log.Printf("db error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if tokens == nil {
			tokens = []APIToken{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokens)
	}
}

func HandleRevoke(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := mw.WithTenantTx(r.Context(), pool, func(tx pgx.Tx) error {
			ct, err := tx.Exec(r.Context(),
				`UPDATE api_tokens SET is_active=false WHERE id=$1`, id)
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
