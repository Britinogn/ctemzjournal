package handler

import (
	"net/http"

	"github.com/google/uuid"
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
}

func (h *Auth) Sync(w http.ResponseWriter, r *http.Request) {
	u := middleware.FromContext(r.Context())
	if u == nil {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	var req syncRequest
	if r.Body != nil {
		_ = response.Decode(r, &req)
	}
	uid, err := uuid.Parse(u.ID)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid user")
		return
	}
	profile, err := h.svc.Sync(r.Context(), uid, req.DisplayName)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "sync failed")
		return
	}
	response.JSON(w, http.StatusOK, profile)
}
