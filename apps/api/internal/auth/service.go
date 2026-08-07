package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

type Service struct {
	pool      *pgxpool.Pool
	jwtSecret string
}

func NewService(pool *pgxpool.Pool, secret string) *Service {
	return &Service{pool: pool, jwtSecret: secret}
}

func (s *Service) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Email == "" || body.Password == "" || body.FullName == "" {
		http.Error(w, "email, password and full_name are required", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	var userID string
	err = s.pool.QueryRow(r.Context(),
		`INSERT INTO users (email, password_hash, full_name) VALUES ($1,$2,$3) RETURNING id`,
		body.Email, string(hash), body.FullName,
	).Scan(&userID)
	if err != nil {
		http.Error(w, "email already registered", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
}

func (s *Service) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		TenantID string `json:"tenant_id"` // optional; scopes the token
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	var userID, hash string
	err := s.pool.QueryRow(r.Context(),
		`SELECT id, password_hash FROM users WHERE email=$1`,
		body.Email,
	).Scan(&userID, &hash)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(body.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	access, refresh, err := s.issueTokenPair(r.Context(), userID, body.TenantID, "member")
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (s *Service) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.RefreshToken), 4) // low cost for lookup
	_ = hash
	// ponytail: using raw token as lookup key (hashed at insert); full rotation in phase 1
	var userID, tenantID, role string
	err = s.pool.QueryRow(r.Context(),
		`SELECT user_id, tenant_id, 'member' FROM refresh_tokens
		 WHERE token_hash = crypt($1, token_hash)
		 AND expires_at > NOW() AND revoked_at IS NULL`,
		body.RefreshToken,
	).Scan(&userID, &tenantID, &role)
	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	access, _, err := s.issueTokenPair(r.Context(), userID, tenantID, role)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": access})
}

func (s *Service) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// For now: client just discards the token. Full revocation (allowlist) when needed.
	// ponytail: stateless logout; add token revocation if session invalidation is required
	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) issueTokenPair(ctx context.Context, userID, tenantID, role string) (access, refresh string, err error) {
	now := time.Now()

	claims := mw.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
	}
	access, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		return
	}

	// Refresh token: opaque random string stored hashed
	rawRefresh := generateToken()
	_, err = s.pool.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, tenant_id, token_hash, expires_at)
		 VALUES ($1, $2, crypt($3, gen_salt('bf', 4)), NOW() + INTERVAL '7 days')`,
		userID, tenantID, rawRefresh,
	)
	if err != nil {
		return
	}
	refresh = rawRefresh
	return
}
