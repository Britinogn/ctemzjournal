package integration

import (
	"context"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/test/fixtures"
	"github.com/britinogn/ctemzjournal/test/helpers"
	"github.com/google/uuid"
)

func TestTradesExportCSV(t *testing.T) {
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	ctx := context.Background()
	user := fixtures.MustUser(t, database, "alice")

	accounts := service.NewAccounts(repository.NewAccounts(database.Queries))
	account, err := accounts.Create(ctx, user, service.AccountCreate{Name: "Main", Type: "demo"})
	if err != nil {
		t.Fatalf("seed account: %v", err)
	}
	tagsSvc := service.NewTags(repository.NewTags(database.Queries))
	tag, err := tagsSvc.Create(ctx, user, service.TagCreate{Name: "FOMO", Kind: "mistake"})
	if err != nil {
		t.Fatalf("seed tag: %v", err)
	}
	trades := service.NewTrades(
		repository.NewTrades(database.Queries),
		repository.NewTradeTags(database.Queries),
		repository.NewAccounts(database.Queries),
		repository.NewSetups(database.Queries),
		repository.NewTags(database.Queries),
		database.Queries,
	)
	entry, sl, exit, lots := 1.0850, 1.0800, 1.0950, 1.0
	if _, err := trades.Create(ctx, user, service.TradeCreate{
		AccountID: account.ID, Pair: "EUR/USD", Direction: "long",
		Entry: &entry, StopLoss: &sl, ExitPrice: &exit, LotSize: &lots,
		Status: "closed", TagIDs: []uuid.UUID{tag.ID},
	}); err != nil {
		t.Fatalf("seed trade: %v", err)
	}

	filename, data, err := trades.ExportCSV(ctx, user, service.TradeFilter{})
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if !strings.HasPrefix(filename, "trades-export-") || !strings.HasSuffix(filename, ".csv") {
		t.Fatalf("filename: %s", filename)
	}
	rows, err := csv.NewReader(strings.NewReader(string(data))).ReadAll()
	if err != nil {
		t.Fatalf("parse csv: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %d, want header + 1", len(rows))
	}
	if rows[0][0] != "id" || rows[0][3] != "pair" {
		t.Fatalf("header: %v", rows[0])
	}
	if rows[1][3] != "EUR/USD" {
		t.Fatalf("pair: %v", rows[1])
	}
	if rows[1][23] != "FOMO" {
		t.Fatalf("tags col: %v", rows[1])
	}
	if rows[1][16] == "" {
		t.Fatalf("pnl missing: %v", rows[1])
	}
}

func TestTradesExportPDF(t *testing.T) {
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	ctx := context.Background()
	user := fixtures.MustUser(t, database, "alice")

	accounts := service.NewAccounts(repository.NewAccounts(database.Queries))
	account, err := accounts.Create(ctx, user, service.AccountCreate{Name: "Main", Type: "demo"})
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
	entry, sl, exit, lots := 1.0850, 1.0800, 1.0950, 1.0
	if _, err := trades.Create(ctx, user, service.TradeCreate{
		AccountID: account.ID, Pair: "EUR/USD", Direction: "long",
		Entry: &entry, StopLoss: &sl, ExitPrice: &exit, LotSize: &lots,
		Status: "closed",
	}); err != nil {
		t.Fatalf("seed trade: %v", err)
	}

	pdf := service.NewPDFExporter(
		trades,
		repository.NewTradeTags(database.Queries),
		repository.NewSetups(database.Queries),
		repository.NewSiteSettings(database.Queries),
		"",
	)
	filename, data, err := pdf.Export(ctx, user, service.TradeFilter{})
	if err != nil {
		t.Fatalf("export pdf: %v", err)
	}
	if !strings.HasPrefix(filename, "trades-export-") || !strings.HasSuffix(filename, ".pdf") {
		t.Fatalf("filename: %s", filename)
	}
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Fatalf("not a PDF (%d bytes)", len(data))
	}
	if !strings.Contains(string(data), "EUR/USD") {
		t.Fatalf("trade missing from PDF body")
	}
}
