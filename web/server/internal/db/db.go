package db

import (
	"context"
	"os"
	"strconv"
	"time"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB bundles the pgx pool with the sqlc Queries.
type DB struct {
	Pool    *pgxpool.Pool
	Queries *sqlc.Queries
}

const (
	defaultMaxConns = 5 // Supabase session pooler allows 15 total; deploys briefly run two copies
	connectAttempts = 5
)

// maxConns reads DB_MAX_CONNS (e.g. set it to 2 for local dev), else the default.
func maxConns() int32 {
	if v, err := strconv.Atoi(os.Getenv("DB_MAX_CONNS")); err == nil && v > 0 {
		return int32(v)
	}
	return defaultMaxConns
}

// Connect opens a pgx pool against DATABASE_URL. If the pooler is full (for
// example while an old copy is still shutting down), it retries with a short
// backoff instead of failing the whole start.
func Connect(ctx context.Context, databaseURL string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = maxConns()
	cfg.MinConns = 0
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = time.Minute

	var pool *pgxpool.Pool
	for attempt := 1; attempt <= connectAttempts; attempt++ {
		pool, err = pgxpool.NewWithConfig(ctx, cfg)
		if err == nil {
			if err = pool.Ping(ctx); err == nil {
				break
			}
			pool.Close()
		}
		if attempt == connectAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(attempt) * 3 * time.Second):
		}
	}
	if err != nil {
		return nil, err
	}
	return &DB{Pool: pool, Queries: sqlc.New(pool)}, nil
}

func (d *DB) Close() {
	if d != nil && d.Pool != nil {
		d.Pool.Close()
	}
}