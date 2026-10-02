package routes

import (
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// DashboardDeps carries the handlers mounted under the authenticated user group.
// This is the single "user dashboard" entry point: /me plus every user-owned
// resource (accounts, setups, tags now; trades/images/stats/export later).
type DashboardDeps struct {
	Me       *handler.Me
	Accounts *handler.Accounts
	Setups   *handler.Setups
	Tags     *handler.Tags
	Trades   *handler.Trades
	Images   *handler.Images
	Auth     *middleware.Auth
}

// MountDashboard registers everything an authenticated user needs.
func MountDashboard(r chi.Router, d DashboardDeps, authHandler *handler.Auth) {
	MountAuth(r, authHandler, d.Auth)
	MountMe(r, d.Me, d.Auth)
	MountAccounts(r, d.Accounts, d.Auth)
	MountSetups(r, d.Setups, d.Auth)
	MountTags(r, d.Tags, d.Auth)
	MountTrades(r, d.Trades, d.Auth)
	MountImages(r, d.Images, d.Auth)
}
