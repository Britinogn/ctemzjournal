package integration

import (
	"context"
	"testing"

	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/test/fixtures"
	"github.com/britinogn/ctemzjournal/test/helpers"
)

func TestAdminOverview(t *testing.T) {
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	ctx := context.Background()

	alice := fixtures.MustUser(t, database, "alice")
	fixtures.MustUser(t, database, "bob")

	accounts := service.NewAccounts(repository.NewAccounts(database.Queries))
	account, err := accounts.Create(ctx, alice, service.AccountCreate{Name: "Main", Type: "demo"})
	if err != nil {
		t.Fatalf("seed account: %v", err)
	}
	trades := service.NewTrades(
		repository.NewTrades(database.Queries),
		repository.NewTradeTags(database.Queries),
		repository.NewAccounts(database.Queries),
		repository.NewSetups(database.Queries),
		repository.NewTags(database.Queries),
		database.Queries,
	)
	entry, sl, lots := 1.0850, 1.0800, 1.0
	for i := 0; i < 3; i++ {
		if _, err := trades.Create(ctx, alice, service.TradeCreate{
			AccountID: account.ID, Pair: "EUR/USD", Direction: "long",
			Entry: &entry, StopLoss: &sl, LotSize: &lots, Status: "open",
		}); err != nil {
			t.Fatalf("seed trade %d: %v", i, err)
		}
	}

	admin := service.NewAdmin(
		repository.NewProfiles(database.Queries),
		repository.NewTrades(database.Queries),
	)
	overview, err := admin.Overview(ctx)
	if err != nil {
		t.Fatalf("overview: %v", err)
	}
	if overview.UserCount != 2 {
		t.Fatalf("users = %d, want 2", overview.UserCount)
	}
	if overview.TradeCount != 3 {
		t.Fatalf("trades = %d, want 3", overview.TradeCount)
	}
	if overview.NewUsersThisWeek != 2 {
		t.Fatalf("new users = %d, want 2", overview.NewUsersThisWeek)
	}
}
