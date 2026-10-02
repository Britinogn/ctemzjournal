package repository

import (
	"context"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/google/uuid"
)

// TradeTags links trades to tags. Callers must verify the trade belongs
// to the user first (Trades.Get is user-scoped).
type TradeTags struct {
	q *sqlc.Queries
}

func NewTradeTags(q *sqlc.Queries) *TradeTags { return &TradeTags{q: q} }

func (t *TradeTags) ListByTrade(ctx context.Context, tradeID uuid.UUID) ([]sqlc.Tag, error) {
	tags, err := t.q.ListTagsByTrade(ctx, tradeID)
	if err != nil {
		return nil, err
	}
	if tags == nil {
		tags = []sqlc.Tag{}
	}
	return tags, nil
}

// SetTags replaces every tag on a trade (clear + add).
func (t *TradeTags) SetTags(ctx context.Context, tradeID uuid.UUID, tagIDs []uuid.UUID) error {
	if err := t.q.ClearTradeTags(ctx, tradeID); err != nil {
		return err
	}
	for _, tagID := range tagIDs {
		if err := t.q.AddTradeTag(ctx, sqlc.AddTradeTagParams{TradeID: tradeID, TagID: tagID}); err != nil {
			return err
		}
	}
	return nil
}
