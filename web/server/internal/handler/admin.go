package handler

import (
	"errors"
	"net/http"

	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/pagination"
	"github.com/britinogn/ctemzjournal/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Admin handles the /admin/* endpoints (role-gated by middleware).
type Admin struct {
	svc *service.Admin
}

func NewAdmin(svc *service.Admin) *Admin { return &Admin{svc: svc} }

func adminIDOf(r *http.Request) (uuid.UUID, bool) {
	u := middleware.FromContext(r.Context())
	if u == nil {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(u.ID)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

// Overview serves GET /admin/overview (user count, trade count, new users).
func (h *Admin) Overview(w http.ResponseWriter, r *http.Request) {
	overview, err := h.svc.Overview(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "overview failed")
		return
	}
	response.JSON(w, http.StatusOK, overview)
}

// ListUsers serves GET /admin/users?search=&limit=&offset=.
func (h *Admin) ListUsers(w http.ResponseWriter, r *http.Request) {
	page := pagination.FromRequest(r)
	users, err := h.svc.ListUsers(r.Context(), r.URL.Query().Get("search"), page.Limit, page.Offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "users failed")
		return
	}
	response.JSON(w, http.StatusOK, users)
}

type userStatusRequest struct {
	Status string `json:"status"`
}

// SetUserStatus serves PATCH /admin/users/{id}/status (suspend/reactivate).
func (h *Admin) SetUserStatus(w http.ResponseWriter, r *http.Request) {
	adminID, ok := adminIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	targetID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req userStatusRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	profile, err := h.svc.SetUserStatus(r.Context(), adminID, targetID, req.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, profile)
}

// ListJournals serves GET /admin/journals (public ones, for moderation).
func (h *Admin) ListJournals(w http.ResponseWriter, r *http.Request) {
	page := pagination.FromRequest(r)
	journals, err := h.svc.ListJournals(r.Context(), page.Limit, page.Offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "journals failed")
		return
	}
	response.JSON(w, http.StatusOK, journals)
}

type hideJournalRequest struct {
	Hidden bool `json:"hidden"`
}

// HideJournal serves PATCH /admin/journals/{id}/hide.
func (h *Admin) HideJournal(w http.ResponseWriter, r *http.Request) {
	adminID, ok := adminIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	tradeID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req hideJournalRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	trade, err := h.svc.HideJournal(r.Context(), adminID, tradeID, req.Hidden)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "journal not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "hide failed")
		return
	}
	response.JSON(w, http.StatusOK, trade)
}

// GetSettings serves GET /admin/site-settings.
func (h *Admin) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.svc.GetSettings(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "settings failed")
		return
	}
	response.JSON(w, http.StatusOK, settings)
}

type settingsUpdateRequest struct {
	SiteName        *string        `json:"site_name"`
	Tagline         *string        `json:"tagline"`
	LogoPath        *string        `json:"logo_path"`
	FaviconPath     *string        `json:"favicon_path"`
	ContactEmail    *string        `json:"contact_email"`
	FooterText      *string        `json:"footer_text"`
	RiskDisclaimer  *string        `json:"risk_disclaimer"`
	SocialLinks     map[string]any `json:"social_links"`
	AllowSignups    *bool          `json:"allow_signups"`
	MaintenanceMode *bool          `json:"maintenance_mode"`
}

// UpdateSettings serves PATCH /admin/site-settings.
func (h *Admin) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	adminID, ok := adminIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	var req settingsUpdateRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	in := service.SettingsUpdate{
		SiteName: req.SiteName, Tagline: req.Tagline, LogoPath: req.LogoPath,
		FaviconPath: req.FaviconPath, ContactEmail: req.ContactEmail, FooterText: req.FooterText,
		RiskDisclaimer: req.RiskDisclaimer, SocialLinks: req.SocialLinks,
		SocialLinksSet: req.SocialLinks != nil,
		AllowSignups:   req.AllowSignups, MaintenanceMode: req.MaintenanceMode,
	}
	settings, err := h.svc.UpdateSettings(r.Context(), adminID, in)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "settings update failed")
		return
	}
	response.JSON(w, http.StatusOK, settings)
}

type logoRequest struct {
	Kind      string `json:"kind"`
	Extension string `json:"extension"`
}

// LogoTarget serves POST /admin/site-settings/logo: reserves the
// site-assets path the browser uploads to itself (no service key here).
func (h *Admin) LogoTarget(w http.ResponseWriter, r *http.Request) {
	var req logoRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	target, err := h.svc.LogoTarget(r.Context(), req.Kind, req.Extension)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, target)
}

// ListAudit serves GET /admin/audit?limit=&offset= (action history).
func (h *Admin) ListAudit(w http.ResponseWriter, r *http.Request) {
	page := pagination.FromRequest(r)
	rows, err := h.svc.ListAudit(r.Context(), page.Limit, page.Offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "audit failed")
		return
	}
	response.JSON(w, http.StatusOK, rows)
}
