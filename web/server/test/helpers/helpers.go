package helpers

import (
	"context"
	"os"
	"testing"

	"github.com/britinogn/ctemzjournal/internal/db"
)

// MustConnect opens the test database or skips the test when no URL is set.
// Set TEST_DATABASE_URL (falls back to DATABASE_URL) to run integration tests:
// committees:
//
//	$env:TEST_DATABASE_URL="postgres://..." ; go test ./test/integration/ -v
func MustConnect(t *testing.T) *db.DB {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		url = os.Getenv("DATABASE_URL")
	}
	if url == "" {
		t.Skip("TEST_DATABASE_URL/DATABASE_URL unset — skipping integration test")
	}
	database, err := db.Connect(context.Background(), url)
	if err != nil {
		t.Fatalf("test db connect: %v", err)
	}
	t.Cleanup(database.Close)
	return database
}

// Truncate wipes user-owned tables between tests (respects FK order).
func Truncate(t *testing.T, database *db.DB) {
	t.Helper()
	ctx := context.Background()
	for _, q := range []string{
		`truncate trade_images, trade_tags, trades, tags, setups, accounts, audit_log cascade`,
	} {
		if _, err := database.Pool.Exec(ctx, q); err != nil {
			t.Fatalf("truncate: %v", err)
		}
	}
}
