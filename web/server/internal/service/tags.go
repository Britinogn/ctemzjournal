package service

import (
	"context"
	"errors"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/model"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/google/uuid"
)

// Tags owns tag rules (kind must be mistake or general).
type Tags struct {
	repo *repository.Tags
}

func NewTags(repo *repository.Tags) *Tags { return &Tags{repo: repo} }

type TagCreate struct {
	Name string
	Kind string
}

func (s *Tags) Create(ctx context.Context, userID uuid.UUID, in TagCreate) (sqlc.Tag, error) {
	if in.Name == "" {
		return sqlc.Tag{}, errors.New("name is required")
	}
	if in.Kind == "" {
		in.Kind = model.TagGeneral
	}
	if in.Kind != model.TagMistake && in.Kind != model.TagGeneral {
		return sqlc.Tag{}, errors.New("kind must be mistake or general")
	}
	return s.repo.Create(ctx, userID, in.Name, in.Kind)
}

func (s *Tags) Get(ctx context.Context, id, userID uuid.UUID) (sqlc.Tag, error) {
	return s.repo.Get(ctx, id, userID)
}

func (s *Tags) List(ctx context.Context, userID uuid.UUID) ([]sqlc.Tag, error) {
	tags, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	if tags == nil {
		tags = []sqlc.Tag{}
	}
	return tags, nil
}

type TagUpdate struct {
	Name *string
	Kind *string
}

func (s *Tags) Update(ctx context.Context, id, userID uuid.UUID, in TagUpdate) (sqlc.Tag, error) {
	if in.Kind != nil && *in.Kind != model.TagMistake && *in.Kind != model.TagGeneral {
		return sqlc.Tag{}, errors.New("kind must be mistake or general")
	}
	return s.repo.Update(ctx, id, userID, repository.TagUpdate(in))
}

func (s *Tags) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.Delete(ctx, id, userID)
}
