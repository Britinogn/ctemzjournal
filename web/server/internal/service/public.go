package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/model"
	"github.com/britinogn/ctemzjournal/internal/rates"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/google/uuid"
)

// Public serves the home-page endpoints (no auth): site settings,
// recent journals, and cached rates.
type Public struct {
	settings    *repository.SiteSettings
	queries     *sqlc.Queries
	images      *repository.TradeImages
	storage     Storage
	rates       *rates.Service
	supabaseURL string
}

func NewPublic(
	settings *repository.SiteSettings,
	queries *sqlc.Queries,
	images *repository.TradeImages,
	storage Storage,
	ratesSvc *rates.Service,
	supabaseURL string,
) *Public {
	return &Public{
		settings: settings, queries: queries, images: images,
		storage: storage, rates: ratesSvc,
		supabaseURL: strings.TrimRight(supabaseURL, "/"),
	}
}

// assetURL composes the public site-assets URL for a stored logo/favicon path.
func (s *Public) assetURL(path *string) *string {
	if path == nil || *path == "" || s.supabaseURL == "" {
		return path
	}
	p := strings.TrimLeft(*path, "/")
	u := s.supabaseURL + "/storage/v1/object/public/site-assets/" + p
	return &u
}

func (s *Public) SiteSettings(ctx context.Context) (model.PublicSettings, error) {
	row, err := s.settings.Get(ctx)
	if err != nil {
		return model.PublicSettings{}, err
	}
	var social any = map[string]any{}
	if len(row.SocialLinks) > 0 {
		_ = json.Unmarshal(row.SocialLinks, &social)
	}
	return model.PublicSettings{
		SiteName:        row.SiteName,
		Tagline:         textVal(row.Tagline),
		LogoURL:         s.assetURL(textVal(row.LogoPath)),
		FaviconURL:      s.assetURL(textVal(row.FaviconPath)),
		ContactEmail:    textVal(row.ContactEmail),
		FooterText:      textVal(row.FooterText),
		RiskDisclaimer:  textVal(row.RiskDisclaimer),
		SocialLinks:     social,
		AllowSignups:    row.AllowSignups,
		MaintenanceMode: row.MaintenanceMode,
	}, nil
}

// Journals returns recent public trades projected to the safe public shape:
// display name, pair, direction, setup, win/loss, R-multiple, first signed
// image, date. Money and lots never leave the server (notes only on detail).
func (s *Public) Journals(ctx context.Context, limit, offset int32) ([]model.PublicJournal, error) {
	if limit < 1 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := s.queries.ListPublicJournals(ctx, sqlc.ListPublicJournalsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]model.PublicJournal, 0, len(rows))
	for _, row := range rows {
		j := model.PublicJournal{
			ID:          row.ID.String(),
			DisplayName: textVal(row.DisplayName),
			Pair:        row.Pair,
			Direction:   row.Direction,
			Timeframe:   textVal(row.Timeframe),
			SetupName:   textVal(row.SetupName),
		}
		if row.Pnl.Valid {
			result := "loss"
			if numericFloat(row.Pnl) > 0 {
				result = "win"
			}
			j.Result = &result
		}
		if row.RMultiple.Valid {
			j.RMultiple = ptr(numericFloat(row.RMultiple))
		}
		if images, err := s.images.ListByTrade(ctx, row.ID); err == nil && len(images) > 0 {
			if u, err := s.storage.DeliveryURL(images[0].PublicID); err == nil {
				j.ImageURL = &u
			}
		}
		date := row.CreatedAt.Time
		if row.ClosedAt.Valid {
			date = row.ClosedAt.Time
		}
		j.Date = date.UTC().Format(time.RFC3339)
		out = append(out, j)
	}
	return out, nil
}

// Journal returns one public trade with notes + every image for the detail
// page. Private, missing or admin-hidden trades are pgx.ErrNoRows — the
// handler maps that to 404 so existence stays undisclosed.
func (s *Public) Journal(ctx context.Context, id uuid.UUID) (model.PublicJournalDetail, error) {
	row, err := s.queries.GetPublicJournal(ctx, id)
	if err != nil {
		return model.PublicJournalDetail{}, err
	}
	j := model.PublicJournalDetail{
		PublicJournal: model.PublicJournal{
			ID:          row.ID.String(),
			DisplayName: textVal(row.DisplayName),
			Pair:        row.Pair,
			Direction:   row.Direction,
			Timeframe:   textVal(row.Timeframe),
			SetupName:   textVal(row.SetupName),
		},
		Notes:  textVal(row.Notes),
		Images: []model.PublicJournalImage{},
	}
	if row.Pnl.Valid {
		result := "loss"
		if numericFloat(row.Pnl) > 0 {
			result = "win"
		}
		j.Result = &result
	}
	if row.RMultiple.Valid {
		j.RMultiple = ptr(numericFloat(row.RMultiple))
	}
	if images, err := s.images.ListByTrade(ctx, row.ID); err == nil {
		for _, img := range images {
			u, err := s.storage.DeliveryURL(img.PublicID)
			if err != nil {
				continue
			}
			kind := textVal(img.Kind)
			j.Images = append(j.Images, model.PublicJournalImage{URL: u, Kind: kind})
		}
	}
	if len(j.Images) > 0 {
		first := j.Images[0].URL
		j.ImageURL = &first
	}
	date := row.CreatedAt.Time
	if row.ClosedAt.Valid {
		date = row.ClosedAt.Time
	}
	j.Date = date.UTC().Format(time.RFC3339)
	return j, nil
}

// Rates returns the cached prices with updated_at and the stale flag.
func (s *Public) Rates() model.RatesResponse {
	resp := model.RatesResponse{Rates: []model.Rate{}, Stale: true}
	if s.rates == nil {
		return resp
	}
	quotes, updatedAt, stale := s.rates.Snapshot()
	resp.Stale = stale
	if !updatedAt.IsZero() {
		u := updatedAt.UTC().Format(time.RFC3339)
		resp.UpdatedAt = &u
	}
	for _, q := range quotes {
		resp.Rates = append(resp.Rates, model.Rate{
			Pair:      q.Pair,
			Price:     q.Price,
			UpdatedAt: q.UpdatedAt.UTC().Format(time.RFC3339),
			Stale:     q.Stale,
		})
	}
	return resp
}
