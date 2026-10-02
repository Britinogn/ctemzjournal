package model

import sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"

// Trade mirrors the trades row returned by sqlc.
type Trade = sqlc.Trade

// Trade directions and statuses (match the DB check constraints).
const (
	DirectionLong  = "long"
	DirectionShort = "short"

	TradeOpen   = "open"
	TradeClosed = "closed"
)
