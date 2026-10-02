package db

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Up applies ./migrate/*.up.sql against databaseURL.
// Reads the golang-migrate mirror (supabase/migrations holds the canonical copy).
func Up(databaseURL string) error {
	m, err := migrate.New("file://migrate", databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// Down rolls back one migration step.
func Down(databaseURL string) error {
	m, err := migrate.New("file://migrate", databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()
	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
