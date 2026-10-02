package handler

import (
	"net/http"

	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/response"
)

// Admin handles the /admin/* endpoints (role-gated by middleware).
type Admin struct {
	svc *service.Admin
}

func NewAdmin(svc *service.Admin) *Admin { return &Admin{svc: svc} }

// Overview serves GET /admin/overview (user count, trade count, new users).
func (h *Admin) Overview(w http.ResponseWriter, r *http.Request) {
	overview, err := h.svc.Overview(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "overview failed")
		return
	}
	response.JSON(w, http.StatusOK, overview)
}
