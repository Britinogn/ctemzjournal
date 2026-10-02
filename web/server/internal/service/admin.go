package service

import (
	"context"
	"time"

	"github.com/britinogn/ctemzjournal/internal/model"
	"github.com/britinogn/ctemzjournal/internal/repository"
)

// Admin serves the /admin/* endpoints. Role gating happens in middleware
// (RequireRole admin); every admin write is recorded in audit_log.
type Admin struct {
	profiles *repository.Profiles
	trades   *repository.Trades
}

func NewAdmin(profiles *repository.Profiles, trades *repository.Trades) *Admin {
	return &Admin{profiles: profiles, trades: trades}
}

// Overview returns user count, trade count and new users this week.
func (s *Admin) Overview(ctx context.Context) (model.AdminOverview, error) {
	var out model.AdminOverview
	var err error
	if out.UserCount, err = s.profiles.Count(ctx); err != nil {
		return model.AdminOverview{}, err
	}
	if out.TradeCount, err = s.trades.CountAll(ctx); err != nil {
		return model.AdminOverview{}, err
	}
	if out.NewUsersThisWeek, err = s.profiles.CountNewSince(ctx, time.Now().UTC().AddDate(0, 0, -7)); err != nil {
		return model.AdminOverview{}, err
	}
	return out, nil
}
