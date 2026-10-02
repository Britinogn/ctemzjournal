package routes

import (
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// MountAuth registers POST /auth/sync behind auth + login rate limit.
// The handle_new_user trigger usually creates the profile; sync is the fallback.
func MountAuth(r chi.Router, h *handler.Auth, auth *middleware.Auth) {
	r.With(auth.Middleware, middleware.LoginLimiter).Post("/auth/sync", h.Sync)
}
