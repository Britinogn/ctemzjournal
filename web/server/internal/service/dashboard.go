package service

import (
	"context"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/model"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/google/uuid"
)

// RecentTradesLimit bounds the recent-trades section of the overview.
const RecentTradesLimit = 10

// Dashboard composes one overview of what's happening in a user's account:
// profile, accounts, headline totals, equity, calendar, breakdowns and
// recent trades — a single GET /dashboard for the first screen.
// The granular endpoints stay for filters and drill-downs.
type Dashboard struct {
	profiles *repository.Profiles
	accounts *Accounts
	stats    *Stats
	trades   *repository.Trades
}

func NewDashboard(
	profiles *repository.Profiles,
	accounts *Accounts,
	stats *Stats,
	trades *repository.Trades,
) *Dashboard {
	return &Dashboard{profiles: profiles, accounts: accounts, stats: stats, trades: trades}
}

// Overview is the aggregated dashboard payload. accountID scopes the
// numbers to one account (?account=); nil covers the whole user.
type Overview struct {
	Me       sqlc.Profile        `json:"me"`
	Accounts []sqlc.Account      `json:"accounts"`
	Summary  model.StatsSummary  `json:"summary"`
	Equity   []model.EquityPoint `json:"equity_curve"`
	Calendar []model.CalendarDay `json:"calendar"`
	BySetup  []model.SetupStat   `json:"by_setup"`
	ByPair   []model.PairStat    `json:"by_pair"`
	Recent   []sqlc.Trade        `json:"recent_trades"`
}

func (s *Dashboard) Overview(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID) (Overview, error) {
	var out Overview
	var err error

	if out.Me, err = s.profiles.GetByID(ctx, userID); err != nil {
		return Overview{}, err
	}
	if out.Accounts, err = s.accounts.List(ctx, userID); err != nil {
		return Overview{}, err
	}
	if out.Summary, err = s.stats.Summary(ctx, userID, accountID); err != nil {
		return Overview{}, err
	}
	if out.Equity, err = s.stats.EquityCurve(ctx, userID, accountID); err != nil {
		return Overview{}, err
	}
	if out.Calendar, err = s.stats.Calendar(ctx, userID, accountID); err != nil {
		return Overview{}, err
	}
	if out.BySetup, err = s.stats.BySetup(ctx, userID, accountID); err != nil {
		return Overview{}, err
	}
	if out.ByPair, err = s.stats.ByPair(ctx, userID, accountID); err != nil {
		return Overview{}, err
	}
	recent, err := s.trades.List(ctx, userID, repository.TradeFilter{Limit: RecentTradesLimit})
	if err != nil {
		return Overview{}, err
	}
	if accountID != nil {
		filtered := recent[:0]
		for _, t := range recent {
			if t.AccountID == *accountID {
				filtered = append(filtered, t)
			}
		}
		recent = filtered
	}
	if recent == nil {
		recent = []sqlc.Trade{}
	}
	out.Recent = recent
	return out, nil
}
