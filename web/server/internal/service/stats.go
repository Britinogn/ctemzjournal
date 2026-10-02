package service

import (
	"context"
	"sort"
	"time"

	"github.com/britinogn/ctemzjournal/internal/model"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/pkg/calc"
	"github.com/google/uuid"
)

// Stats serves the dashboard aggregates. Closed-trade rows come from the
// repository; all math runs through pkg/calc (unit tested).
type Stats struct {
	stats    *repository.Stats
	profiles *repository.Profiles
	setups   *repository.Setups
	trades   *repository.Trades
}

func NewStats(stats *repository.Stats, profiles *repository.Profiles, setups *repository.Setups, trades *repository.Trades) *Stats {
	return &Stats{stats: stats, profiles: profiles, setups: setups, trades: trades}
}

// locationOf resolves the user's display timezone (default Africa/Lagos),
// falling back to UTC on unknown values. Calendar grouping uses it.
func (s *Stats) locationOf(ctx context.Context, userID uuid.UUID) *time.Location {
	name := "Africa/Lagos"
	if profile, err := s.profiles.GetByID(ctx, userID); err == nil && profile.Timezone != "" {
		name = profile.Timezone
	}
	if loc, err := time.LoadLocation(name); err == nil {
		return loc
	}
	if loc, err := time.LoadLocation("Africa/Lagos"); err == nil {
		return loc
	}
	return time.UTC
}

func (s *Stats) closed(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID) ([]calc.ClosedTrade, error) {
	rows, err := s.stats.ClosedTrades(ctx, userID, accountID)
	if err != nil {
		return nil, err
	}
	out := make([]calc.ClosedTrade, 0, len(rows))
	for _, r := range rows {
		ct := calc.ClosedTrade{
			Net:           numericFloat(r.Pnl),
			R:             numericFloat(r.RMultiple),
			FollowedRules: r.FollowedRules.Bool,
			Pair:          r.Pair,
			AccountID:     r.AccountID.String(),
		}
		if r.SetupID.Valid {
			ct.SetupID = uuid.UUID(r.SetupID.Bytes).String()
		}
		if r.ClosedAt.Valid {
			ct.ClosedAt = r.ClosedAt.Time
		} else if r.OpenedAt.Valid {
			ct.ClosedAt = r.OpenedAt.Time
		}
		out = append(out, ct)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ClosedAt.Before(out[j].ClosedAt) })
	return out, nil
}

func (s *Stats) Summary(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID) (model.StatsSummary, error) {
	trades, err := s.closed(ctx, userID, accountID)
	if err != nil {
		return model.StatsSummary{}, err
	}
	sum := calc.Summarize(trades)
	out := model.StatsSummary{
		Total: sum.Total, Wins: sum.Wins, Losses: sum.Losses,
		WinRate: sum.WinRate, AvgR: sum.AvgR, Expectancy: sum.Expectancy,
		TotalPnl: sum.TotalPnl, MaxDrawdown: sum.MaxDrawdown, RuleRate: sum.RuleRate,
	}
	// Headline totals (all trades taken, currently open). Whole-user by
	// default; account-scoped when ?account= is set.
	if accountID == nil {
		if total, err := s.trades.CountByUser(ctx, userID); err == nil {
			out.TotalTrades = total
		}
		if open, err := s.trades.CountOpenByUser(ctx, userID); err == nil {
			out.OpenTrades = open
		}
	} else {
		if open, err := s.trades.CountOpenByAccount(ctx, userID, *accountID); err == nil {
			out.OpenTrades = open
			out.TotalTrades = int64(sum.Total) + open
		} else {
			out.TotalTrades = int64(sum.Total)
		}
	}
	return out, nil
}

func (s *Stats) EquityCurve(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID) ([]model.EquityPoint, error) {
	trades, err := s.closed(ctx, userID, accountID)
	if err != nil {
		return nil, err
	}
	loc := s.locationOf(ctx, userID)
	out := make([]model.EquityPoint, 0, len(trades))
	var equity float64
	for _, t := range trades {
		equity += t.Net
		out = append(out, model.EquityPoint{
			Date:   t.ClosedAt.In(loc).Format("2006-01-02"),
			Equity: equity,
		})
	}
	return out, nil
}

func (s *Stats) Calendar(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID) ([]model.CalendarDay, error) {
	trades, err := s.closed(ctx, userID, accountID)
	if err != nil {
		return nil, err
	}
	loc := s.locationOf(ctx, userID)
	return groupCalendar(trades, loc), nil
}

// groupCalendar buckets closed trades by calendar day in loc (pure, unit tested).
// The frontend colours each day by result from wins/losses/pnl.
func groupCalendar(trades []calc.ClosedTrade, loc *time.Location) []model.CalendarDay {
	days := map[string]*model.CalendarDay{}
	order := []string{}
	for _, t := range trades {
		key := t.ClosedAt.In(loc).Format("2006-01-02")
		d, ok := days[key]
		if !ok {
			d = &model.CalendarDay{Date: key}
			days[key] = d
			order = append(order, key)
		}
		d.Trades++
		d.Pnl += t.Net
		if t.Net > 0 {
			d.Wins++
		} else {
			d.Losses++
		}
	}
	sort.Strings(order)
	out := make([]model.CalendarDay, 0, len(order))
	for _, k := range order {
		out = append(out, *days[k])
	}
	return out
}

func (s *Stats) BySetup(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID) ([]model.SetupStat, error) {
	trades, err := s.closed(ctx, userID, accountID)
	if err != nil {
		return nil, err
	}
	names := map[string]string{}
	setups, err := s.setups.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, st := range setups {
		names[st.ID.String()] = st.Name
	}
	groups := calc.GroupBy(trades, func(c calc.ClosedTrade) string { return c.SetupID })
	out := make([]model.SetupStat, 0, len(groups))
	for _, g := range groups {
		st := model.SetupStat{Trades: g.Trades, Pnl: g.Pnl, AvgR: g.AvgR}
		if g.Key != "" {
			id := g.Key
			st.SetupID = &id
			st.SetupName = names[g.Key]
		} else {
			st.SetupName = "No setup"
		}
		out = append(out, st)
	}
	return out, nil
}

func (s *Stats) ByPair(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID) ([]model.PairStat, error) {
	trades, err := s.closed(ctx, userID, accountID)
	if err != nil {
		return nil, err
	}
	groups := calc.GroupBy(trades, func(c calc.ClosedTrade) string { return c.Pair })
	out := make([]model.PairStat, 0, len(groups))
	for _, g := range groups {
		out = append(out, model.PairStat{Pair: g.Key, Trades: g.Trades, Pnl: g.Pnl, AvgR: g.AvgR})
	}
	return out, nil
}
