package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
)

// DashboardDeps carries the handlers mounted under the authenticated user group.
// Later phases add accounts/setups/tags/trades/images/stats/export here;
// the group itself stays the single "user dashboard" entry point.
type DashboardDeps struct {
	Me   *handler.Me
	Auth *middleware.Auth
}

// MountDashboard registers everything an authenticated user needs.
// Today: /me (+ /auth/sync). Phase 2+ adds the resource routes here.
func MountDashboard(r chi.Router, d DashboardDeps, authHandler *handler.Auth) {
	MountAuth(r, authHandler, d.Auth)
	MountMe(r, d.Me, d.Auth)
}
