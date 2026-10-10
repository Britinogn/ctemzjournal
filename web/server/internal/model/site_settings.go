package model

import (
	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
)

// SiteSettings mirrors the single-row site_settings table.
type SiteSettings = sqlc.SiteSetting

// PublicSettings is the home-page view of site settings.
type PublicSettings struct {
	SiteName        string  `json:"site_name"`
	Tagline         *string `json:"tagline"`
	LogoURL         *string `json:"logo_url"`
	FaviconURL      *string `json:"favicon_url"`
	ContactEmail    *string `json:"contact_email"`
	FooterText      *string `json:"footer_text"`
	RiskDisclaimer  *string `json:"risk_disclaimer"`
	SocialLinks     any     `json:"social_links"`
	AllowSignups    bool    `json:"allow_signups"`
	MaintenanceMode bool    `json:"maintenance_mode"`
}

// PublicJournal is one recent public trade for the home page.
// Money amounts and lot size are NEVER included.
type PublicJournal struct {
	ID          string   `json:"id"`
	DisplayName *string  `json:"display_name"`
	Pair        string   `json:"pair"`
	Direction   string   `json:"direction"`
	Timeframe   *string  `json:"timeframe"`
	SetupName   *string  `json:"setup_name"`
	Result      *string  `json:"result"` // "win" | "loss" | null (open)
	RMultiple   *float64 `json:"r_multiple"`
	ImageURL    *string  `json:"image_url"` // first image, signed URL
	Date        string   `json:"date"`
}

// PublicJournalImage is one signed trade image on the journal detail page.
type PublicJournalImage struct {
	URL  string  `json:"url"`
	Kind *string `json:"kind,omitempty"`
}

// PublicJournalDetail is the GET /public/journals/{id} shape: the public
// fields plus the trader's notes and every image. Money and lots stay out.
// Notes here are public to anyone on the internet — the trader opted in by
// making the trade public.
type PublicJournalDetail struct {
	PublicJournal
	Notes  *string              `json:"notes"`
	Images []PublicJournalImage `json:"images"`
}
