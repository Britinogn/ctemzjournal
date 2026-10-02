package integration

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/rates"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/test/fixtures"
	"github.com/britinogn/ctemzjournal/test/helpers"
	"github.com/google/uuid"
)

type stubRates struct {
	values map[string]float64
	err    error
}

func (s stubRates) Fetch(_ context.Context, _ []string) (map[string]float64, error) {
	return s.values, s.err
}

type publicEnv struct {
	public  *service.Public
	queries *sqlc.Queries
	trades  *service.Trades
	user    uuid.UUID
	account uuid.UUID
}

func newPublicEnv(t *testing.T) (*publicEnv, context.Context) {
	t.Helper()
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	user := fixtures.MustUser(t, database, "alice")
	accountsRepo := repository.NewAccounts(database.Queries)
	accounts := service.NewAccounts(accountsRepo)
	ctx := context.Background()
	account, err := accounts.Create(ctx, user, service.AccountCreate{Name: "Main", Type: "demo"})
	if err != nil {
		t.Fatalf("seed account: %v", err)
	}
	ratesSvc := rates.NewService(
		stubRates{values: map[string]float64{"EUR/USD": 1.08}},
		stubRates{err: errors.New("down")},
		[]string{"EUR/USD"},
	)
	if err := ratesSvc.Refresh(ctx); err != nil {
		t.Fatalf("seed rates: %v", err)
	}
	trades := service.NewTrades(
		repository.NewTrades(database.Queries),
		repository.NewTradeTags(database.Queries),
		accountsRepo,
		repository.NewSetups(database.Queries),
		repository.NewTags(database.Queries),
		database.Queries,
	)
	public := service.NewPublic(
		repository.NewSiteSettings(database.Queries),
		database.Queries,
		repository.NewTradeImages(database.Queries),
		&fakeStorage{},
		ratesSvc,
		"https://example.supabase.co",
	)
	return &publicEnv{public: public, queries: database.Queries, trades: trades, user: user, account: account.ID}, ctx
}

func TestPublicJournals(t *testing.T) {
	env, ctx := newPublicEnv(t)

	entry, sl, exit, lots := 1.0850, 1.0800, 1.0950, 1.0
	trade, err := env.trades.Create(ctx, env.user, service.TradeCreate{
		AccountID: env.account, Pair: "EUR/USD", Direction: "long",
		Entry: &entry, StopLoss: &sl, ExitPrice: &exit, LotSize: &lots, Status: "closed",
	})
	if err != nil {
		t.Fatalf("seed trade: %v", err)
	}
	// Private trades never appear.
	journals, err := env.public.Journals(ctx, 20, 0)
	if err != nil || len(journals) != 0 {
		t.Fatalf("private trade leaked: %+v (%v)", journals, err)
	}

	if _, err := env.trades.SetVisibility(ctx, trade.ID, env.user, true); err != nil {
		t.Fatalf("visibility: %v", err)
	}
	journals, err = env.public.Journals(ctx, 20, 0)
	if err != nil || len(journals) != 1 {
		t.Fatalf("journals: %+v (%v)", journals, err)
	}
	j := journals[0]
	if j.Result == nil || *j.Result != "win" || j.RMultiple == nil {
		t.Fatalf("journal result: %+v", j)
	}
	if j.Date == "" || j.Pair != "EUR/USD" {
		t.Fatalf("journal fields: %+v", j)
	}
	// Money, lots and notes must never serialize.
	raw, _ := json.Marshal(j)
	for _, leaked := range []string{"pnl", "lot", "notes", "commission", "entry", "risk"} {
		if strings.Contains(strings.ToLower(string(raw)), `"`+leaked) {
			t.Fatalf("leaked field %q in %s", leaked, raw)
		}
	}
}

func TestPublicSiteSettingsAndRates(t *testing.T) {
	env, ctx := newPublicEnv(t)

	settings, err := env.public.SiteSettings(ctx)
	if err != nil {
		t.Fatalf("settings: %v", err)
	}
	if settings.SiteName == "" {
		t.Fatalf("empty site name: %+v", settings)
	}

	resp := env.public.Rates()
	if resp.Stale || len(resp.Rates) != 1 || resp.Rates[0].Price != 1.08 {
		t.Fatalf("rates: %+v", resp)
	}
	if resp.UpdatedAt == nil {
		t.Fatalf("rates missing updated_at")
	}
}
