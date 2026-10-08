package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/test/fixtures"
	"github.com/britinogn/ctemzjournal/test/helpers"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// stubKF stands in for the Supabase JWKS: HS256 with a known secret.
type stubKF struct {
	secret []byte
}

func (s stubKF) Keyfunc(_ *jwt.Token) (any, error) { return s.secret, nil }
func (s stubKF) KeyfuncCtx(_ context.Context) jwt.Keyfunc {
	return s.Keyfunc
}
func (s stubKF) Storage() jwkset.Storage { return nil }
func (s stubKF) VerificationKeySet(_ context.Context) (jwt.VerificationKeySet, error) {
	return jwt.VerificationKeySet{}, nil
}

var _ keyfunc.Keyfunc = stubKF{}

func mintToken(t *testing.T, secret []byte, sub uuid.UUID) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub.String(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("mint token: %v", err)
	}
	return signed
}

func TestMaintenanceGate(t *testing.T) {
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	ctx := context.Background()

	secret := []byte("test-secret-for-maintenance")
	adminID := fixtures.MustUser(t, database, "admin")
	userID := fixtures.MustUser(t, database, "alice")
	if _, err := database.Pool.Exec(ctx, `update profiles set role = 'admin' where id = $1`, adminID); err != nil {
		t.Fatalf("promote admin: %v", err)
	}

	settings := repository.NewSiteSettings(database.Queries)
	setting := func(on bool) {
		t.Helper()
		if _, err := settings.Update(ctx, repository.SiteSettingsUpdate{MaintenanceMode: &on}); err != nil {
			t.Fatalf("toggle maintenance: %v", err)
		}
	}

	gate := middleware.NewMaintenance(database.Queries, stubKF{secret: secret})
	ok := func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }

	call := func(path, token string) int {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		rec := httptest.NewRecorder()
		gate.RequireAvailable(http.HandlerFunc(ok)).ServeHTTP(rec, req)
		return rec.Code
	}

	setting(false)
	if code := call("/me", ""); code != http.StatusOK {
		t.Fatalf("flag off, anonymous = %d, want 200", code)
	}

	setting(true)
	if code := call("/me", ""); code != http.StatusServiceUnavailable {
		t.Fatalf("flag on, anonymous = %d, want 503", code)
	}
	if code := call("/me", mintToken(t, secret, userID)); code != http.StatusServiceUnavailable {
		t.Fatalf("flag on, user = %d, want 503", code)
	}
	if code := call("/me", mintToken(t, secret, adminID)); code != http.StatusOK {
		t.Fatalf("flag on, admin = %d, want 200", code)
	}
	if code := call("/public/site-settings", ""); code != http.StatusOK {
		t.Fatalf("settings exempt = %d, want 200", code)
	}
	if code := call("/healthz", ""); code != http.StatusOK {
		t.Fatalf("healthz exempt = %d, want 200", code)
	}

	setting(false)
}
