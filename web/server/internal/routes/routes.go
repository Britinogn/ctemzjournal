package routes

import (
	"log"

	"github.com/britinogn/ctemzjournal/config"
	infra "github.com/britinogn/ctemzjournal/internal/cloudinary"
	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/handler"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/britinogn/ctemzjournal/internal/rates"
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
	Rates   *rates.Service   // may be nil in tests; /public/rates then serves empty+stale
}

// New builds the chi router: CORS, logging, health, public home-page routes,
// then (with DB + JWKS) the authenticated dashboard group and admin routes.
func New(d Deps) chi.Router {
	r := chi.NewRouter()
	r.Use(chiMW.RequestID)
	r.Use(chiMW.RealIP)
	r.Use(chiMW.Logger)
	r.Use(chiMW.Recoverer)
	r.Use(middleware.CORS(d.Config.AllowedOrigins()))

	MountHealth(r)

	cld, err := infra.New(d.Config.CloudinaryCloudName, d.Config.CloudinaryAPIKey, d.Config.CloudinaryAPISecret)
	if err != nil {
		log.Printf("cloudinary disabled: %v", err)
		cld = &infra.Client{}
	}

	if d.Queries != nil {
		MountPublic(r, handler.NewPublic(service.NewPublic(
			repository.NewSiteSettings(d.Queries),
			d.Queries,
			repository.NewTradeImages(d.Queries),
			cld,
			d.Rates,
			d.Config.SupabaseURL,
		)))
	}

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
	imagesRepo := repository.NewTradeImages(d.Queries)
	imagesSvc := service.NewImages(tradesRepo, imagesRepo, cld)
	tradesHandler := handler.NewTrades(service.NewTrades(
		tradesRepo,
		repository.NewTradeTags(d.Queries),
		repository.NewAccounts(d.Queries),
		repository.NewSetups(d.Queries),
		repository.NewTags(d.Queries),
		d.Queries,
	), imagesSvc)
	imagesHandler := handler.NewImages(imagesSvc)
	statsHandler := handler.NewStats(service.NewStats(
		repository.NewStats(d.Queries),
		profiles,
		repository.NewSetups(d.Queries),
		tradesRepo,
	))
	adminHandler := handler.NewAdmin(service.NewAdmin(profiles, tradesRepo))

	MountDashboard(r, DashboardDeps{
		Me:       meHandler,
		Accounts: accountsHandler,
		Setups:   setupsHandler,
		Tags:     tagsHandler,
		Trades:   tradesHandler,
		Images:   imagesHandler,
		Stats:    statsHandler,
		Auth:     d.Auth,
	}, authHandler)
	MountAdmin(r, adminHandler, d.Auth)

	return r
}
