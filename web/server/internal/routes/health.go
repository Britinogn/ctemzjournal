package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/britinogn/ctemzjournal/internal/handler"
)

// MountHealth registers GET /healthz (no auth).
func MountHealth(r chi.Router) {
	r.Get("/healthz", handler.Health)
}
