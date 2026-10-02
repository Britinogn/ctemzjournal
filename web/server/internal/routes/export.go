package routes

import (
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// mountExportCSV registers GET /trades/export.csv inside the /trades group
// (called by MountTrades: chi panics on two Route("/trades") blocks, and the
// static segment wins over /{id}). Exports are rate-limited (heavy query).
func mountExportCSV(r chi.Router, h *handler.Export) {
	r.With(middleware.ExportLimiter).Get("/export.csv", h.CSV)
}
