package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Frankfurter is the daily fallback (no key). It quotes against a base
// currency, so XXX/USD pairs are inverted. XAU/USD is not covered and
// stays on its last cached value.
type Frankfurter struct {
	BaseURL string // default https://api.frankfurter.dev/v1
	Client  *http.Client
}

func (f Frankfurter) Fetch(ctx context.Context, pairs []string) (map[string]float64, error) {
	base := f.BaseURL
	if base == "" {
		base = "https://api.frankfurter.dev/v1"
	}
	client := f.Client
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	u := base + "/latest?base=USD&symbols=EUR,GBP,JPY,CHF,CAD,AUD,NGN"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var body struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	out := map[string]float64{}
	inv := func(sym string) (float64, bool) {
		v, ok := body.Rates[sym]
		if !ok || v <= 0 {
			return 0, false
		}
		return 1 / v, true
	}
	direct := func(sym string) (float64, bool) {
		v, ok := body.Rates[sym]
		if !ok || v <= 0 {
			return 0, false
		}
		return v, true
	}
	for _, pair := range pairs {
		parts := strings.Split(pair, "/")
		if len(parts) != 2 {
			continue
		}
		from, to := parts[0], parts[1]
		var v float64
		var ok bool
		switch {
		case from == "USD":
			v, ok = direct(to) // USD/JPY, USD/CHF, USD/CAD, USD/NGN
		case to == "USD":
			v, ok = inv(from) // EUR/USD, GBP/USD, AUD/USD
		}
		if !ok {
			continue // e.g. XAU/USD: not covered, keep last cache
		}
		if v <= 0 {
			return nil, fmt.Errorf("frankfurter %s: bad rate", pair)
		}
		out[pair] = v
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("frankfurter: no pairs covered")
	}
	return out, nil
}
