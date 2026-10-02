package handler

import (
	"net/http"

	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/response"
	"github.com/google/uuid"
)

// Dashboard serves GET /dashboard: one overview of the user's account.
type Dashboard struct {
	svc *service.Dashboard
}

func NewDashboard(svc *service.Dashboard) *Dashboard { return &Dashboard{svc: svc} }

// Get returns profile, accounts, headline totals, equity, calendar,
// breakdowns and recent trades. Accepts the same ?account= filter as /stats/*.
func (h *Dashboard) Get(w http.ResponseWriter, r *http.Request) {
	u := middleware.FromContext(r.Context())
	if u == nil {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	uid, err := uuid.Parse(u.ID)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid user")
		return
	}
	var accountID *uuid.UUID
	if s := r.URL.Query().Get("account"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid account")
			return
		}
		accountID = &id
	}
	overview, err := h.svc.Overview(r.Context(), uid, accountID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "dashboard failed")
		return
	}
	response.JSON(w, http.StatusOK, overview)
}
