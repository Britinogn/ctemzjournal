package integration

import (
	"context"
	"errors"
	"testing"

	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/test/fixtures"
	"github.com/britinogn/ctemzjournal/test/helpers"
	"github.com/jackc/pgx/v5"
)

func TestSetupsCRUDAndIsolation(t *testing.T) {
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	userA := fixtures.MustUser(t, database, "alice")
	userB := fixtures.MustUser(t, database, "bob")
	svc := service.NewSetups(repository.NewSetups(database.Queries))
	ctx := context.Background()

	rules := "only A+ setups"
	created, err := svc.Create(ctx, userA, service.SetupCreate{Name: "Breakout", Rules: &rules})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Get(ctx, created.ID, userB); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("user B read user A setup: %v", err)
	}

	list, err := svc.List(ctx, userA)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %v (n=%d)", err, len(list))
	}

	name := "Breakout v2"
	updated, err := svc.Update(ctx, created.ID, userA, service.SetupUpdate{Name: &name})
	if err != nil || updated.Name != "Breakout v2" {
		t.Fatalf("update: %v (%+v)", err, updated)
	}

	if _, err := svc.Create(ctx, userA, service.SetupCreate{}); err == nil {
		t.Fatalf("expected name-required error")
	}

	if err := svc.Delete(ctx, created.ID, userA); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
