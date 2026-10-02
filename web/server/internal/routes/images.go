package routes

import (
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// MountImages registers the trade image flow behind auth:
// POST /trades/{id}/images/upload-signature (rate-limited),
// POST /trades/{id}/images (confirm),
// GET  /trades/{id}/images,
// DELETE /trades/{id}/images/{imageId}.
func MountImages(r chi.Router, h *handler.Images, auth *middleware.Auth) {
	r.Route("/trades/{id}/images", func(r chi.Router) {
		r.Use(auth.Middleware)
		r.With(middleware.UploadLimiter).Post("/upload-signature", h.RequestSignature)
		r.Post("/", h.Confirm)
		r.Get("/", h.List)
		r.Delete("/{imageId}", h.Delete)
	})
}
