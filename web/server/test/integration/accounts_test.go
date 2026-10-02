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

func newAccountsService(t *testing.T) (*service.Accounts, uuidPair) {
	t.Helper()
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	return service.NewAccounts(repository.NewAccounts(database.Queries)),
		uuidPair{A: fixtures.MustUser(t, database, "alice"), B: fixtures.MustUser(t, database, "bob")}
}

func TestAccountsCRUD(t *testing.T) {
	svc, users := newAccountsService(t)
	ctx := context.Background()

	created, err := svc.Create(ctx, users.A, service.AccountCreate{
		Name: "Demo One", Type: "demo", Currency: "USD", StartingBalance: 1000,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.Name != "Demo One" || created.Type != "demo" {
		t.Fatalf("unexpected account: %+v", created)
	}

	got, err := svc.Get(ctx, created.ID, users.A)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("id mismatch")
	}

	list, err := svc.List(ctx, users.A)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %v (n=%d)", err, len(list))
	}

	name := "Renamed"
	updated, err := svc.Update(ctx, created.ID, users.A, service.AccountUpdate{Name: &name})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Renamed" {
		t.Fatalf("name not updated: %s", updated.Name)
	}

	if err := svc.Delete(ctx, created.ID, users.A); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := svc.Get(ctx, created.ID, users.A); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestAccountsIsolation(t *testing.T) {
	svc, users := newAccountsService(t)
	ctx := context.Background()

	created, err := svc.Create(ctx, users.A, service.AccountCreate{
		Name: "Private", Type: "live", Currency: "USD",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, err := svc.Get(ctx, created.ID, users.B); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("user B read user A account: %v", err)
	}
	name := "hijack"
	if _, err := svc.Update(ctx, created.ID, users.B, service.AccountUpdate{Name: &name}); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("user B updated user A account: %v", err)
	}
	// Delete is idempotent; B's delete must not remove A's row.
	_ = svc.Delete(ctx, created.ID, users.B)
	if _, err := svc.Get(ctx, created.ID, users.A); err != nil {
		t.Fatalf("A's account gone after B's delete: %v", err)
	}
}

func TestAccountsValidation(t *testing.T) {
	svc, users := newAccountsService(t)
	ctx := context.Background()

	cases := []service.AccountCreate{
		{Name: "", Type: "demo"},
		{Name: "x", Type: "paper"},
		{Name: "x", Type: "demo", StartingBalance: -5},
	}
	for i, c := range cases {
		if _, err := svc.Create(ctx, users.A, c); err == nil {
			t.Fatalf("case %d: expected validation error", i)
		}
	}
}
