package service

import (
	"context"

	"github.com/google/uuid"
	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/repository"
)

// Me serves GET /me and PATCH /me.
// The frontend calls GET /me after login to redirect:
// admin -> /admin, user -> /dashboard.
type Me struct {
	profiles *repository.Profiles
}

func NewMe(profiles *repository.Profiles) *Me { return &Me{profiles: profiles} }

func (s *Me) Get(ctx context.Context, userID uuid.UUID) (sqlc.Profile, error) {
	return s.profiles.GetByID(ctx, userID)
}

type MeUpdate struct {
	DisplayName *string
	Timezone    *string
	AvatarPath  *string
}

func (s *Me) Update(ctx context.Context, userID uuid.UUID, in MeUpdate) (sqlc.Profile, error) {
	return s.profiles.Update(ctx, userID, in.DisplayName, in.Timezone, in.AvatarPath)
}
