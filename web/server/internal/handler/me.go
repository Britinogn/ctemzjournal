package handler

import (
	"net/http"

	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/response"
	"github.com/google/uuid"
)

// Me handles GET /me and PATCH /me.
type Me struct {
	svc *service.Me
}

func NewMe(svc *service.Me) *Me { return &Me{svc: svc} }

// Get returns the profile and role (frontend redirects on role).
func (h *Me) Get(w http.ResponseWriter, r *http.Request) {
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
	profile, err := h.svc.Get(r.Context(), uid)
	if err != nil {
		response.Error(w, http.StatusNotFound, "profile not found")
		return
	}
	response.JSON(w, http.StatusOK, profile)
}

type meUpdateRequest struct {
	DisplayName *string `json:"display_name"`
	Timezone    *string `json:"timezone"`
	AvatarPath  *string `json:"avatar_path"`
}

func (h *Me) Patch(w http.ResponseWriter, r *http.Request) {
	u := middleware.FromContext(r.Context())
	if u == nil {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	var req meUpdateRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	uid, err := uuid.Parse(u.ID)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid user")
		return
	}
	profile, err := h.svc.Update(r.Context(), uid, service.MeUpdate{
		DisplayName: req.DisplayName,
		Timezone:    req.Timezone,
		AvatarPath:  req.AvatarPath,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "update failed")
		return
	}
	response.JSON(w, http.StatusOK, profile)
}
