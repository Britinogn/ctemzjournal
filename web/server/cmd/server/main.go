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
	"github.com/britinogn/ctemzjournal/internal/routes"
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
}
