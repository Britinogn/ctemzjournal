package integration

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"os"
	"testing"
	"time"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/calc"
	"github.com/britinogn/ctemzjournal/test/fixtures"
	"github.com/britinogn/ctemzjournal/test/helpers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// fixtureTrade mirrors test/fixtures/five_trades.json.
type fixtureTrade struct {
	Name          string  `json:"name"`
	Pair          string  `json:"pair"`
	Direction     string  `json:"direction"`
	Entry         float64 `json:"entry"`
	StopLoss      float64 `json:"stop_loss"`
	ExitPrice     float64 `json:"exit_price"`
	LotSize       float64 `json:"lot_size"`
	Commission    float64 `json:"commission"`
	Swap          float64 `json:"swap"`
	FollowedRules bool    `json:"followed_rules"`
	ClosedAt      string  `json:"closed_at"`
	Expected      struct {
		Risk float64 `json:"risk"`
		Net  float64 `json:"net"`
		R    float64 `json:"r"`
	} `json:"expected"`
}

func loadFixture(t *testing.T) []fixtureTrade {
	t.Helper()
	raw, err := os.ReadFile("../fixtures/five_trades.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var trades []fixtureTrade
	if err := json.Unmarshal(raw, &trades); err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	if len(trades) != 5 {
		t.Fatalf("fixture must hold 5 trades, got %d", len(trades))
	}
	return trades
}

func numToFloat(n pgtype.Numeric) float64 {
	if !n.Valid || n.NaN || n.Int == nil {
		return 0
	}
	rat := new(big.Rat).SetInt(n.Int)
	if n.Exp != 0 {
		pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil)
		if n.Exp > 0 {
			rat.Mul(rat, new(big.Rat).SetInt(pow))
		} else {
			rat.Quo(rat, new(big.Rat).SetInt(pow))
		}
	}
	f, _ := rat.Float64()
	return f
}

func closeEnough(a, b float64) bool { return math.Abs(a-b) < 1e-6 }

type tradesEnv struct {
	svc     *service.Trades
	queries *sqlc.Queries
	userA   uuid.UUID
	userB   uuid.UUID
	account uuid.UUID
}

func newTradesEnv(t *testing.T) (*tradesEnv, context.Context) {
	t.Helper()
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	userA := fixtures.MustUser(t, database, "alice")
	userB := fixtures.MustUser(t, database, "bob")
	accounts := service.NewAccounts(repository.NewAccounts(database.Queries))
	ctx := context.Background()
	account, err := accounts.Create(ctx, userA, service.AccountCreate{Name: "Main", Type: "demo", Currency: "USD"})
	if err != nil {
		t.Fatalf("seed account: %v", err)
	}
	svc := service.NewTrades(
		repository.NewTrades(database.Queries),
		repository.NewTradeTags(database.Queries),
		repository.NewAccounts(database.Queries),
		repository.NewSetups(database.Queries),
		repository.NewTags(database.Queries),
		database.Queries,
	)
	return &tradesEnv{svc: svc, queries: database.Queries, userA: userA, userB: userB, account: account.ID}, ctx
}

func createFixtureTrade(t *testing.T, env *tradesEnv, ctx context.Context, f fixtureTrade) sqlc.Trade {
	t.Helper()
	closedAt, err := time.Parse(time.RFC3339, f.ClosedAt)
	if err != nil {
		t.Fatalf("bad fixture date: %v", err)
	}
	trade, err := env.svc.Create(ctx, env.userA, service.TradeCreate{
		AccountID: env.account, Pair: f.Pair, Direction: f.Direction,
		Entry: &f.Entry, StopLoss: &f.StopLoss, ExitPrice: &f.ExitPrice,
		LotSize: &f.LotSize, Commission: f.Commission, Swap: f.Swap,
		FollowedRules: &f.FollowedRules, Status: "closed", ClosedAt: &closedAt,
	})
	if err != nil {
		t.Fatalf("create %s: %v", f.Name, err)
	}
	if !closeEnough(numToFloat(trade.RiskAmount), f.Expected.Risk) ||
		!closeEnough(numToFloat(trade.Pnl), f.Expected.Net) ||
		!closeEnough(numToFloat(trade.RMultiple), f.Expected.R) {
		t.Fatalf("%s: stored risk=%v pnl=%v r=%v, want %v/%v/%v",
			f.Name, numToFloat(trade.RiskAmount), numToFloat(trade.Pnl), numToFloat(trade.RMultiple),
			f.Expected.Risk, f.Expected.Net, f.Expected.R)
	}
	return trade
}

// TestTradesFiveFixture is the docx §12 hand-calculation check: stored
// server-side numbers must match the fixture, and the summary must match
// win rate 60%, avg R 0.8336, expectancy 15185.6, drawdown 1312.
func TestTradesFiveFixture(t *testing.T) {
	env, ctx := newTradesEnv(t)
	for _, f := range loadFixture(t) {
		createFixtureTrade(t, env, ctx, f)
	}

	rows, err := env.queries.StatsClosedTrades(ctx, sqlc.StatsClosedTradesParams{UserID: env.userA, AccountID: pgtype.UUID{Valid: false}})
	if err != nil {
		t.Fatalf("stats rows: %v", err)
	}
	if len(rows) != 5 {
		t.Fatalf("closed rows = %d, want 5", len(rows))
	}
	closed := make([]calc.ClosedTrade, 0, len(rows))
	for _, r := range rows {
		closed = append(closed, calc.ClosedTrade{
			Net: numToFloat(r.Pnl), R: numToFloat(r.RMultiple),
			FollowedRules: r.FollowedRules.Bool,
			ClosedAt:      r.ClosedAt.Time,
		})
	}
	s := calc.Summarize(closed)
	if !closeEnough(s.WinRate, 0.6) || !closeEnough(s.AvgR, 0.8336) ||
		!closeEnough(s.Expectancy, 15185.6) || !closeEnough(s.MaxDrawdown, 1312) {
		t.Fatalf("summary mismatch: %+v", s)
	}
}

func TestTradesIsolationAndVisibility(t *testing.T) {
	env, ctx := newTradesEnv(t)
	entry, sl, lots := 1.0850, 1.0800, 1.0
	trade, err := env.svc.Create(ctx, env.userA, service.TradeCreate{
		AccountID: env.account, Pair: "EUR/USD", Direction: "long",
		Entry: &entry, StopLoss: &sl, LotSize: &lots, Status: "open",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if trade.IsPublic {
		t.Fatalf("is_public must default to false")
	}
	if _, err := env.svc.Get(ctx, trade.ID, env.userB); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("user B read user A trade: %v", err)
	}

	made, err := env.svc.SetVisibility(ctx, trade.ID, env.userA, true)
	if err != nil || !made.IsPublic {
		t.Fatalf("visibility on: %v (%+v)", err, made)
	}
	if _, err := env.svc.SetVisibility(ctx, trade.ID, env.userB, false); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("user B toggled user A trade: %v", err)
	}
	off, err := env.svc.SetVisibility(ctx, trade.ID, env.userA, false)
	if err != nil || off.IsPublic {
		t.Fatalf("visibility off: %v (%+v)", err, off)
	}
}

func TestTradesValidation(t *testing.T) {
	env, ctx := newTradesEnv(t)
	entry := 1.0850
	badAccount := uuid.New()

	cases := map[string]service.TradeCreate{
		"empty pair":      {AccountID: env.account, Direction: "long", Status: "open"},
		"bad direction":   {AccountID: env.account, Pair: "EUR/USD", Direction: "sideways", Status: "open"},
		"close no exit":   {AccountID: env.account, Pair: "EUR/USD", Direction: "long", Entry: &entry, Status: "closed"},
		"foreign account": {AccountID: badAccount, Pair: "EUR/USD", Direction: "long", Status: "open"},
	}
	for name, c := range cases {
		if _, err := env.svc.Create(ctx, env.userA, c); err == nil {
			t.Fatalf("%s: expected error", name)
		}
	}
}

func TestTradesReopen(t *testing.T) {
	env, ctx := newTradesEnv(t)
	entry, sl, exit, lots := 1.0850, 1.0800, 1.0950, 1.0
	trade, err := env.svc.Create(ctx, env.userA, service.TradeCreate{
		AccountID: env.account, Pair: "EUR/USD", Direction: "long",
		Entry: &entry, StopLoss: &sl, ExitPrice: &exit, LotSize: &lots, Status: "closed",
	})
	if err != nil {
		t.Fatalf("create closed: %v", err)
	}
	if !trade.Pnl.Valid {
		t.Fatalf("closed trade must store pnl")
	}
	status := "open"
	reopened, err := env.svc.Update(ctx, trade.ID, env.userA, service.TradeUpdate{Status: &status})
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if reopened.Status != "open" || reopened.Pnl.Valid || reopened.ExitPrice.Valid || reopened.ClosedAt.Valid {
		t.Fatalf("reopened trade keeps realized fields: %+v", reopened)
	}
}
