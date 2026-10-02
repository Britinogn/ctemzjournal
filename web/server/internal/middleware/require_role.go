package middleware

import (
	"net/http"

	"github.com/britinogn/ctemzjournal/pkg/response"
)

// RequireRole allows only the given role (reads the DB-loaded role from context).
// Admin routes use RequireRole("admin"); user routes use RequireRole("user").
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := FromContext(r.Context())
			if u == nil {
				response.Error(w, http.StatusUnauthorized, "unauthenticated")
				return
			}
			if u.Role != role {
				response.Error(w, http.StatusForbidden, "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAuth allows any authenticated, non-suspended user.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if FromContext(r.Context()) == nil {
			response.Error(w, http.StatusUnauthorized, "unauthenticated")
			return
		}
		next.ServeHTTP(w, r)
	})
}
