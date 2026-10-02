package repository

import (
	"context"

	"github.com/google/uuid"
	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
)

// Tags wraps the sqlc tag queries. Every method is scoped by user_id.
type Tags struct {
	q *sqlc.Queries
}

func NewTags(q *sqlc.Queries) *Tags { return &Tags{q: q} }

func (t *Tags) Create(ctx context.Context, userID uuid.UUID, name, kind string) (sqlc.Tag, error) {
	return t.q.CreateTag(ctx, sqlc.CreateTagParams{UserID: userID, Name: name, Kind: kind})
}

func (t *Tags) Get(ctx context.Context, id, userID uuid.UUID) (sqlc.Tag, error) {
	return t.q.GetTag(ctx, sqlc.GetTagParams{ID: id, UserID: userID})
}

func (t *Tags) List(ctx context.Context, userID uuid.UUID) ([]sqlc.Tag, error) {
	return t.q.ListTagsByUser(ctx, userID)
}

type TagUpdate struct {
	Name *string
	Kind *string
}

func (t *Tags) Update(ctx context.Context, id, userID uuid.UUID, in TagUpdate) (sqlc.Tag, error) {
	return t.q.UpdateTag(ctx, sqlc.UpdateTagParams{
		Name:   textArg(in.Name),
		Kind:   textArg(in.Kind),
		ID:     id,
		UserID: userID,
	})
}

func (t *Tags) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return t.q.DeleteTag(ctx, sqlc.DeleteTagParams{ID: id, UserID: userID})
}
