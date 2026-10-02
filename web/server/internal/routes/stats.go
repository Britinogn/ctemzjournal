package routes

import (
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// MountStats registers the dashboard aggregates behind auth:
// /stats/summary, /stats/equity-curve, /stats/calendar,
// /stats/by-setup, /stats/by-pair (all accept ?account=).
func MountStats(r chi.Router, h *handler.Stats, auth *middleware.Auth) {
	r.Route("/stats", func(r chi.Router) {
		r.Use(auth.Middleware)
		r.Get("/summary", h.Summary)
		r.Get("/equity-curve", h.EquityCurve)
		r.Get("/calendar", h.Calendar)
		r.Get("/by-setup", h.BySetup)
		r.Get("/by-pair", h.ByPair)
	})
}
