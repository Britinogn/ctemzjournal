package integration

import (
	"context"
	"testing"

	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/test/fixtures"
	"github.com/britinogn/ctemzjournal/test/helpers"
	"github.com/google/uuid"
)

func TestAdminOverview(t *testing.T) {
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	ctx := context.Background()

	alice := fixtures.MustUser(t, database, "alice")
	fixtures.MustUser(t, database, "bob")

	accounts := service.NewAccounts(repository.NewAccounts(database.Queries))
	account, err := accounts.Create(ctx, alice, service.AccountCreate{Name: "Main", Type: "demo"})
	if err != nil {
		t.Fatalf("seed account: %v", err)
	}
	trades := service.NewTrades(
		repository.NewTrades(database.Queries),
		repository.NewTradeTags(database.Queries),
		repository.NewAccounts(database.Queries),
		repository.NewSetups(database.Queries),
		repository.NewTags(database.Queries),
		database.Queries,
	)
	entry, sl, lots := 1.0850, 1.0800, 1.0
	for i := 0; i < 3; i++ {
		if _, err := trades.Create(ctx, alice, service.TradeCreate{
			AccountID: account.ID, Pair: "EUR/USD", Direction: "long",
			Entry: &entry, StopLoss: &sl, LotSize: &lots, Status: "open",
		}); err != nil {
			t.Fatalf("seed trade %d: %v", i, err)
		}
	}

	admin := service.NewAdmin(
		repository.NewProfiles(database.Queries),
		repository.NewTrades(database.Queries),
		repository.NewSiteSettings(database.Queries),
		repository.NewAuditLog(database.Queries),
		database.Queries,
	)
	overview, err := admin.Overview(ctx)
	if err != nil {
		t.Fatalf("overview: %v", err)
	}
	if overview.UserCount != 2 {
		t.Fatalf("users = %d, want 2", overview.UserCount)
	}
	if overview.TradeCount != 3 {
		t.Fatalf("trades = %d, want 3", overview.TradeCount)
	}
	if overview.NewUsersThisWeek != 2 {
		t.Fatalf("new users = %d, want 2", overview.NewUsersThisWeek)
	}
	if overview.ActiveUsers != 2 || overview.SuspendedUsers != 0 {
		t.Fatalf("active/suspended = %d/%d, want 2/0", overview.ActiveUsers, overview.SuspendedUsers)
	}
	if overview.PublicTrades != 0 || overview.HiddenJournals != 0 {
		t.Fatalf("public/hidden = %d/%d, want 0/0", overview.PublicTrades, overview.HiddenJournals)
	}
	if overview.NewUsersPrevWeek != 0 {
		t.Fatalf("prev week = %d, want 0", overview.NewUsersPrevWeek)
	}
	var signed int64
	for _, d := range overview.SignupsLast7D {
		signed += d.Count
	}
	if signed != 2 {
		t.Fatalf("signup chart total = %d, want 2", signed)
	}
}

func newAdminEnv(t *testing.T) (*service.Admin, *service.Trades, uuid.UUID, uuid.UUID, uuid.UUID, context.Context) {
	t.Helper()
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	ctx := context.Background()
	adminID := fixtures.MustUser(t, database, "admin")
	userID := fixtures.MustUser(t, database, "alice")
	accounts := service.NewAccounts(repository.NewAccounts(database.Queries))
	account, err := accounts.Create(ctx, userID, service.AccountCreate{Name: "Main", Type: "demo"})
	if err != nil {
		t.Fatalf("seed account: %v", err)
	}
	trades := service.NewTrades(
		repository.NewTrades(database.Queries),
		repository.NewTradeTags(database.Queries),
		repository.NewAccounts(database.Queries),
		repository.NewSetups(database.Queries),
		repository.NewTags(database.Queries),
		database.Queries,
	)
	admin := service.NewAdmin(
		repository.NewProfiles(database.Queries),
		repository.NewTrades(database.Queries),
		repository.NewSiteSettings(database.Queries),
		repository.NewAuditLog(database.Queries),
		database.Queries,
	)
	return admin, trades, adminID, userID, account.ID, ctx
}

func TestAdminUserStatusAndAudit(t *testing.T) {
	admin, _, adminID, userID, _, ctx := newAdminEnv(t)

	users, err := admin.ListUsers(ctx, "", 20, 0)
	if err != nil || len(users) != 2 {
		t.Fatalf("users: %+v (%v)", users, err)
	}
	found, err := admin.ListUsers(ctx, "ali", 20, 0)
	if err != nil || len(found) != 1 {
		t.Fatalf("search: %+v (%v)", found, err)
	}
	if found[0].TradeCount != 0 {
		t.Fatalf("alice trade_count = %d, want 0", found[0].TradeCount)
	}

	if _, err := admin.SetUserStatus(ctx, adminID, adminID, "suspended"); err == nil {
		t.Fatalf("expected self-suspend block")
	}
	if _, err := admin.SetUserStatus(ctx, adminID, userID, "banned"); err == nil {
		t.Fatalf("expected status validation error")
	}
	suspended, err := admin.SetUserStatus(ctx, adminID, userID, "suspended")
	if err != nil || suspended.Status != "suspended" {
		t.Fatalf("suspend: %+v (%v)", suspended, err)
	}
	if _, err := admin.SetUserStatus(ctx, adminID, userID, "active"); err != nil {
		t.Fatalf("reactivate: %v", err)
	}

	rows, err := admin.ListAudit(ctx, 20, 0)
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	statusWrites := 0
	for _, row := range rows {
		if row.Action == "user.status" {
			statusWrites++
		}
	}
	if statusWrites != 2 {
		t.Fatalf("audit user.status rows = %d, want 2 (suspend + reactivate)", statusWrites)
	}
}

func TestAdminJournalHide(t *testing.T) {
	admin, trades, adminID, userID, accountID, ctx := newAdminEnv(t)

	entry, sl, exit, lots := 1.0850, 1.0800, 1.0950, 1.0
	trade, err := trades.Create(ctx, userID, service.TradeCreate{
		AccountID: accountID, Pair: "EUR/USD", Direction: "long",
		Entry: &entry, StopLoss: &sl, ExitPrice: &exit, LotSize: &lots, Status: "closed",
	})
	if err != nil {
		t.Fatalf("seed trade: %v", err)
	}
	if _, err := trades.SetVisibility(ctx, trade.ID, userID, true); err != nil {
		t.Fatalf("visibility: %v", err)
	}

	journals, err := admin.ListJournals(ctx, 20, 0)
	if err != nil || len(journals) != 1 {
		t.Fatalf("journals: %+v (%v)", journals, err)
	}
	hidden, err := admin.HideJournal(ctx, adminID, trade.ID, true)
	if err != nil || !hidden.HiddenByAdmin {
		t.Fatalf("hide: %+v (%v)", hidden, err)
	}
	journals, err = admin.ListJournals(ctx, 20, 0)
	if err != nil || len(journals) != 0 {
		t.Fatalf("hidden journal still listed: %+v (%v)", journals, err)
	}
	restored, err := admin.HideJournal(ctx, adminID, trade.ID, false)
	if err != nil || restored.HiddenByAdmin {
		t.Fatalf("restore: %+v (%v)", restored, err)
	}

	rows, err := admin.ListAudit(ctx, 20, 0)
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	hides := 0
	for _, row := range rows {
		if row.Action == "journal.hide" {
			hides++
		}
	}
	if hides != 2 {
		t.Fatalf("audit journal.hide rows = %d, want 2", hides)
	}
}

func TestAdminSettingsAndLogo(t *testing.T) {
	admin, _, adminID, _, _, ctx := newAdminEnv(t)

	settings, err := admin.GetSettings(ctx)
	if err != nil || settings.ID != 1 {
		t.Fatalf("settings: %+v (%v)", settings, err)
	}
	name, allow := "My Journal", false
	updated, err := admin.UpdateSettings(ctx, adminID, service.SettingsUpdate{
		SiteName: &name, AllowSignups: &allow,
	})
	if err != nil || updated.SiteName != "My Journal" || updated.AllowSignups {
		t.Fatalf("update: %+v (%v)", updated, err)
	}

	target, err := admin.LogoTarget(ctx, "logo", "png")
	if err != nil || target.Bucket != "site-assets" || target.Path != "site-settings/logo.png" {
		t.Fatalf("logo target: %+v (%v)", target, err)
	}
	if _, err := admin.LogoTarget(ctx, "logo", "exe"); err == nil {
		t.Fatalf("expected extension validation error")
	}
	if _, err := admin.LogoTarget(ctx, "banner", "png"); err == nil {
		t.Fatalf("expected kind validation error")
	}

	rows, err := admin.ListAudit(ctx, 20, 0)
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	found := false
	for _, row := range rows {
		if row.Action == "settings.update" {
			found = true
		}
	}
	if !found {
		t.Fatalf("settings.update missing from audit")
	}
}
