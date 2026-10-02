package repository

import (
	"context"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// AuditLog records every admin write (docx §6).
type AuditLog struct {
	q *sqlc.Queries
}

func NewAuditLog(q *sqlc.Queries) *AuditLog { return &AuditLog{q: q} }

// Record writes one audit row. meta is pre-marshalled JSON (nil keeps '{}').
func (a *AuditLog) Record(ctx context.Context, adminID *uuid.UUID, action, targetType, targetID string, meta []byte) error {
	var admin pgtype.UUID
	if adminID != nil {
		admin = pgtype.UUID{Bytes: *adminID, Valid: true}
	}
	if meta == nil {
		meta = []byte("{}")
	}
	_, err := a.q.CreateAuditLog(ctx, sqlc.CreateAuditLogParams{
		AdminID: admin, Action: action, TargetType: targetType, TargetID: targetID, Meta: meta,
	})
	return err
}

// List returns audit rows newest-first.
func (a *AuditLog) List(ctx context.Context, limit, offset int32) ([]sqlc.AuditLog, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := a.q.ListAuditLog(ctx, sqlc.ListAuditLogParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []sqlc.AuditLog{}
	}
	return rows, nil
}
