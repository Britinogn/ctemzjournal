package repository

import (
	"context"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
)

// SiteSettings wraps the single-row site_settings table (id = 1).
type SiteSettings struct {
	q *sqlc.Queries
}

func NewSiteSettings(q *sqlc.Queries) *SiteSettings { return &SiteSettings{q: q} }

func (s *SiteSettings) Get(ctx context.Context) (sqlc.SiteSetting, error) {
	return s.q.GetSiteSettings(ctx)
}

type SiteSettingsUpdate struct {
	SiteName        *string
	Tagline         *string
	LogoPath        *string
	FaviconPath     *string
	ContactEmail    *string
	FooterText      *string
	RiskDisclaimer  *string
	SocialLinks     []byte
	AllowSignups    *bool
	MaintenanceMode *bool
}

func (s *SiteSettings) Update(ctx context.Context, in SiteSettingsUpdate) (sqlc.SiteSetting, error) {
	return s.q.UpdateSiteSettings(ctx, sqlc.UpdateSiteSettingsParams{
		SiteName:        textArg(in.SiteName),
		Tagline:         textArg(in.Tagline),
		LogoPath:        textArg(in.LogoPath),
		FaviconPath:     textArg(in.FaviconPath),
		ContactEmail:    textArg(in.ContactEmail),
		FooterText:      textArg(in.FooterText),
		RiskDisclaimer:  textArg(in.RiskDisclaimer),
		SocialLinks:     in.SocialLinks, // nil keeps the current value
		AllowSignups:    boolOrNull(in.AllowSignups),
		MaintenanceMode: boolOrNull(in.MaintenanceMode),
	})
}
