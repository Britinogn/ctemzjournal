package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/model"
	"github.com/britinogn/ctemzjournal/pkg/response"
)

type ctxKey struct{}

// Auth verifies the Supabase JWT via JWKS, loads profiles.role from the
// database (never token metadata), and puts *model.User in the context.
// Suspended users get 403 on every authenticated route.
type Auth struct {
	queries *sqlc.Queries
	kf      keyfunc.Keyfunc
}

// NewAuth builds the middleware. jwksURL is SUPABASE_JWKS_URL.
func NewAuth(ctx context.Context, q *sqlc.Queries, jwksURL string) (*Auth, error) {
	kf, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, err
	}
	return &Auth{queries: q, kf: kf}, nil
}

// FromContext returns the authenticated user, or nil.
func FromContext(ctx context.Context) *model.User {
	u, _ := ctx.Value(ctxKey{}).(*model.User)
	return u
}

func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Error(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		raw := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(raw, a.kf.Keyfunc)
		if err != nil || !token.Valid {
			response.Error(w, http.StatusUnauthorized, "invalid token")
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(w, http.StatusUnauthorized, "invalid claims")
			return
		}
		sub, _ := claims["sub"].(string)
		uid, err := uuid.Parse(sub)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "invalid subject")
			return
		}

		profile, err := a.queries.GetProfileByID(r.Context(), uid)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "profile not found")
			return
		}
		user := &model.User{
			ID:     profile.ID.String(),
			Role:   profile.Role,
			Status: profile.Status,
		}
		if user.IsSuspended() {
			response.Error(w, http.StatusForbidden, "account suspended")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, user)))
	})
}
