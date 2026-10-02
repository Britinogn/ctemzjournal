package repository

import (
	"context"
	"strconv"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Accounts wraps the sqlc account queries. Every method is scoped by user_id.
type Accounts struct {
	q *sqlc.Queries
}

func NewAccounts(q *sqlc.Queries) *Accounts { return &Accounts{q: q} }

func textArg(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func numericArg(f *float64) pgtype.Numeric {
	var n pgtype.Numeric
	if f == nil {
		return n
	}
	_ = n.Scan(strconv.FormatFloat(*f, 'f', 2, 64))
	return n
}

func (a *Accounts) Create(ctx context.Context, userID uuid.UUID, name, typ, currency string, startingBalance float64) (sqlc.Account, error) {
	return a.q.CreateAccount(ctx, sqlc.CreateAccountParams{
		UserID:          userID,
		Name:            name,
		Type:            typ,
		Currency:        currency,
		StartingBalance: numericArg(&startingBalance),
	})
}

func (a *Accounts) Get(ctx context.Context, id, userID uuid.UUID) (sqlc.Account, error) {
	return a.q.GetAccount(ctx, sqlc.GetAccountParams{ID: id, UserID: userID})
}

func (a *Accounts) List(ctx context.Context, userID uuid.UUID) ([]sqlc.Account, error) {
	return a.q.ListAccountsByUser(ctx, userID)
}

type AccountUpdate struct {
	Name            *string
	Type            *string
	Currency        *string
	StartingBalance *float64
}

func (a *Accounts) Update(ctx context.Context, id, userID uuid.UUID, in AccountUpdate) (sqlc.Account, error) {
	return a.q.UpdateAccount(ctx, sqlc.UpdateAccountParams{
		Name:            textArg(in.Name),
		Type:            textArg(in.Type),
		Currency:        textArg(in.Currency),
		StartingBalance: numericArg(in.StartingBalance),
		ID:              id,
		UserID:          userID,
	})
}

func (a *Accounts) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return a.q.DeleteAccount(ctx, sqlc.DeleteAccountParams{ID: id, UserID: userID})
}
