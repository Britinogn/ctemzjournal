package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/model"
	"github.com/britinogn/ctemzjournal/pkg/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ctxKey struct{}
type userIDKey struct{}

// Auth verifies the Supabase JWT via JWKS, loads profiles.role from the
// database (never token metadata), and puts *model.User in the context.
// Suspended users get 403 on every authenticated route.
//
// POST /auth/sync is the one exception: it uses VerifyOnly (valid token,
// no profile row required) so first logins can bootstrap the profile
// instead of deadlocking on "profile not found".
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

// UserIDFromContext returns the token subject saved by VerifyOnly, or false.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	uid, ok := ctx.Value(userIDKey{}).(uuid.UUID)
	return uid, ok
}

// verifyToken checks the Bearer token and returns the subject.
// The message is human-readable: it surfaces in API responses.
func (a *Auth) verifyToken(r *http.Request) (uid uuid.UUID, msg string, ok bool) {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return uuid.Nil, "No session found. Please log in again.", false
	}
	raw := strings.TrimPrefix(header, "Bearer ")

	token, err := jwt.Parse(raw, a.kf.Keyfunc)
	if err != nil || !token.Valid {
		return uuid.Nil, "Your session has expired or is invalid. Please log in again.", false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, "Your session could not be read. Please log in again.", false
	}
	sub, _ := claims["sub"].(string)
	uid, err = uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, "Your session has no valid user. Please log in again.", false
	}
	return uid, "", true
}

// VerifyOnly checks the token and stores its subject, without requiring a
// profile row. Only POST /auth/sync uses this (first-login bootstrap).
func (a *Auth) VerifyOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, msg, ok := a.verifyToken(r)
		if !ok {
			response.Error(w, http.StatusUnauthorized, msg)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey{}, uid)))
	})
}

func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, msg, ok := a.verifyToken(r)
		if !ok {
			response.Error(w, http.StatusUnauthorized, msg)
			return
		}

		profile, err := a.queries.GetProfileByID(r.Context(), uid)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Account record not found. Please sync your account and try again.")
			return
		}
		user := &model.User{
			ID:     profile.ID.String(),
			Role:   profile.Role,
			Status: profile.Status,
		}
		if user.IsSuspended() {
			response.Error(w, http.StatusForbidden, "This account has been suspended. Contact support if this is a mistake.")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, user)))
	})
}
