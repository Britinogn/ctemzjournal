package integration

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/test/fixtures"
	"github.com/britinogn/ctemzjournal/test/helpers"
	"github.com/google/uuid"
)

func TestDashboardOverview(t *testing.T) {
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	ctx := context.Background()
	user := fixtures.MustUser(t, database, "alice")

	accountsSvc := service.NewAccounts(repository.NewAccounts(database.Queries))
	account, err := accountsSvc.Create(ctx, user, service.AccountCreate{Name: "Main", Type: "demo"})
	if err != nil {
		t.Fatalf("seed account: %v", err)
	}
	tradesRepo := repository.NewTrades(database.Queries)
	tradesSvc := service.NewTrades(
		tradesRepo,
		repository.NewTradeTags(database.Queries),
		repository.NewAccounts(database.Queries),
		repository.NewSetups(database.Queries),
		repository.NewTags(database.Queries),
		database.Queries,
	)

	raw, err := os.ReadFile("../fixtures/five_trades.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var fixture []fixtureTrade
	if err := json.Unmarshal(raw, &fixture); err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	for _, f := range fixture {
		closedAt, _ := time.Parse(time.RFC3339, f.ClosedAt)
		if _, err := tradesSvc.Create(ctx, user, service.TradeCreate{
			AccountID: account.ID, Pair: f.Pair, Direction: f.Direction,
			Entry: &f.Entry, StopLoss: &f.StopLoss, ExitPrice: &f.ExitPrice,
			LotSize: &f.LotSize, Commission: f.Commission, Swap: f.Swap,
			FollowedRules: &f.FollowedRules, Status: "closed", ClosedAt: &closedAt,
		}); err != nil {
			t.Fatalf("seed %s: %v", f.Name, err)
		}
	}
	// One open trade on top.
	entry, sl, lots := 1.0850, 1.0800, 1.0
	if _, err := tradesSvc.Create(ctx, user, service.TradeCreate{
		AccountID: account.ID, Pair: "EUR/USD", Direction: "long",
		Entry: &entry, StopLoss: &sl, LotSize: &lots, Status: "open",
	}); err != nil {
		t.Fatalf("seed open: %v", err)
	}

	profilesRepo := repository.NewProfiles(database.Queries)
	statsSvc := service.NewStats(
		repository.NewStats(database.Queries),
		profilesRepo,
		repository.NewSetups(database.Queries),
		tradesRepo,
	)
	dashboard := service.NewDashboard(profilesRepo, accountsSvc, statsSvc, tradesRepo)
	overview, err := dashboard.Overview(ctx, user, nil)
	if err != nil {
		t.Fatalf("overview: %v", err)
	}

	if overview.Me.ID != user {
		t.Fatalf("me mismatch: %v", overview.Me.ID)
	}
	if len(overview.Accounts) != 1 {
		t.Fatalf("accounts = %d, want 1", len(overview.Accounts))
	}
	if overview.Summary.TotalTrades != 6 || overview.Summary.OpenTrades != 1 {
		t.Fatalf("totals: %+v", overview.Summary)
	}
	if !closeEnough(overview.Summary.WinRate, 0.6) {
		t.Fatalf("win rate: %+v", overview.Summary)
	}
	if len(overview.Equity) != 5 || len(overview.Calendar) != 5 {
		t.Fatalf("equity/calendar: %d/%d", len(overview.Equity), len(overview.Calendar))
	}
	if len(overview.ByPair) != 4 {
		t.Fatalf("by-pair: %d", len(overview.ByPair))
	}
	if len(overview.Recent) != 6 {
		t.Fatalf("recent = %d, want 6", len(overview.Recent))
	}

	// Account-scoped overview stays consistent.
	scoped, err := dashboard.Overview(ctx, user, &account.ID)
	if err != nil {
		t.Fatalf("scoped overview: %v", err)
	}
	if scoped.Summary.TotalTrades != 6 || len(scoped.Recent) != 6 {
		t.Fatalf("scoped totals: %+v (%d)", scoped.Summary, len(scoped.Recent))
	}
	other := uuid.New()
	empty, err := dashboard.Overview(ctx, user, &other)
	if err != nil {
		t.Fatalf("foreign overview: %v", err)
	}
	if empty.Summary.TotalTrades != 0 || len(empty.Recent) != 0 {
		t.Fatalf("foreign account leaked: %+v", empty.Summary)
	}
}
