package model

import sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"

// Profile mirrors the profiles row returned by sqlc.
type Profile = sqlc.Profile

// Roles and statuses (match the DB check constraints).
const (
	RoleUser  = "user"
	RoleAdmin = "admin"

	StatusActive    = "active"
	StatusSuspended = "suspended"
)
