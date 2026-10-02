package repository

import (
	"context"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Stats wraps the stats queries. Money aggregates are computed in Go
// (pkg/calc, unit tested) from the closed-trade rows — the SQL sum/avg
// helpers stay unused so fractional pnls are never truncated.
type Stats struct {
	q *sqlc.Queries
}

func NewStats(q *sqlc.Queries) *Stats { return &Stats{q: q} }

// ClosedTrades returns closed trades oldest-first, optionally for one account.
func (s *Stats) ClosedTrades(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID) ([]sqlc.StatsClosedTradesRow, error) {
	var account pgtype.UUID
	if accountID != nil {
		account = pgtype.UUID{Bytes: *accountID, Valid: true}
	}
	rows, err := s.q.StatsClosedTrades(ctx, sqlc.StatsClosedTradesParams{UserID: userID, AccountID: account})
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []sqlc.StatsClosedTradesRow{}
	}
	return rows, nil
}
