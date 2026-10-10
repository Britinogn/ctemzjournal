package handler

import (
	"net/http"

	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/response"
)

// Auth handles POST /auth/sync (profile bootstrap on first login).
type Auth struct {
	svc *service.Auth
}

func NewAuth(svc *service.Auth) *Auth { return &Auth{svc: svc} }

type syncRequest struct {
	DisplayName string `json:"display_name"`
	Timezone    string `json:"timezone"`
}

func (h *Auth) Sync(w http.ResponseWriter, r *http.Request) {
	// Token-only identity: the profile may not exist yet — creating it
	// is this endpoint's whole job.
	uid, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "No session found. Please log in again.")
		return
	}
	var req syncRequest
	if r.Body != nil {
		_ = response.Decode(r, &req)
	}
	profile, err := h.svc.Sync(r.Context(), uid, req.DisplayName, req.Timezone)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Could not set up your account. Please try again.")
		return
	}
	response.JSON(w, http.StatusOK, profile)
}
