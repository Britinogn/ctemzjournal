package middleware

import (
	"net/http"

	"github.com/MicahParks/keyfunc/v3"
	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/pkg/response"
)

// Maintenance locks non-admin traffic while site_settings.maintenance_mode
// is on. Mounted globally (before auth), so it must verify tokens itself.
// Exempt: /healthz (deploy probes) and GET /public/site-settings (clients
// discover the flag through it — without this the frontend could never
// render the maintenance screen). Everyone else, anonymous included,
// gets 503; admins pass. A settings-read failure fails OPEN (a DB hiccup
// must not lock the site by itself).
type Maintenance struct {
	queries *sqlc.Queries
	kf      keyfunc.Keyfunc
}

// NewMaintenance builds the gate from live dependencies.
func NewMaintenance(q *sqlc.Queries, kf keyfunc.Keyfunc) *Maintenance {
	return &Maintenance{queries: q, kf: kf}
}

// Keyfunc exposes the Auth verifier for wiring (routes package).
func (a *Auth) Keyfunc() keyfunc.Keyfunc { return a.kf }

func (m *Maintenance) RequireAvailable(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/public/site-settings" {
			next.ServeHTTP(w, r)
			return
		}
		settings, err := m.queries.GetSiteSettings(r.Context())
		if err != nil || !settings.MaintenanceMode {
			next.ServeHTTP(w, r)
			return
		}
		uid, role, ok := m.adminIdentity(r)
		if !ok {
			_ = uid
			_ = role
			response.Error(w, http.StatusServiceUnavailable, "The site is under maintenance. Please check back shortly.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// adminIdentity verifies the token and loads the DB role.
// Only admins pass during maintenance.
func (m *Maintenance) adminIdentity(r *http.Request) (uid, role string, ok bool) {
	auth := &Auth{queries: m.queries, kf: m.kf}
	id, _, valid := auth.verifyToken(r)
	if !valid {
		return "", "", false
	}
	profile, err := m.queries.GetProfileByID(r.Context(), id)
	if err != nil || profile.Role != "admin" {
		return "", "", false
	}
	return id.String(), profile.Role, true
}
