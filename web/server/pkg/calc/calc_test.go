package calc

import (
	"math"
	"testing"
	"time"
)

func approx(a, b float64) bool { return math.Abs(a-b) < 1e-6 }

func TestPipAndContractSize(t *testing.T) {
	cases := []struct {
		pair     string
		pip      float64
		contract float64
	}{
		{"EUR/USD", 0.0001, 100_000},
		{"GBP/USD", 0.0001, 100_000},
		{"USD/JPY", 0.01, 100_000},
		{"USD/CHF", 0.0001, 100_000},
		{"XAU/USD", 0.01, 100},
		{"USD/NGN", 0.0001, 100_000},
		{"BTC/USD", 0.01, 1},
		{"ETH/USD", 0.01, 1},
		{"SOL/USD", 0.01, 1},
	}
	for _, c := range cases {
		if got := PipSize(c.pair); got != c.pip {
			t.Errorf("PipSize(%s) = %v, want %v", c.pair, got, c.pip)
		}
		if got := ContractSize(c.pair); got != c.contract {
			t.Errorf("ContractSize(%s) = %v, want %v", c.pair, got, c.contract)
		}
	}
}

func TestMetricsOf(t *testing.T) {
	cases := []struct {
		name string
		in   TradeInput
		want Metrics
	}{
		{
			name: "EURUSD long win",
			in: TradeInput{Pair: "EUR/USD", Direction: "long",
				Entry: 1.0850, StopLoss: 1.0800, ExitPrice: 1.0950, HasExit: true,
				LotSize: 1, Commission: 7, Swap: 3},
			want: Metrics{Risk: 500, Gross: 1000, Net: 990, R: 1.98},
		},
		{
			name: "USDJPY short win",
			in: TradeInput{Pair: "USD/JPY", Direction: "short",
				Entry: 150.00, StopLoss: 151.00, ExitPrice: 148.50, HasExit: true,
				LotSize: 0.5},
			want: Metrics{Risk: 50000, Gross: 75000, Net: 75000, R: 1.5},
		},
		{
			name: "GBPUSD long loss with costs",
			in: TradeInput{Pair: "GBP/USD", Direction: "long",
				Entry: 1.2700, StopLoss: 1.2650, ExitPrice: 1.2660, HasExit: true,
				LotSize: 2, Commission: 10, Swap: 2},
			want: Metrics{Risk: 1000, Gross: -800, Net: -812, R: -0.812},
		},
		{
			name: "no stop loss means no R",
			in: TradeInput{Pair: "EUR/USD", Direction: "long",
				Entry: 1.1000, StopLoss: 1.1000, ExitPrice: 1.1100, HasExit: true,
				LotSize: 1},
			want: Metrics{Risk: 0, Gross: 1000, Net: 1000, R: 0},
		},
		{
			name: "open trade has no R",
			in: TradeInput{Pair: "EUR/USD", Direction: "long",
				Entry: 1.0850, StopLoss: 1.0800, HasExit: false, LotSize: 1},
			want: Metrics{Risk: 500, Gross: 0, Net: 0, R: 0},
		},
		{
			name: "short loss when price rises",
			in: TradeInput{Pair: "EUR/USD", Direction: "short",
				Entry: 1.0850, StopLoss: 1.0900, ExitPrice: 1.0920, HasExit: true,
				LotSize: 1},
			want: Metrics{Risk: 500, Gross: -700, Net: -700, R: -1.4},
		},
		{
			name: "BTCUSD long win (lots are coins)",
			in: TradeInput{Pair: "BTC/USD", Direction: "long",
				Entry: 67000, StopLoss: 66000, ExitPrice: 69000, HasExit: true,
				LotSize: 0.5},
			want: Metrics{Risk: 500, Gross: 1000, Net: 1000, R: 2},
		},
		{
			name: "ETHUSD short loss",
			in: TradeInput{Pair: "ETH/USD", Direction: "short",
				Entry: 3500, StopLoss: 3600, ExitPrice: 3550, HasExit: true,
				LotSize: 2},
			want: Metrics{Risk: 200, Gross: -100, Net: -100, R: -0.5},
		},
	}
	for _, c := range cases {
		got := MetricsOf(c.in)
		if !approx(got.Risk, c.want.Risk) || !approx(got.Gross, c.want.Gross) ||
			!approx(got.Net, c.want.Net) || !approx(got.R, c.want.R) {
			t.Errorf("%s: got %+v, want %+v", c.name, got, c.want)
		}
	}
}

// TestSummarizeFiveTrades mirrors test/fixtures/five_trades.json:
// hand-calculated numbers the stats endpoints must match (docx §12).
func TestSummarizeFiveTrades(t *testing.T) {
	day := func(d int) time.Time { return time.Date(2026, 9, d, 12, 0, 0, 0, time.UTC) }
	trades := []ClosedTrade{
		{Net: 990, R: 1.98, FollowedRules: true, ClosedAt: day(1), Pair: "EUR/USD"},
		{Net: 75000, R: 1.5, FollowedRules: true, ClosedAt: day(2), Pair: "USD/JPY"},
		{Net: -812, R: -0.812, FollowedRules: false, ClosedAt: day(3), Pair: "GBP/USD"},
		{Net: -500, R: -1, FollowedRules: true, ClosedAt: day(4), Pair: "EUR/USD"},
		{Net: 1250, R: 2.5, FollowedRules: false, ClosedAt: day(5), Pair: "XAU/USD"},
	}
	s := Summarize(trades)

	if s.Total != 5 || s.Wins != 3 || s.Losses != 2 {
		t.Fatalf("counts: %+v", s)
	}
	if !approx(s.WinRate, 0.6) {
		t.Errorf("win rate = %v, want 0.6", s.WinRate)
	}
	if !approx(s.AvgR, 0.8336) {
		t.Errorf("avg R = %v, want 0.8336", s.AvgR)
	}
	if !approx(s.TotalPnl, 75928) || !approx(s.Expectancy, 15185.6) {
		t.Errorf("pnl = %v expectancy = %v, want 75928 / 15185.6", s.TotalPnl, s.Expectancy)
	}
	if !approx(s.MaxDrawdown, 1312) {
		t.Errorf("max drawdown = %v, want 1312", s.MaxDrawdown)
	}
	if !approx(s.RuleRate, 0.6) {
		t.Errorf("rule rate = %v, want 0.6", s.RuleRate)
	}
	wantCurve := []float64{990, 75990, 75178, 74678, 75928}
	if len(s.EquityCurve) != len(wantCurve) {
		t.Fatalf("curve len %d", len(s.EquityCurve))
	}
	for i := range wantCurve {
		if !approx(s.EquityCurve[i], wantCurve[i]) {
			t.Errorf("curve[%d] = %v, want %v", i, s.EquityCurve[i], wantCurve[i])
		}
	}
}

func TestSummarizeEmpty(t *testing.T) {
	s := Summarize(nil)
	if s.Total != 0 || len(s.EquityCurve) != 0 || s.MaxDrawdown != 0 {
		t.Fatalf("empty: %+v", s)
	}
}

func TestGroupBy(t *testing.T) {
	trades := []ClosedTrade{
		{Net: 100, R: 1, Pair: "EUR/USD"},
		{Net: -50, R: -0.5, Pair: "EUR/USD"},
		{Net: 200, R: 2, Pair: "GBP/USD"},
	}
	groups := GroupBy(trades, func(c ClosedTrade) string { return c.Pair })
	if len(groups) != 2 {
		t.Fatalf("groups: %+v", groups)
	}
	if groups[0].Key != "EUR/USD" || groups[0].Trades != 2 ||
		!approx(groups[0].Pnl, 50) || !approx(groups[0].AvgR, 0.25) {
		t.Errorf("EUR group: %+v", groups[0])
	}
	if groups[1].Key != "GBP/USD" || !approx(groups[1].Pnl, 200) {
		t.Errorf("GBP group: %+v", groups[1])
	}
}
