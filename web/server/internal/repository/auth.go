package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Auth handles profile bootstrap on first login.
// (Supabase Auth owns signup/login; the handle_new_user trigger usually
// creates the profile — this is the fallback used by POST /auth/sync.)
type Auth struct {
	pool *pgxpool.Pool
}

func NewAuth(pool *pgxpool.Pool) *Auth { return &Auth{pool: pool} }

// EnsureProfile inserts the profile row if missing and returns its id.
func (a *Auth) EnsureProfile(ctx context.Context, id uuid.UUID, displayName string) error {
	_, err := a.pool.Exec(ctx,
		`insert into profiles (id, display_name) values ($1, $2) on conflict (id) do nothing`,
		id, displayName)
	if err == pgx.ErrNoRows {
		return nil
	}
	return err
}
