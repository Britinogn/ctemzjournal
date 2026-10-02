package service

import (
	"context"
	"errors"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/model"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/google/uuid"
)

// Accounts owns account business rules (type/currency validation, user scoping).
type Accounts struct {
	repo *repository.Accounts
}

func NewAccounts(repo *repository.Accounts) *Accounts { return &Accounts{repo: repo} }

type AccountCreate struct {
	Name            string  `validate:"required,max=100"`
	Type            string  `validate:"required,oneof=demo live"`
	Currency        string  `validate:"max=10"`
	StartingBalance float64 `validate:"gte=0"`
}

func (s *Accounts) Create(ctx context.Context, userID uuid.UUID, in AccountCreate) (sqlc.Account, error) {
	if in.Name == "" {
		return sqlc.Account{}, errors.New("name is required")
	}
	if in.Type != model.AccountDemo && in.Type != model.AccountLive {
		return sqlc.Account{}, errors.New("type must be demo or live")
	}
	if in.Currency == "" {
		in.Currency = "USD"
	}
	if in.StartingBalance < 0 {
		return sqlc.Account{}, errors.New("starting_balance must be >= 0")
	}
	return s.repo.Create(ctx, userID, in.Name, in.Type, in.Currency, in.StartingBalance)
}

func (s *Accounts) Get(ctx context.Context, id, userID uuid.UUID) (sqlc.Account, error) {
	return s.repo.Get(ctx, id, userID)
}

func (s *Accounts) List(ctx context.Context, userID uuid.UUID) ([]sqlc.Account, error) {
	accounts, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	if accounts == nil {
		accounts = []sqlc.Account{}
	}
	return accounts, nil
}

type AccountUpdate struct {
	Name            *string
	Type            *string
	Currency        *string
	StartingBalance *float64
}

func (s *Accounts) Update(ctx context.Context, id, userID uuid.UUID, in AccountUpdate) (sqlc.Account, error) {
	if in.Type != nil && *in.Type != model.AccountDemo && *in.Type != model.AccountLive {
		return sqlc.Account{}, errors.New("type must be demo or live")
	}
	if in.StartingBalance != nil && *in.StartingBalance < 0 {
		return sqlc.Account{}, errors.New("starting_balance must be >= 0")
	}
	return s.repo.Update(ctx, id, userID, repository.AccountUpdate(in))
}

func (s *Accounts) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.Delete(ctx, id, userID)
}
