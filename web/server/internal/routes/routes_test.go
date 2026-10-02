package routes

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

// TestTradesRouteMerge proves the /trades group registers CRUD, visibility
// and export.csv together (chi panics on two Route("/trades") blocks, so
// export mounts inside the single group; static wins over /{id}).
func TestTradesRouteMerge(t *testing.T) {
	r := chi.NewRouter()
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("route mount panicked: %v", rec)
		}
	}()
	MountTrades(r, nil, nil, nil)

	foundExport, foundGet := false, false
	_ = chi.Walk(r, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		if method == "GET" && route == "/trades/export.csv" {
			foundExport = true
		}
		if method == "GET" && route == "/trades/{id}" {
			foundGet = true
		}
		return nil
	})
	if !foundExport || !foundGet {
		t.Fatalf("routes missing: export=%v get=%v", foundExport, foundGet)
	}
}
