package handler

import (
	"errors"
	"net/http"

	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Images handles the upload flow for a trade:
// POST signature -> browser uploads to Cloudinary -> POST confirm.
// Reads always go through short-lived signed delivery URLs.
type Images struct {
	svc *service.Images
}

func NewImages(svc *service.Images) *Images { return &Images{svc: svc} }

func tradeIDOf(r *http.Request) (uuid.UUID, uuid.UUID, bool) {
	uid, ok := userIDOf(r)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	tid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return uuid.Nil, uuid.Nil, false
	}
	return uid, tid, true
}

// RequestSignature returns the signed Cloudinary payload (rate-limited route).
func (h *Images) RequestSignature(w http.ResponseWriter, r *http.Request) {
	uid, tid, ok := tradeIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	sig, err := h.svc.RequestSignature(r.Context(), uid, tid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "trade not found")
			return
		}
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, sig)
}

type confirmImageRequest struct {
	PublicID string  `json:"public_id"`
	Kind     *string `json:"kind"`
}

// Confirm saves the Cloudinary public_id after the direct browser upload.
func (h *Images) Confirm(w http.ResponseWriter, r *http.Request) {
	uid, tid, ok := tradeIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	var req confirmImageRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	img, err := h.svc.Confirm(r.Context(), uid, tid, req.PublicID, req.Kind)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "trade not found")
			return
		}
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, img)
}

func (h *Images) List(w http.ResponseWriter, r *http.Request) {
	uid, tid, ok := tradeIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	images, err := h.svc.List(r.Context(), uid, tid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "trade not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "list failed")
		return
	}
	response.JSON(w, http.StatusOK, images)
}

func (h *Images) Delete(w http.ResponseWriter, r *http.Request) {
	uid, tid, ok := tradeIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	imageID, err := uuid.Parse(chi.URLParam(r, "imageId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid image id")
		return
	}
	if err := h.svc.Delete(r.Context(), uid, tid, imageID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "image not found")
			return
		}
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
