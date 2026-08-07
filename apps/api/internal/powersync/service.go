package powersync

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	mw "github.com/Suggus1899/chiguire/api/internal/middleware"
)

// Service issues short-lived PowerSync tokens scoped to user + tenant.
// PowerSync validates these against the secret configured in config.yaml.
type Service struct {
	secret string
}

func NewService(secret string) *Service {
	return &Service{secret: secret}
}

type psClaims struct {
	jwt.RegisteredClaims
	// PowerSync reads these fields from the token
	Sub      string            `json:"sub"`      // user_id
	Parameters map[string]string `json:"parameters"` // available in sync rules as token_parameters.*
}

func (s *Service) HandleToken(w http.ResponseWriter, r *http.Request) {
	userID := mw.UserIDFrom(r.Context())
	tenantID := mw.TenantIDFrom(r.Context())

	if tenantID == "" {
		http.Error(w, "tenant_id required in JWT to get PowerSync token", http.StatusBadRequest)
		return
	}

	now := time.Now()
	claims := psClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute)),
		},
		Sub: userID,
		Parameters: map[string]string{
			"tenant_id": tenantID,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.secret))
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      token,
		"expires_at": claims.ExpiresAt.Time.Unix(),
	})
}
