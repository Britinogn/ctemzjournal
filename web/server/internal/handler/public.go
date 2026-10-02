package handler

import (
	"net/http"

	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/pagination"
	"github.com/britinogn/ctemzjournal/pkg/response"
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

func (h *Public) Rates(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, h.svc.Rates())
}
