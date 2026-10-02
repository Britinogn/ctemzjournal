package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// Pairs tracked by the rates service (USD/NGN included when the provider
// serves it — docx open decision; missing pairs simply stay absent).
var Pairs = []string{
	"EUR/USD", "GBP/USD", "USD/JPY", "USD/CHF",
	"AUD/USD", "USD/CAD", "XAU/USD", "USD/NGN",
}

// Quote is one cached price. Browsers never call providers directly.
type Quote struct {
	Pair      string    `json:"pair"`
	Price     float64   `json:"price"`
	UpdatedAt time.Time `json:"updated_at"`
	Stale     bool      `json:"stale"`
}

// Fetcher fetches spot prices for pairs (partial maps allowed).
type Fetcher interface {
	Fetch(ctx context.Context, pairs []string) (map[string]float64, error)
}

// Service holds the in-memory cache with primary/fallback sourcing.
type Service struct {
	mu        sync.RWMutex
	quotes    map[string]Quote
	updatedAt time.Time
	stale     bool
	primary   Fetcher
	fallback  Fetcher
	pairs     []string
}

func NewService(primary, fallback Fetcher, pairs []string) *Service {
	if pairs == nil {
		pairs = Pairs
	}
	return &Service{quotes: map[string]Quote{}, primary: primary, fallback: fallback, pairs: pairs, stale: true}
}

// Refresh fetches all pairs: primary first, fallback for anything missing.
// stale is false only when the primary served every pair. On total failure
// the last values keep serving (stale: true) — never an empty outage.
func (s *Service) Refresh(ctx context.Context) error {
	now := time.Now().UTC()
	got := map[string]float64{}
	primaryOK := true
	if s.primary != nil {
		m, err := s.primary.Fetch(ctx, s.pairs)
		if err != nil {
			primaryOK = false
		} else {
			got = m
		}
	} else {
		primaryOK = false
	}

	missing := []string{}
	for _, p := range s.pairs {
		if _, ok := got[p]; !ok {
			missing = append(missing, p)
		}
	}
	fromFallback := map[string]bool{}
	if len(missing) > 0 && s.fallback != nil {
		if m, err := s.fallback.Fetch(ctx, missing); err == nil {
			for p, v := range m {
				got[p] = v
				fromFallback[p] = true
			}
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if len(got) == 0 {
		s.stale = true // total failure: keep serving last values
		if len(s.quotes) == 0 {
			return fmt.Errorf("rates: primary and fallback failed")
		}
		return nil
	}
	allPrimary := primaryOK
	for p, v := range got {
		stale := fromFallback[p]
		if stale {
			allPrimary = false
		}
		s.quotes[p] = Quote{Pair: p, Price: v, UpdatedAt: now, Stale: stale}
	}
	// Pairs served by neither source keep their last quote, flagged stale.
	for _, p := range s.pairs {
		if _, ok := got[p]; !ok {
			if q, ok := s.quotes[p]; ok {
				q.Stale = true
				s.quotes[p] = q
			}
			allPrimary = false
		}
	}
	s.updatedAt = now
	s.stale = !allPrimary
	return nil
}

// Snapshot returns a copy of the cache with its updated_at and stale flag.
func (s *Service) Snapshot() (quotes []Quote, updatedAt time.Time, stale bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	quotes = make([]Quote, 0, len(s.pairs))
	for _, p := range s.pairs {
		if q, ok := s.quotes[p]; ok {
			q.Stale = q.Stale || s.stale
			quotes = append(quotes, q)
		}
	}
	return quotes, s.updatedAt, s.stale
}

// TwelveData is the primary rates source (free plan: 800 req/day).
// One request per pair per cycle: 8 pairs x 3/hour x 24h = 576/day.
type TwelveData struct {
	APIKey  string
	BaseURL string // default https://api.twelvedata.com
	Client  *http.Client
}

func (t TwelveData) Fetch(ctx context.Context, pairs []string) (map[string]float64, error) {
	base := t.BaseURL
	if base == "" {
		base = "https://api.twelvedata.com"
	}
	client := t.Client
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	out := map[string]float64{}
	for _, pair := range pairs {
		u := base + "/price?symbol=" + url.QueryEscape(pair) + "&apikey=" + url.QueryEscape(t.APIKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return nil, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		var body struct {
			Price   string `json:"price"`
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if body.Code != 0 || body.Price == "" {
			return nil, fmt.Errorf("twelve data %s: %s", pair, body.Message)
		}
		price, err := strconv.ParseFloat(body.Price, 64)
		if err != nil || price <= 0 {
			return nil, fmt.Errorf("twelve data %s: bad price %q", pair, body.Price)
		}
		out[pair] = price
	}
	return out, nil
}
