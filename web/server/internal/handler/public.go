package handler

import (
	"errors"
	"net/http"

	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/britinogn/ctemzjournal/pkg/pagination"
	"github.com/britinogn/ctemzjournal/pkg/response"
	"github.com/jackc/pgx/v5"
)

// Public serves the home-page endpoints (no auth).
type Public struct {
	svc *service.Public
}

func NewPublic(svc *service.Public) *Public { return &Public{svc: svc} }

func (h *Public) SiteSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.svc.SiteSettings(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "settings failed")
		return
	}
	response.JSON(w, http.StatusOK, settings)
}

func (h *Public) Journals(w http.ResponseWriter, r *http.Request) {
	page := pagination.FromRequest(r)
	journals, err := h.svc.Journals(r.Context(), page.Limit, page.Offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "journals failed")
		return
	}
	response.JSON(w, http.StatusOK, journals)
}

// JournalDetail serves GET /public/journals/{id} with notes + images.
// Private, missing or admin-hidden trades are 404 (existence undisclosed).
func (h *Public) JournalDetail(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusNotFound, "journal not found")
		return
	}
	journal, err := h.svc.Journal(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "journal not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "journal failed")
		return
	}
	response.JSON(w, http.StatusOK, journal)
}

func (h *Public) Rates(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, h.svc.Rates())
}
