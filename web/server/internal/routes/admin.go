package routes

import (
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// MountAdmin registers the admin area behind auth + the admin role gate.
// Reads (overview) are open to admins; every admin write (later phases)
// is recorded in audit_log.
func MountAdmin(r chi.Router, h *handler.Admin, auth *middleware.Auth) {
	r.Route("/admin", func(r chi.Router) {
		r.Use(auth.Middleware)
		r.Use(middleware.RequireRole("admin"))
		r.Get("/overview", h.Overview)
	})
}
