package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/britinogn/ctemzjournal/config"
	"github.com/britinogn/ctemzjournal/internal/db"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/britinogn/ctemzjournal/internal/rates"
	"github.com/britinogn/ctemzjournal/internal/routes"
	"github.com/britinogn/ctemzjournal/internal/worker"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var (
		database *db.DB
		authMW   *middleware.Auth
	)
	if cfg.DatabaseURL == "" {
		log.Println("warn: DATABASE_URL unset — serving /healthz only")
	} else {
		var err error
		database, err = db.Connect(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("db connect: %v", err)
		}
		defer database.Close()
		log.Println("db connected")

		if cfg.SupabaseJWKSURL == "" {
			log.Println("warn: SUPABASE_JWKS_URL unset — auth routes disabled")
		} else {
			authMW, err = middleware.NewAuth(ctx, database.Queries, cfg.SupabaseJWKSURL)
			if err != nil {
				log.Fatalf("jwks: %v", err)
			}
			log.Println("jwks loaded")
		}
	}

	deps := routes.Deps{Config: cfg, Auth: authMW}
	if database != nil {
		deps.Pool = database.Pool
		deps.Queries = database.Queries
	}

	// Live prices: background goroutine fetches ~8 pairs on a timer into an
	// in-memory cache (browsers never call the provider). Refresh every
	// RATES_REFRESH_MINUTES (20 min default = ~576 Twelve Data req/day).
	ratesSvc := rates.NewService(
		rates.TwelveData{APIKey: cfg.TwelveDataAPIKey},
		rates.Frankfurter{},
		nil,
	)
	deps.Rates = ratesSvc
	workers := worker.NewGroup(ctx)
	interval := time.Duration(cfg.RatesRefreshMinutes) * time.Minute
	workers.Go(func(ctx context.Context) error { return worker.RunRates(ctx, ratesSvc, interval) })
	defer workers.Stop()

	r := routes.New(deps)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Printf("listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	workers.Stop()
	if err := workers.Err(); err != nil {
		log.Printf("worker error: %v", err)
	}
}
