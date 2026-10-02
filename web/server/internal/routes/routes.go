package routes

import (
	"github.com/britinogn/ctemzjournal/config"
	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Deps wires everything the router needs.
type Deps struct {
	Config  *config.Config
	Pool    *pgxpool.Pool
	Queries *sqlc.Queries
	Auth    *middleware.Auth // nil when SUPABASE_JWKS_URL is unset (local boot)
}

// New builds the chi router: CORS, logging, health, then the user dashboard
// group (/me, /auth/sync). Public/admin/resource routes mount in later phases.
func New(d Deps) chi.Router {
	r := chi.NewRouter()
	r.Use(chiMW.RequestID)
	r.Use(chiMW.RealIP)
	r.Use(chiMW.Logger)
	r.Use(chiMW.Recoverer)
	r.Use(middleware.CORS(d.Config.AllowedOrigins()))

	MountHealth(r)

	if d.Auth == nil {
		return r
	}

	profiles := repository.NewProfiles(d.Queries)
	authRepo := repository.NewAuth(d.Pool)

	meHandler := handler.NewMe(service.NewMe(profiles))
	authHandler := handler.NewAuth(service.NewAuth(authRepo, profiles))
	accountsHandler := handler.NewAccounts(service.NewAccounts(repository.NewAccounts(d.Queries)))
	setupsHandler := handler.NewSetups(service.NewSetups(repository.NewSetups(d.Queries)))
	tagsHandler := handler.NewTags(service.NewTags(repository.NewTags(d.Queries)))
	tradesRepo := repository.NewTrades(d.Queries)
	tradesHandler := handler.NewTrades(service.NewTrades(
		tradesRepo,
		repository.NewTradeTags(d.Queries),
		repository.NewAccounts(d.Queries),
		repository.NewSetups(d.Queries),
		repository.NewTags(d.Queries),
		d.Queries,
	))

	MountDashboard(r, DashboardDeps{
		Me:       meHandler,
		Accounts: accountsHandler,
		Setups:   setupsHandler,
		Tags:     tagsHandler,
		Trades:   tradesHandler,
		Auth:     d.Auth,
	}, authHandler)

	return r
}
