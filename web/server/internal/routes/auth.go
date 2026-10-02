package routes

import (
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// MountAuth registers POST /auth/sync behind token verification + login
// rate limit. VerifyOnly (not the full profile middleware): the profile
// may not exist yet, and creating it is this endpoint's job.
// The handle_new_user trigger usually creates the profile; sync is the fallback.
func MountAuth(r chi.Router, h *handler.Auth, auth *middleware.Auth) {
	r.With(auth.VerifyOnly, middleware.LoginLimiter).Post("/auth/sync", h.Sync)
}
