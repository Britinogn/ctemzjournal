package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
)

// MountAccounts registers /accounts CRUD behind auth.
func MountAccounts(r chi.Router, h *handler.Accounts, auth *middleware.Auth) {
	r.Route("/accounts", func(r chi.Router) {
		r.Use(auth.Middleware)
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Patch("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
