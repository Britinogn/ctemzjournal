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
// Money amounts, lot size and notes are NEVER included.
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
