package model

import sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"

// Tag mirrors the tags row returned by sqlc.
type Tag = sqlc.Tag

// Tag kinds (match the DB check constraint).
const (
	TagMistake = "mistake"
	TagGeneral = "general"
)
