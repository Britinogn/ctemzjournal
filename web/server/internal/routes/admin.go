package routes

import (
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// MountAdmin registers the admin area behind auth + the admin role gate.
// Reads are open to admins; every admin write is recorded in audit_log.
func MountAdmin(r chi.Router, h *handler.Admin, auth *middleware.Auth) {
	r.Route("/admin", func(r chi.Router) {
		r.Use(auth.Middleware)
		r.Use(middleware.RequireRole("admin"))
		r.Get("/overview", h.Overview)
		r.Get("/users", h.ListUsers)
		r.Patch("/users/{id}/status", h.SetUserStatus)
		r.Get("/journals", h.ListJournals)
		r.Patch("/journals/{id}/hide", h.HideJournal)
		r.Get("/site-settings", h.GetSettings)
		r.Patch("/site-settings", h.UpdateSettings)
		r.Post("/site-settings/logo", h.LogoTarget)
		r.Get("/audit", h.ListAudit)
	})
}
