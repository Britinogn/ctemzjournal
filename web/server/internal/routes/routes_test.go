package routes

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

// TestTradesRouteMerge proves the /trades group registers CRUD, visibility
// and both exports together (chi panics on two Route("/trades") blocks, so
// exports mount inside the single group; static wins over /{id}).
func TestTradesRouteMerge(t *testing.T) {
	r := chi.NewRouter()
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("route mount panicked: %v", rec)
		}
	}()
	MountTrades(r, nil, nil, nil)

	foundCSV, foundPDF, foundGet := false, false, false
	_ = chi.Walk(r, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		if method == "GET" && route == "/trades/export.csv" {
			foundCSV = true
		}
		if method == "GET" && route == "/trades/export.pdf" {
			foundPDF = true
		}
		if method == "GET" && route == "/trades/{id}" {
			foundGet = true
		}
		return nil
	})
	if !foundCSV || !foundPDF || !foundGet {
		t.Fatalf("routes missing: csv=%v pdf=%v get=%v", foundCSV, foundPDF, foundGet)
	}
}
