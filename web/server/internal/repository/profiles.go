package repository

import (
	"context"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Profiles wraps the sqlc profile queries.
type Profiles struct {
	q *sqlc.Queries
}

func NewProfiles(q *sqlc.Queries) *Profiles { return &Profiles{q: q} }

func (p *Profiles) GetByID(ctx context.Context, id uuid.UUID) (sqlc.Profile, error) {
	return p.q.GetProfileByID(ctx, id)
}

// Update applies display_name / timezone / avatar_path patch.
func (p *Profiles) Update(ctx context.Context, id uuid.UUID, displayName, timezone, avatarPath *string) (sqlc.Profile, error) {
	toText := func(s *string) pgtype.Text {
		if s == nil {
			return pgtype.Text{Valid: false}
		}
		return pgtype.Text{String: *s, Valid: true}
	}
	return p.q.UpdateProfile(ctx, sqlc.UpdateProfileParams{
		ID:          id,
		DisplayName: toText(displayName),
		Timezone:    toText(timezone),
		AvatarPath:  toText(avatarPath),
	})
}

// SetStatus suspends or reactivates a user (admin only).
func (p *Profiles) SetStatus(ctx context.Context, id uuid.UUID, status string) (sqlc.Profile, error) {
	return p.q.UpdateProfileStatus(ctx, sqlc.UpdateProfileStatusParams{ID: id, Status: status})
}
