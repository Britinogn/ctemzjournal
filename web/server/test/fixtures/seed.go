package fixtures

import (
	"context"
	"testing"

	"github.com/britinogn/ctemzjournal/internal/db"
	"github.com/google/uuid"
)

// MustUser creates a profiles row directly (bypasses auth.users trigger)
// and returns its id. Callers clean up via Truncate.
func MustUser(t *testing.T, database *db.DB, displayName string) uuid.UUID {
	t.Helper()
	id := uuid.New()
	_, err := database.Pool.Exec(context.Background(),
		`insert into profiles (id, display_name) values ($1, $2)`, id, displayName)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return id
}
