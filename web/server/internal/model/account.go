package model

import sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"

// Account mirrors the accounts row returned by sqlc.
type Account = sqlc.Account

// Account types (match the DB check constraint).
const (
	AccountDemo = "demo"
	AccountLive = "live"
)
