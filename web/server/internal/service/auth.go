package service

import (
	"context"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/google/uuid"
)

// Auth bootstraps the profile row on first login.
type Auth struct {
	auth     *repository.Auth
	profiles *repository.Profiles
}

func NewAuth(auth *repository.Auth, profiles *repository.Profiles) *Auth {
	return &Auth{auth: auth, profiles: profiles}
}

// Sync ensures the profile exists (the handle_new_user trigger usually
// creates it) and returns it. timezone is only applied on insert — an
// existing custom timezone is never overwritten by a later login.
func (s *Auth) Sync(ctx context.Context, userID uuid.UUID, displayName, timezone string) (sqlc.Profile, error) {
	if err := s.auth.EnsureProfile(ctx, userID, displayName, normalizeTimezone(timezone)); err != nil {
		return sqlc.Profile{}, err
	}
	return s.profiles.GetByID(ctx, userID)
}
