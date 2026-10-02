package repository

import (
	"context"

	"github.com/google/uuid"
	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
)

// Setups wraps the sqlc setup queries. Every method is scoped by user_id.
type Setups struct {
	q *sqlc.Queries
}

func NewSetups(q *sqlc.Queries) *Setups { return &Setups{q: q} }

func (s *Setups) Create(ctx context.Context, userID uuid.UUID, name string, rules, invalidation *string) (sqlc.Setup, error) {
	return s.q.CreateSetup(ctx, sqlc.CreateSetupParams{
		UserID:       userID,
		Name:         name,
		Rules:        textArg(rules),
		Invalidation: textArg(invalidation),
	})
}

func (s *Setups) Get(ctx context.Context, id, userID uuid.UUID) (sqlc.Setup, error) {
	return s.q.GetSetup(ctx, sqlc.GetSetupParams{ID: id, UserID: userID})
}

func (s *Setups) List(ctx context.Context, userID uuid.UUID) ([]sqlc.Setup, error) {
	return s.q.ListSetupsByUser(ctx, userID)
}

type SetupUpdate struct {
	Name         *string
	Rules        *string
	Invalidation *string
}

func (s *Setups) Update(ctx context.Context, id, userID uuid.UUID, in SetupUpdate) (sqlc.Setup, error) {
	return s.q.UpdateSetup(ctx, sqlc.UpdateSetupParams{
		Name:         textArg(in.Name),
		Rules:        textArg(in.Rules),
		Invalidation: textArg(in.Invalidation),
		ID:           id,
		UserID:       userID,
	})
}

func (s *Setups) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.q.DeleteSetup(ctx, sqlc.DeleteSetupParams{ID: id, UserID: userID})
}
