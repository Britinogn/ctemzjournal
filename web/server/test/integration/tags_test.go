package integration

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/test/fixtures"
	"github.com/britinogn/ctemzjournal/test/helpers"
)

func TestTagsCRUDAndIsolation(t *testing.T) {
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	userA := fixtures.MustUser(t, database, "alice")
	userB := fixtures.MustUser(t, database, "bob")
	svc := service.NewTags(repository.NewTags(database.Queries))
	ctx := context.Background()

	created, err := svc.Create(ctx, userA, service.TagCreate{Name: "FOMO", Kind: "mistake"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.Kind != "mistake" {
		t.Fatalf("kind not kept: %s", created.Kind)
	}
	if _, err := svc.Get(ctx, created.ID, userB); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("user B read user A tag: %v", err)
	}

	// empty kind defaults to general
	def, err := svc.Create(ctx, userA, service.TagCreate{Name: "London"})
	if err != nil || def.Kind != "general" {
		t.Fatalf("default kind: %v (%+v)", err, def)
	}

	if _, err := svc.Create(ctx, userA, service.TagCreate{Name: "x", Kind: "nope"}); err == nil {
		t.Fatalf("expected kind validation error")
	}

	kind := "general"
	updated, err := svc.Update(ctx, created.ID, userA, service.TagUpdate{Kind: &kind})
	if err != nil || updated.Kind != "general" {
		t.Fatalf("update: %v (%+v)", err, updated)
	}

	if err := svc.Delete(ctx, created.ID, userA); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
