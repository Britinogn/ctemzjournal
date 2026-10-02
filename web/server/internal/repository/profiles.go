package repository

import (
	"context"
	"time"

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

// Count returns the total user count (admin overview).
func (p *Profiles) Count(ctx context.Context) (int64, error) {
	return p.q.CountUsers(ctx)
}

// CountNewSince returns users created at or after since (admin overview).
func (p *Profiles) CountNewSince(ctx context.Context, since time.Time) (int64, error) {
	return p.q.CountNewUsersSince(ctx, pgtype.Timestamptz{Time: since, Valid: true})
}

// ListUsers returns profiles with display-name search (admin only).
func (p *Profiles) ListUsers(ctx context.Context, search string, limit, offset int32) ([]sqlc.Profile, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	users, err := p.q.ListUsers(ctx, sqlc.ListUsersParams{Column1: search, Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	if users == nil {
		users = []sqlc.Profile{}
	}
	return users, nil
}
