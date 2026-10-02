package routes

import (
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// MountMe registers GET /me and PATCH /me behind auth.
// Frontend calls GET /me after login: admin -> /admin, user -> /dashboard.
func MountMe(r chi.Router, h *handler.Me, auth *middleware.Auth) {
	r.With(auth.Middleware).Get("/me", h.Get)
	r.With(auth.Middleware).Patch("/me", h.Patch)
}
