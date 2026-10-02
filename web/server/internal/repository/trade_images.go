package repository

import (
	"context"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// TradeImages wraps the sqlc image queries. Trade ownership is verified
// by callers via Trades.Get (user-scoped) before touching images.
type TradeImages struct {
	q *sqlc.Queries
}

func NewTradeImages(q *sqlc.Queries) *TradeImages { return &TradeImages{q: q} }

func (t *TradeImages) Create(ctx context.Context, tradeID uuid.UUID, publicID string, kind *string, position int) (sqlc.TradeImage, error) {
	var k pgtype.Text
	if kind != nil {
		k = pgtype.Text{String: *kind, Valid: true}
	}
	return t.q.CreateTradeImage(ctx, sqlc.CreateTradeImageParams{
		TradeID: tradeID, PublicID: publicID, Kind: k, Position: int32(position),
	})
}

func (t *TradeImages) ListByTrade(ctx context.Context, tradeID uuid.UUID) ([]sqlc.TradeImage, error) {
	images, err := t.q.ListImagesByTrade(ctx, tradeID)
	if err != nil {
		return nil, err
	}
	if images == nil {
		images = []sqlc.TradeImage{}
	}
	return images, nil
}

func (t *TradeImages) CountByTrade(ctx context.Context, tradeID uuid.UUID) (int64, error) {
	return t.q.CountImagesByTrade(ctx, tradeID)
}

func (t *TradeImages) Get(ctx context.Context, id uuid.UUID) (sqlc.TradeImage, error) {
	return t.q.GetTradeImage(ctx, id)
}

func (t *TradeImages) Delete(ctx context.Context, id uuid.UUID) error {
	return t.q.DeleteTradeImage(ctx, id)
}
