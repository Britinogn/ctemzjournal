package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims are the app's own JWT claims (password reset / app sessions).
// Supabase Auth tokens are verified separately via JWKS in middleware/auth.go.
type Claims struct {
	UserID string `json:"uid"`
	Kind   string `json:"kind"`
	jwt.RegisteredClaims
}

// Sign mints an HS256 app token (kind e.g. "reset").
func Sign(secret, userID, kind string, ttl time.Duration, now time.Time) (string, error) {
	claims := Claims{
		UserID: userID,
		Kind:   kind,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// Verify checks an app token and returns its claims.
func Verify(secret, token string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
