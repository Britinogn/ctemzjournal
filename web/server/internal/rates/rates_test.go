package rates

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func twelveServer(t *testing.T, prices map[string]string, code int, message string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sym := r.URL.Query().Get("symbol")
		if code != 0 {
			_ = json.NewEncoder(w).Encode(map[string]any{"code": code, "message": message})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"price": prices[sym]})
	}))
}

func TestTwelveDataFetch(t *testing.T) {
	srv := twelveServer(t, map[string]string{"EUR/USD": "1.08512"}, 0, "")
	defer srv.Close()
	td := TwelveData{APIKey: "k", BaseURL: srv.URL, Client: srv.Client()}
	m, err := td.Fetch(context.Background(), []string{"EUR/USD"})
	if err != nil || m["EUR/USD"] != 1.08512 {
		t.Fatalf("fetch: %v (%v)", m, err)
	}
}

func TestTwelveDataError(t *testing.T) {
	srv := twelveServer(t, nil, 429, "rate limit")
	defer srv.Close()
	td := TwelveData{APIKey: "k", BaseURL: srv.URL, Client: srv.Client()}
	if _, err := td.Fetch(context.Background(), []string{"EUR/USD"}); err == nil {
		t.Fatalf("expected rate-limit error")
	}
}

func TestTwelveDataDefaultPace(t *testing.T) {
	var td TwelveData
	if got := td.paceGap(); got != 8*time.Second {
		t.Fatalf("default pace = %v, want 8s", got)
	}
}

// TestTwelveDataPacing proves pair requests are spaced out so a refresh
// cycle cannot breach the per-minute credit cap.
func TestTwelveDataPacing(t *testing.T) {
	var mu sync.Mutex
	var arrivals []time.Time
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		arrivals = append(arrivals, time.Now())
		mu.Unlock()
		_ = json.NewEncoder(w).Encode(map[string]string{"price": "1.0"})
	}))
	defer srv.Close()
	td := TwelveData{APIKey: "k", BaseURL: srv.URL, Client: srv.Client(), pace: 120 * time.Millisecond}
	pairs := []string{"EUR/USD", "GBP/USD", "USD/JPY"}
	if _, err := td.Fetch(context.Background(), pairs); err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if len(arrivals) != len(pairs) {
		t.Fatalf("requests = %d, want %d", len(arrivals), len(pairs))
	}
	for i := 1; i < len(arrivals); i++ {
		if gap := arrivals[i].Sub(arrivals[i-1]); gap < 80*time.Millisecond {
			t.Fatalf("gap %d = %v, want >= ~120ms", i, gap)
		}
	}
}

func TestTwelveDataPacingCancel(t *testing.T) {
	srv := twelveServer(t, map[string]string{"EUR/USD": "1.0"}, 0, "")
	defer srv.Close()
	td := TwelveData{APIKey: "k", BaseURL: srv.URL, Client: srv.Client(), pace: time.Hour}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := td.Fetch(ctx, []string{"EUR/USD", "GBP/USD"}); err == nil {
		t.Fatalf("expected cancellation during pacing")
	}
}

func frankServer(t *testing.T, rates map[string]float64) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"base": "USD", "rates": rates})
	}))
}

func TestFrankfurterConvert(t *testing.T) {
	srv := frankServer(t, map[string]float64{"EUR": 0.92, "JPY": 150.0, "NGN": 1500.0})
	defer srv.Close()
	f := Frankfurter{BaseURL: srv.URL, Client: srv.Client()}
	m, err := f.Fetch(context.Background(), []string{"EUR/USD", "USD/JPY", "USD/NGN", "XAU/USD"})
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if m["EUR/USD"] != 1/0.92 || m["USD/JPY"] != 150.0 || m["USD/NGN"] != 1500.0 {
		t.Fatalf("convert: %v", m)
	}
	if _, ok := m["XAU/USD"]; ok {
		t.Fatalf("XAU must stay uncovered: %v", m)
	}
}

type stubFetcher struct {
	values map[string]float64
	err    error
}

func (s stubFetcher) Fetch(_ context.Context, _ []string) (map[string]float64, error) {
	return s.values, s.err
}

func TestServicePrimaryFresh(t *testing.T) {
	svc := NewService(
		stubFetcher{values: map[string]float64{"EUR/USD": 1.08}},
		stubFetcher{err: errors.New("down")},
		[]string{"EUR/USD"},
	)
	if err := svc.Refresh(context.Background()); err != nil {
		t.Fatalf("refresh: %v", err)
	}
	quotes, _, stale := svc.Snapshot()
	if stale || len(quotes) != 1 || quotes[0].Price != 1.08 || quotes[0].Stale {
		t.Fatalf("snapshot: %+v stale=%v", quotes, stale)
	}
}

func TestServiceFallbackStale(t *testing.T) {
	svc := NewService(
		stubFetcher{err: errors.New("primary down")},
		stubFetcher{values: map[string]float64{"EUR/USD": 1.07}},
		[]string{"EUR/USD"},
	)
	if err := svc.Refresh(context.Background()); err != nil {
		t.Fatalf("refresh: %v", err)
	}
	quotes, _, stale := svc.Snapshot()
	if !stale || len(quotes) != 1 || quotes[0].Price != 1.07 {
		t.Fatalf("snapshot: %+v stale=%v", quotes, stale)
	}
}

func TestServiceTotalFailureKeepsCache(t *testing.T) {
	svc := NewService(
		stubFetcher{values: map[string]float64{"EUR/USD": 1.08}},
		stubFetcher{err: errors.New("down")},
		[]string{"EUR/USD"},
	)
	if err := svc.Refresh(context.Background()); err != nil {
		t.Fatalf("refresh 1: %v", err)
	}
	svc.primary = stubFetcher{err: errors.New("down")}
	svc.fallback = stubFetcher{err: errors.New("down")}
	if err := svc.Refresh(context.Background()); err != nil {
		t.Fatalf("refresh 2 must keep cache: %v", err)
	}
	quotes, _, stale := svc.Snapshot()
	if !stale || len(quotes) != 1 || quotes[0].Price != 1.08 {
		t.Fatalf("snapshot: %+v stale=%v", quotes, stale)
	}
}

func TestServiceEmptyCacheFails(t *testing.T) {
	svc := NewService(
		stubFetcher{err: errors.New("down")},
		stubFetcher{err: errors.New("down")},
		[]string{"EUR/USD"},
	)
	if err := svc.Refresh(context.Background()); err == nil {
		t.Fatalf("expected error on empty cache")
	}
	quotes, _, stale := svc.Snapshot()
	if !stale || len(quotes) != 0 {
		t.Fatalf("snapshot: %+v stale=%v", quotes, stale)
	}
}
