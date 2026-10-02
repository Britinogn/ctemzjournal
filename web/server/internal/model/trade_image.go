package model

import sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"

// TradeImage mirrors the trade_images row (only the Cloudinary public_id
// is stored, so the provider can change later).
type TradeImage = sqlc.TradeImage

// Image kinds (match the DB check constraint).
const (
	ImageEntry = "entry"
	ImageExit  = "exit"
)
