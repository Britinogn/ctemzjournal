package routes

import (
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/go-chi/chi/v5"
)

// MountPublic registers the home-page endpoints (no auth):
// /public/site-settings, /public/journals, /public/rates.
func MountPublic(r chi.Router, h *handler.Public) {
	r.Route("/public", func(r chi.Router) {
		r.Get("/site-settings", h.SiteSettings)
		r.Get("/journals", h.Journals)
		r.Get("/rates", h.Rates)
	})
}
