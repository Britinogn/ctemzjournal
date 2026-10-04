package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/model"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/google/uuid"
)

// Admin serves the /admin/* endpoints. Role gating happens in middleware
// (RequireRole admin); every admin write is recorded in audit_log.
type Admin struct {
	profiles *repository.Profiles
	trades   *repository.Trades
	settings *repository.SiteSettings
	audit    *repository.AuditLog
	queries  *sqlc.Queries
}

func NewAdmin(
	profiles *repository.Profiles,
	trades *repository.Trades,
	settings *repository.SiteSettings,
	audit *repository.AuditLog,
	queries *sqlc.Queries,
) *Admin {
	return &Admin{profiles: profiles, trades: trades, settings: settings, audit: audit, queries: queries}
}

// Overview returns the admin cards: totals, active/suspended split,
// public/hidden journal counts, new users this week + previous week,
// and per-day signups for the last 7 days.
func (s *Admin) Overview(ctx context.Context) (model.AdminOverview, error) {
	var out model.AdminOverview
	var err error
	now := time.Now().UTC()
	if out.UserCount, err = s.profiles.Count(ctx); err != nil {
		return model.AdminOverview{}, err
	}
	if out.ActiveUsers, err = s.profiles.CountByStatus(ctx, model.StatusActive); err != nil {
		return model.AdminOverview{}, err
	}
	if out.SuspendedUsers, err = s.profiles.CountByStatus(ctx, model.StatusSuspended); err != nil {
		return model.AdminOverview{}, err
	}
	if out.TradeCount, err = s.trades.CountAll(ctx); err != nil {
		return model.AdminOverview{}, err
	}
	if out.PublicTrades, err = s.trades.CountPublic(ctx); err != nil {
		return model.AdminOverview{}, err
	}
	if out.HiddenJournals, err = s.trades.CountHidden(ctx); err != nil {
		return model.AdminOverview{}, err
	}
	if out.NewUsersThisWeek, err = s.profiles.CountNewSince(ctx, now.AddDate(0, 0, -7)); err != nil {
		return model.AdminOverview{}, err
	}
	if out.NewUsersPrevWeek, err = s.profiles.CountNewSince(ctx, now.AddDate(0, 0, -14)); err != nil {
		return model.AdminOverview{}, err
	} else {
		out.NewUsersPrevWeek -= out.NewUsersThisWeek
	}
	rows, err := s.profiles.SignupsSince(ctx, now.AddDate(0, 0, -7))
	if err != nil {
		return model.AdminOverview{}, err
	}
	out.SignupsLast7D = make([]model.SignupDay, 0, len(rows))
	for _, r := range rows {
		day := ""
		if r.Day.Valid {
			day = r.Day.Time.Format("2006-01-02")
		}
		out.SignupsLast7D = append(out.SignupsLast7D, model.SignupDay{Date: day, Count: r.Signups})
	}
	return out, nil
}

// ListUsers returns profiles with email + trade counts and name-or-email search.
func (s *Admin) ListUsers(ctx context.Context, search string, limit, offset int32) ([]sqlc.ListUsersAdminRow, error) {
	return s.profiles.ListUsersAdmin(ctx, search, limit, offset)
}

// SetUserStatus suspends or reactivates a user and audits the write.
// Admins cannot suspend themselves.
func (s *Admin) SetUserStatus(ctx context.Context, adminID, targetID uuid.UUID, status string) (sqlc.Profile, error) {
	if status != model.StatusActive && status != model.StatusSuspended {
		return sqlc.Profile{}, errors.New("status must be active or suspended")
	}
	if adminID == targetID {
		return sqlc.Profile{}, errors.New("cannot change your own status")
	}
	profile, err := s.profiles.SetStatus(ctx, targetID, status)
	if err != nil {
		return sqlc.Profile{}, err
	}
	meta, _ := json.Marshal(map[string]string{"status": status})
	_ = s.audit.Record(ctx, &adminID, "user.status", "profile", targetID.String(), meta)
	return profile, nil
}

// ListJournals returns public journals for moderation, INCLUDING hidden
// ones so they can be restored. The public endpoint keeps its own filter.
func (s *Admin) ListJournals(ctx context.Context, limit, offset int32) ([]sqlc.ListAdminJournalsRow, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := s.queries.ListAdminJournals(ctx, sqlc.ListAdminJournalsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []sqlc.ListAdminJournalsRow{}
	}
	return rows, nil
}

// HideJournal hides (or restores) a public journal and audits the write.
// Hidden journals never appear on the public endpoint.
func (s *Admin) HideJournal(ctx context.Context, adminID, tradeID uuid.UUID, hidden bool) (sqlc.Trade, error) {
	trade, err := s.trades.Hide(ctx, tradeID, hidden)
	if err != nil {
		return sqlc.Trade{}, err
	}
	meta, _ := json.Marshal(map[string]bool{"hidden_by_admin": hidden})
	_ = s.audit.Record(ctx, &adminID, "journal.hide", "trade", tradeID.String(), meta)
	return trade, nil
}

// GetSettings returns the single site-settings row.
func (s *Admin) GetSettings(ctx context.Context) (sqlc.SiteSetting, error) {
	return s.settings.Get(ctx)
}

type SettingsUpdate struct {
	SiteName        *string
	Tagline         *string
	LogoPath        *string
	FaviconPath     *string
	ContactEmail    *string
	FooterText      *string
	RiskDisclaimer  *string
	SocialLinks     map[string]any
	SocialLinksSet  bool
	AllowSignups    *bool
	MaintenanceMode *bool
}

// UpdateSettings patches site settings and audits the write.
// Saving here updates the home page and header everywhere (frontend
// caches settings with a long stale time).
func (s *Admin) UpdateSettings(ctx context.Context, adminID uuid.UUID, in SettingsUpdate) (sqlc.SiteSetting, error) {
	var social []byte
	if in.SocialLinksSet {
		social, _ = json.Marshal(in.SocialLinks)
		if social == nil {
			social = []byte("{}")
		}
	}
	row, err := s.settings.Update(ctx, repository.SiteSettingsUpdate{
		SiteName: in.SiteName, Tagline: in.Tagline, LogoPath: in.LogoPath,
		FaviconPath: in.FaviconPath, ContactEmail: in.ContactEmail, FooterText: in.FooterText,
		RiskDisclaimer: in.RiskDisclaimer, SocialLinks: social,
		AllowSignups: in.AllowSignups, MaintenanceMode: in.MaintenanceMode,
	})
	if err != nil {
		return sqlc.SiteSetting{}, err
	}
	meta, _ := json.Marshal(map[string]bool{"updated": true})
	_ = s.audit.Record(ctx, &adminID, "settings.update", "site_settings", "1", meta)
	return row, nil
}

// LogoTarget reserves the canonical logo/favicon path in the public
// site-assets bucket and returns its public URL.
//
// NOTE: the Supabase service key was intentionally left out of this
// backend, so the browser uploads the file itself with Supabase JS
// (bucket site-assets, path below) and then PATCHes settings with the
// returned path. Only the path reservation happens here.
func (s *Admin) LogoTarget(_ context.Context, kind, extension string) (model.LogoTarget, error) {
	if kind != "logo" && kind != "favicon" {
		return model.LogoTarget{}, errors.New("kind must be logo or favicon")
	}
	extension = strings.ToLower(strings.TrimPrefix(extension, "."))
	allowed := map[string]bool{"png": true, "jpg": true, "jpeg": true, "webp": true, "svg": true, "ico": true}
	if !allowed[extension] {
		return model.LogoTarget{}, errors.New("extension must be png, jpg, jpeg, webp, svg or ico")
	}
	path := "site-settings/" + kind + "." + extension
	return model.LogoTarget{Bucket: "site-assets", Path: path}, nil
}

// ListAudit returns admin action history newest-first.
func (s *Admin) ListAudit(ctx context.Context, limit, offset int32) ([]sqlc.AuditLog, error) {
	return s.audit.List(ctx, limit, offset)
}
