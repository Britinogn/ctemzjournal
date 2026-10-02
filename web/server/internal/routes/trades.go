package routes

import (
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// MountTrades registers /trades CRUD, PATCH /trades/{id}/visibility
// ("Make public" switch, off by default) and GET /trades/export.csv
// behind auth.
func MountTrades(r chi.Router, h *handler.Trades, exp *handler.Export, auth *middleware.Auth) {
	r.Route("/trades", func(r chi.Router) {
		r.Use(auth.Middleware)
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Patch("/{id}", h.Update)
		r.Patch("/{id}/visibility", h.SetVisibility)
		r.Delete("/{id}", h.Delete)
		mountExportCSV(r, exp)
	})
}
