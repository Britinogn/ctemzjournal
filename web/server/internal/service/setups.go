package service

import (
	"context"
	"errors"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/google/uuid"
)

// Setups owns setup rules (unique name per user is enforced by the DB).
type Setups struct {
	repo *repository.Setups
}

func NewSetups(repo *repository.Setups) *Setups { return &Setups{repo: repo} }

type SetupCreate struct {
	Name         string
	Rules        *string
	Invalidation *string
}

func (s *Setups) Create(ctx context.Context, userID uuid.UUID, in SetupCreate) (sqlc.Setup, error) {
	if in.Name == "" {
		return sqlc.Setup{}, errors.New("name is required")
	}
	return s.repo.Create(ctx, userID, in.Name, in.Rules, in.Invalidation)
}

func (s *Setups) Get(ctx context.Context, id, userID uuid.UUID) (sqlc.Setup, error) {
	return s.repo.Get(ctx, id, userID)
}

func (s *Setups) List(ctx context.Context, userID uuid.UUID) ([]sqlc.Setup, error) {
	setups, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	if setups == nil {
		setups = []sqlc.Setup{}
	}
	return setups, nil
}

type SetupUpdate struct {
	Name         *string
	Rules        *string
	Invalidation *string
}

func (s *Setups) Update(ctx context.Context, id, userID uuid.UUID, in SetupUpdate) (sqlc.Setup, error) {
	return s.repo.Update(ctx, id, userID, repository.SetupUpdate(in))
}

func (s *Setups) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.Delete(ctx, id, userID)
}
