// Package calc holds the server-side trade math (docx §7).
// The server calculates these values — the frontend preview is display only.
// Everything here is pure and covered by table-driven tests in calc_test.go,
// because wrong numbers would mislead the trader (highest-risk part).
package calc

import (
	"math"
	"sort"
	"strings"
	"time"
)

// --- per-trade math ---

// PipSize returns the price increment of one pip for a pair.
// JPY pairs quote to 2 decimals (0.01), most forex to 4 (0.0001),
// XAU/USD (gold) and crypto (BTC/ETH/SOL, quoted in dollars/cents) to 2 (0.01).
func PipSize(pair string) float64 {
	p := strings.ToUpper(pair)
	if strings.Contains(p, "JPY") {
		return 0.01
	}
	if strings.Contains(p, "XAU") {
		return 0.01
	}
	if strings.Contains(p, "BTC") || strings.Contains(p, "ETH") || strings.Contains(p, "SOL") {
		return 0.01
	}
	return 0.0001
}

// ContractSize returns units per 1.0 lot: 100_000 for forex, 100 oz for
// gold, 1 coin for crypto (so lot_size means coins there).
func ContractSize(pair string) float64 {
	u := strings.ToUpper(pair)
	if strings.Contains(u, "XAU") {
		return 100
	}
	if strings.Contains(u, "BTC") || strings.Contains(u, "ETH") || strings.Contains(u, "SOL") {
		return 1
	}
	return 100_000
}

// TradeInput carries one trade's price fields in price units.
type TradeInput struct {
	Pair       string
	Direction  string // "long" or "short"
	Entry      float64
	StopLoss   float64
	ExitPrice  float64
	HasExit    bool
	LotSize    float64
	Commission float64
	Swap       float64
}

// Metrics holds the server-calculated money values (quote-currency units;
// conversion to account currency is a documented V1 limitation).
type Metrics struct {
	Risk  float64
	Gross float64
	Net   float64
	R     float64
}

// RiskAmount = |entry - stop_loss| x lot size x contract size.
func RiskAmount(in TradeInput) float64 {
	return math.Abs(in.Entry-in.StopLoss) * in.LotSize * ContractSize(in.Pair)
}

// GrossPnl is the raw price move in quote-currency units.
// Long: (exit - entry); short: (entry - exit).
// Open trades (no exit yet) have no realized pnl: 0.
func GrossPnl(in TradeInput) float64 {
	if !in.HasExit {
		return 0
	}
	move := in.ExitPrice - in.Entry
	if strings.ToLower(in.Direction) == "short" {
		move = in.Entry - in.ExitPrice
	}
	return move * in.LotSize * ContractSize(in.Pair)
}

// NetPnl = gross pnl - commission - swap.
func NetPnl(in TradeInput) float64 {
	return GrossPnl(in) - in.Commission - in.Swap
}

// RMultiple = net pnl / risk_amount (0 when there is no risk or no exit).
func RMultiple(in TradeInput) float64 {
	if !in.HasExit {
		return 0
	}
	risk := RiskAmount(in)
	if risk <= 0 {
		return 0
	}
	return NetPnl(in) / risk
}

// MetricsOf computes risk, gross, net and R together.
func MetricsOf(in TradeInput) Metrics {
	return Metrics{
		Risk:  RiskAmount(in),
		Gross: GrossPnl(in),
		Net:   NetPnl(in),
		R:     RMultiple(in),
	}
}

// --- stats math (closed trades only) ---

// ClosedTrade is one closed trade for aggregation.
type ClosedTrade struct {
	Net           float64
	R             float64
	FollowedRules bool
	ClosedAt      time.Time
	Pair          string
	SetupID       string
	AccountID     string
}

// Summary aggregates win rate, average R, expectancy, drawdown, rule rate.
type Summary struct {
	Total       int
	Wins        int
	Losses      int
	WinRate     float64
	AvgR        float64
	Expectancy  float64 // average net pnl per trade (quote units)
	TotalPnl    float64
	MaxDrawdown float64
	RuleRate    float64
	EquityCurve []float64
}

// Summarize aggregates closed trades (sorted by ClosedAt internally).
func Summarize(trades []ClosedTrade) Summary {
	sorted := append([]ClosedTrade(nil), trades...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ClosedAt.Before(sorted[j].ClosedAt) })

	var s Summary
	s.Total = len(sorted)
	if s.Total == 0 {
		s.EquityCurve = []float64{}
		return s
	}
	curve := make([]float64, 0, len(sorted))
	var equity, rSum, netSum float64
	var followed int
	for _, t := range sorted {
		if t.Net > 0 {
			s.Wins++
		} else {
			s.Losses++
		}
		rSum += t.R
		netSum += t.Net
		equity += t.Net
		curve = append(curve, equity)
		if t.FollowedRules {
			followed++
		}
	}
	s.WinRate = float64(s.Wins) / float64(s.Total)
	s.AvgR = rSum / float64(s.Total)
	s.Expectancy = netSum / float64(s.Total)
	s.TotalPnl = netSum
	s.RuleRate = float64(followed) / float64(s.Total)
	s.MaxDrawdown = MaxDrawdown(curve)
	s.EquityCurve = curve
	return s
}

// MaxDrawdown is the largest peak-to-trough drop of an equity curve.
func MaxDrawdown(equity []float64) float64 {
	var peak, maxDD float64
	for _, e := range equity {
		if e > peak {
			peak = e
		}
		if dd := peak - e; dd > maxDD {
			maxDD = dd
		}
	}
	return maxDD
}

// ByGroup aggregates trade count, total pnl and average R per key
// (setup id or pair) for the by-setup / by-pair endpoints.
type GroupStat struct {
	Key    string
	Trades int
	Pnl    float64
	AvgR   float64
}

// GroupBy aggregates with keyFn extracting the group key per trade.
func GroupBy(trades []ClosedTrade, keyFn func(ClosedTrade) string) []GroupStat {
	groups := map[string]*GroupStat{}
	order := []string{}
	for _, t := range trades {
		k := keyFn(t)
		g, ok := groups[k]
		if !ok {
			g = &GroupStat{Key: k}
			groups[k] = g
			order = append(order, k)
		}
		g.Trades++
		g.Pnl += t.Net
		g.AvgR += t.R
	}
	out := make([]GroupStat, 0, len(groups))
	for _, k := range order {
		g := groups[k]
		g.AvgR /= float64(g.Trades)
		out = append(out, *g)
	}
	return out
}
