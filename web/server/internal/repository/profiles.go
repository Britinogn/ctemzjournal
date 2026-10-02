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

// toProfile maps the explicit-column profile rows (email lives only in
// ListUsersAdmin / the model) back to sqlc.Profile, keeping every
// existing caller unchanged.
func toProfile(
	id uuid.UUID, displayName pgtype.Text, role, status, timezone string,
	avatarPath pgtype.Text, createdAt, updatedAt pgtype.Timestamptz,
) sqlc.Profile {
	return sqlc.Profile{
		ID: id, DisplayName: displayName, Role: role, Status: status,
		Timezone: timezone, AvatarPath: avatarPath,
		CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
}

func (p *Profiles) GetByID(ctx context.Context, id uuid.UUID) (sqlc.Profile, error) {
	row, err := p.q.GetProfileByID(ctx, id)
	if err != nil {
		return sqlc.Profile{}, err
	}
	return toProfile(row.ID, row.DisplayName, row.Role, row.Status, row.Timezone, row.AvatarPath, row.CreatedAt, row.UpdatedAt), nil
}

// Update applies display_name / timezone / avatar_path patch.
func (p *Profiles) Update(ctx context.Context, id uuid.UUID, displayName, timezone, avatarPath *string) (sqlc.Profile, error) {
	toText := func(s *string) pgtype.Text {
		if s == nil {
			return pgtype.Text{Valid: false}
		}
		return pgtype.Text{String: *s, Valid: true}
	}
	row, err := p.q.UpdateProfile(ctx, sqlc.UpdateProfileParams{
		ID:          id,
		DisplayName: toText(displayName),
		Timezone:    toText(timezone),
		AvatarPath:  toText(avatarPath),
	})
	if err != nil {
		return sqlc.Profile{}, err
	}
	return toProfile(row.ID, row.DisplayName, row.Role, row.Status, row.Timezone, row.AvatarPath, row.CreatedAt, row.UpdatedAt), nil
}

// SetStatus suspends or reactivates a user (admin only).
func (p *Profiles) SetStatus(ctx context.Context, id uuid.UUID, status string) (sqlc.Profile, error) {
	row, err := p.q.UpdateProfileStatus(ctx, sqlc.UpdateProfileStatusParams{ID: id, Status: status})
	if err != nil {
		return sqlc.Profile{}, err
	}
	return toProfile(row.ID, row.DisplayName, row.Role, row.Status, row.Timezone, row.AvatarPath, row.CreatedAt, row.UpdatedAt), nil
}

// Count returns the total user count (admin overview).
func (p *Profiles) Count(ctx context.Context) (int64, error) {
	return p.q.CountUsers(ctx)
}

// CountNewSince returns users created at or after since (admin overview).
func (p *Profiles) CountNewSince(ctx context.Context, since time.Time) (int64, error) {
	return p.q.CountNewUsersSince(ctx, pgtype.Timestamptz{Time: since, Valid: true})
}

// CountByStatus returns users with the given status (active/suspended split).
func (p *Profiles) CountByStatus(ctx context.Context, status string) (int64, error) {
	return p.q.CountUsersByStatus(ctx, status)
}

// SignupsSince returns per-day signup counts at/after since (signup chart).
func (p *Profiles) SignupsSince(ctx context.Context, since time.Time) ([]sqlc.SignupsByDayRow, error) {
	rows, err := p.q.SignupsByDay(ctx, pgtype.Timestamptz{Time: since, Valid: true})
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []sqlc.SignupsByDayRow{}
	}
	return rows, nil
}

// ListUsersAdmin returns profiles with email + trade counts and
// name-or-email search (admin users table).
func (p *Profiles) ListUsersAdmin(ctx context.Context, search string, limit, offset int32) ([]sqlc.ListUsersAdminRow, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := p.q.ListUsersAdmin(ctx, sqlc.ListUsersAdminParams{Search: search, PageOffset: offset, PageLimit: limit})
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []sqlc.ListUsersAdminRow{}
	}
	return rows, nil
}

// ListUsers returns profiles with display-name search (admin only).
func (p *Profiles) ListUsers(ctx context.Context, search string, limit, offset int32) ([]sqlc.Profile, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := p.q.ListUsers(ctx, sqlc.ListUsersParams{Column1: search, Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	users := make([]sqlc.Profile, 0, len(rows))
	for _, row := range rows {
		users = append(users, toProfile(row.ID, row.DisplayName, row.Role, row.Status, row.Timezone, row.AvatarPath, row.CreatedAt, row.UpdatedAt))
	}
	return users, nil
}
