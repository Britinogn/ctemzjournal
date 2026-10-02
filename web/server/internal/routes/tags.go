package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
)

// MountTags registers /tags CRUD behind auth.
func MountTags(r chi.Router, h *handler.Tags, auth *middleware.Auth) {
	r.Route("/tags", func(r chi.Router) {
		r.Use(auth.Middleware)
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Patch("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
