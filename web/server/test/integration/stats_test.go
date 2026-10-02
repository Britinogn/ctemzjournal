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

type statsEnv struct {
	stats   *service.Stats
	trades  *service.Trades
	user    uuid.UUID
	account uuid.UUID
}

func newStatsEnv(t *testing.T) (*statsEnv, context.Context) {
	t.Helper()
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	user := fixtures.MustUser(t, database, "alice")
	accounts := service.NewAccounts(repository.NewAccounts(database.Queries))
	ctx := context.Background()
	account, err := accounts.Create(ctx, user, service.AccountCreate{Name: "Main", Type: "demo"})
	if err != nil {
		t.Fatalf("seed account: %v", err)
	}
	tradesSvc := service.NewTrades(
		repository.NewTrades(database.Queries),
		repository.NewTradeTags(database.Queries),
		repository.NewAccounts(database.Queries),
		repository.NewSetups(database.Queries),
		repository.NewTags(database.Queries),
		database.Queries,
	)
	statsSvc := service.NewStats(
		repository.NewStats(database.Queries),
		repository.NewProfiles(database.Queries),
		repository.NewSetups(database.Queries),
	)
	return &statsEnv{stats: statsSvc, trades: tradesSvc, user: user, account: account.ID}, ctx
}

func seedFiveTrades(t *testing.T, env *statsEnv, ctx context.Context) {
	t.Helper()
	raw, err := os.ReadFile("../fixtures/five_trades.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var trades []fixtureTrade
	if err := json.Unmarshal(raw, &trades); err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	for _, f := range trades {
		closedAt, _ := time.Parse(time.RFC3339, f.ClosedAt)
		if _, err := env.trades.Create(ctx, env.user, service.TradeCreate{
			AccountID: env.account, Pair: f.Pair, Direction: f.Direction,
			Entry: &f.Entry, StopLoss: &f.StopLoss, ExitPrice: &f.ExitPrice,
			LotSize: &f.LotSize, Commission: f.Commission, Swap: f.Swap,
			FollowedRules: &f.FollowedRules, Status: "closed", ClosedAt: &closedAt,
		}); err != nil {
			t.Fatalf("seed %s: %v", f.Name, err)
		}
	}
}

func TestStatsSummary(t *testing.T) {
	env, ctx := newStatsEnv(t)
	seedFiveTrades(t, env, ctx)

	s, err := env.stats.Summary(ctx, env.user, nil)
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if s.Total != 5 || s.Wins != 3 || s.Losses != 2 {
		t.Fatalf("counts: %+v", s)
	}
	if !closeEnough(s.WinRate, 0.6) || !closeEnough(s.AvgR, 0.8336) ||
		!closeEnough(s.Expectancy, 15185.6) || !closeEnough(s.MaxDrawdown, 1312) ||
		!closeEnough(s.RuleRate, 0.6) {
		t.Fatalf("summary mismatch: %+v", s)
	}

	// Unknown account filter matches nothing.
	other := uuid.New()
	empty, err := env.stats.Summary(ctx, env.user, &other)
	if err != nil || empty.Total != 0 {
		t.Fatalf("foreign account filter: %+v (%v)", empty, err)
	}
}

func TestStatsEquityAndCalendar(t *testing.T) {
	env, ctx := newStatsEnv(t)
	seedFiveTrades(t, env, ctx)

	points, err := env.stats.EquityCurve(ctx, env.user, nil)
	if err != nil || len(points) != 5 {
		t.Fatalf("equity: %v (n=%d)", err, len(points))
	}
	wantEquity := []float64{990, 75990, 75178, 74678, 75928}
	for i, w := range wantEquity {
		if !closeEnough(points[i].Equity, w) {
			t.Fatalf("point %d = %+v, want equity %v", i, points[i], w)
		}
		if points[i].Date == "" {
			t.Fatalf("point %d missing date", i)
		}
	}

	days, err := env.stats.Calendar(ctx, env.user, nil)
	if err != nil || len(days) != 5 {
		t.Fatalf("calendar: %v (n=%d)", err, len(days))
	}
	if days[0].Date != "2026-09-01" || days[0].Trades != 1 || days[0].Wins != 1 {
		t.Fatalf("day 1: %+v", days[0])
	}
}

func TestStatsByPairAndSetup(t *testing.T) {
	env, ctx := newStatsEnv(t)
	seedFiveTrades(t, env, ctx)

	pairs, err := env.stats.ByPair(ctx, env.user, nil)
	if err != nil || len(pairs) != 4 {
		t.Fatalf("by-pair: %v (%+v)", err, pairs)
	}

	setups, err := env.stats.BySetup(ctx, env.user, nil)
	if err != nil || len(setups) != 1 {
		t.Fatalf("by-setup: %v (%+v)", err, setups)
	}
	if setups[0].SetupName != "No setup" || setups[0].Trades != 5 {
		t.Fatalf("no-setup group: %+v", setups[0])
	}
}
