package handler

import (
	"net/http"

	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/response"
	"github.com/google/uuid"
)

// Stats serves the dashboard aggregates:
// summary, equity-curve, calendar, by-setup, by-pair.
type Stats struct {
	svc *service.Stats
}

func NewStats(svc *service.Stats) *Stats { return &Stats{svc: svc} }

// accountOf parses the optional ?account= filter.
func accountOf(r *http.Request) (*uuid.UUID, bool) {
	s := r.URL.Query().Get("account")
	if s == "" {
		return nil, true
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, false
	}
	return &id, true
}

func (h *Stats) Summary(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	accountID, ok := accountOf(r)
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid account")
		return
	}
	summary, err := h.svc.Summary(r.Context(), uid, accountID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "summary failed")
		return
	}
	response.JSON(w, http.StatusOK, summary)
}

func (h *Stats) EquityCurve(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	accountID, ok := accountOf(r)
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid account")
		return
	}
	points, err := h.svc.EquityCurve(r.Context(), uid, accountID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "equity curve failed")
		return
	}
	response.JSON(w, http.StatusOK, points)
}

func (h *Stats) Calendar(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	accountID, ok := accountOf(r)
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid account")
		return
	}
	days, err := h.svc.Calendar(r.Context(), uid, accountID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "calendar failed")
		return
	}
	response.JSON(w, http.StatusOK, days)
}

func (h *Stats) BySetup(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	accountID, ok := accountOf(r)
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid account")
		return
	}
	stats, err := h.svc.BySetup(r.Context(), uid, accountID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "by-setup failed")
		return
	}
	response.JSON(w, http.StatusOK, stats)
}

func (h *Stats) ByPair(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	accountID, ok := accountOf(r)
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid account")
		return
	}
	stats, err := h.svc.ByPair(r.Context(), uid, accountID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "by-pair failed")
		return
	}
	response.JSON(w, http.StatusOK, stats)
}
