package routes

import (
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/go-chi/chi/v5"
)

// MountHealth registers GET /healthz (no auth).
func MountHealth(r chi.Router) {
	r.Get("/healthz", handler.Health)
}
