package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/pkg/calc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jung-kurt/gofpdf"
)

// PDFExporter builds the designed trade-statement PDF: branded header with
// the site logo, summary strip, full trade table and totals footer.
// It pages through Trades.List, so filters (including tag/result) match
// the trade list and the CSV export exactly.
type PDFExporter struct {
	trades      *Trades
	tags        *repository.TradeTags
	setups      *repository.Setups
	settings    *repository.SiteSettings
	supabaseURL string
	http        *http.Client
}

func NewPDFExporter(
	trades *Trades,
	tags *repository.TradeTags,
	setups *repository.Setups,
	settings *repository.SiteSettings,
	supabaseURL string,
) *PDFExporter {
	return &PDFExporter{
		trades: trades, tags: tags, setups: setups, settings: settings,
		supabaseURL: strings.TrimRight(supabaseURL, "/"),
		http:        &http.Client{Timeout: 10 * time.Second},
	}
}

type pdfTrade struct {
	row      sqlc.Trade
	setup    string
	tags     string
	openedAt string
}

// Export returns (filename, pdf bytes) for the user's filtered trades.
func (e *PDFExporter) Export(ctx context.Context, userID uuid.UUID, filter TradeFilter) (string, []byte, error) {
	const pageSize = 100
	var all []sqlc.Trade
	for offset := int32(0); ; offset += pageSize {
		f := filter
		f.Limit, f.Offset = pageSize, offset
		page, err := e.trades.List(ctx, userID, f)
		if err != nil {
			return "", nil, err
		}
		all = append(all, page...)
		if len(page) < pageSize {
			break
		}
	}

	setupNames := map[string]string{}
	if setups, err := e.setups.List(ctx, userID); err == nil {
		for _, s := range setups {
			setupNames[s.ID.String()] = s.Name
		}
	}

	rows := make([]pdfTrade, 0, len(all))
	var closed []calc.ClosedTrade
	for _, t := range all {
		pt := pdfTrade{row: t}
		if t.SetupID.Valid {
			pt.setup = setupNames[uuid.UUID(t.SetupID.Bytes).String()]
			if pt.setup == "" {
				pt.setup = "Setup"
			}
		} else {
			pt.setup = "No setup"
		}
		if tags, err := e.tags.ListByTrade(ctx, t.ID); err == nil {
			names := make([]string, 0, len(tags))
			for _, tag := range tags {
				names = append(names, tag.Name)
			}
			pt.tags = strings.Join(names, ", ")
		}
		if t.OpenedAt.Valid {
			pt.openedAt = t.OpenedAt.Time.UTC().Format("2006-01-02")
		}
		if t.Status == "closed" {
			ct := calc.ClosedTrade{Net: numericFloat(t.Pnl), R: numericFloat(t.RMultiple), FollowedRules: t.FollowedRules.Bool, Pair: t.Pair}
			if t.ClosedAt.Valid {
				ct.ClosedAt = t.ClosedAt.Time
			}
			closed = append(closed, ct)
		}
		rows = append(rows, pt)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].openedAt < rows[j].openedAt })

	siteName, tagline := "Trading Journal", ""
	if settings, err := e.settings.Get(ctx); err == nil {
		if settings.SiteName != "" {
			siteName = settings.SiteName
		}
		if settings.Tagline.Valid {
			tagline = settings.Tagline.String
		}
	}
	logo, logoKind := e.fetchLogo(ctx)
	sum := calc.Summarize(closed)

	var buf bytes.Buffer
	if err := renderStatement(&buf, statementData{
		siteName: siteName, tagline: tagline, logo: logo, logoKind: logoKind,
		generatedAt: time.Now().UTC(), total: len(rows),
		summary: sum, rows: rows,
	}); err != nil {
		return "", nil, err
	}
	return "trades-export-" + time.Now().UTC().Format("20060102-150405") + ".pdf", buf.Bytes(), nil
}

// fetchLogo downloads the current site logo (PNG/JPEG only; WebP/SVG are
// skipped because the PDF writer cannot embed them). Never fails the export.
func (e *PDFExporter) fetchLogo(ctx context.Context) (data []byte, kind string) {
	if e.supabaseURL == "" {
		return nil, ""
	}
	settings, err := e.settings.Get(ctx)
	if err != nil || !settings.LogoPath.Valid || settings.LogoPath.String == "" {
		return nil, ""
	}
	lower := strings.ToLower(settings.LogoPath.String)
	switch {
	case strings.HasSuffix(lower, ".png"):
		kind = "png"
	case strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"):
		kind = "jpg"
	default:
		return nil, ""
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		e.supabaseURL+"/storage/v1/object/public/site-assets/"+strings.TrimLeft(settings.LogoPath.String, "/"), nil)
	if err != nil {
		return nil, ""
	}
	resp, err := e.http.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		return nil, ""
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil || len(data) == 0 {
		return nil, ""
	}
	return data, kind
}

type statementData struct {
	siteName    string
	tagline     string
	logo        []byte
	logoKind    string
	generatedAt time.Time
	total       int
	summary     calc.Summary
	rows        []pdfTrade
}

func renderStatement(buf *bytes.Buffer, d statementData) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetCompression(false) // keep text searchable/selectable
	pdf.SetMargins(12, 12, 12)
	pdf.SetAutoPageBreak(true, 15)
	pdf.AliasNbPages("{pages}")
	pdf.SetFooterFunc(func() {
		pdf.SetY(-12)
		pdf.SetFont("Helvetica", "", 7)
		pdf.SetTextColor(130, 130, 130)
		pdf.CellFormat(0, 5, fmt.Sprintf("Generated %s  |  Page %d of {pages}",
			d.generatedAt.Format("2006-01-02 15:04 UTC"), pdf.PageNo()), "", 0, "C", false, 0, "")
	})

	pdf.AddPage()
	top := pdf.GetY()
	if len(d.logo) > 0 {
		pdf.RegisterImageOptionsReader("sitelogo",
			gofpdf.ImageOptions{ImageType: d.logoKind, ReadDpi: true},
			bytes.NewReader(d.logo))
		pdf.Image("sitelogo", 12, top, 18, 0, false, "", 0, "")
	}
	pdf.SetXY(34, top)
	pdf.SetFont("Helvetica", "B", 17)
	pdf.SetTextColor(15, 23, 42)
	pdf.CellFormat(0, 8, d.siteName, "", 1, "L", false, 0, "")
	pdf.SetX(34)
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(100, 116, 139)
	pdf.CellFormat(0, 5, "Trade statement  ·  "+d.generatedAt.Format("2006-01-02 15:04 UTC"), "", 1, "L", false, 0, "")
	if d.tagline != "" {
		pdf.SetX(34)
		pdf.SetFont("Helvetica", "I", 9)
		pdf.CellFormat(0, 5, d.tagline, "", 1, "L", false, 0, "")
	}

	// Summary strip.
	pdf.Ln(4)
	pdf.SetFillColor(246, 248, 251)
	pdf.SetFont("Helvetica", "", 8.5)
	pdf.SetTextColor(15, 23, 42)
	stats := []string{
		fmt.Sprintf("Trades: %d", d.total),
		fmt.Sprintf("Wins: %d   Losses: %d", d.summary.Wins, d.summary.Losses),
		fmt.Sprintf("Win rate: %.1f%%", d.summary.WinRate*100),
		fmt.Sprintf("Avg R: %+.2fR", d.summary.AvgR),
		fmt.Sprintf("Net P&L: %+.2f", d.summary.TotalPnl),
		fmt.Sprintf("Max DD: %.2f", d.summary.MaxDrawdown),
	}
	colW := 273.0 / float64(len(stats))
	y0 := pdf.GetY()
	for i, s := range stats {
		pdf.SetXY(12+colW*float64(i), y0)
		pdf.CellFormat(colW, 9, s, "", 0, "C", true, 0, "")
	}
	pdf.SetY(y0 + 11)

	cols := []struct {
		title string
		w     float64
	}{
		{"Opened", 20}, {"Pair", 22}, {"Dir", 12}, {"Setup", 38},
		{"Entry", 24}, {"SL", 24}, {"TP", 24}, {"Exit", 24},
		{"Lots", 14}, {"P&L", 24}, {"R", 16}, {"Status", 16},
	}
	drawTableHead := func() {
		pdf.SetFont("Helvetica", "B", 7.5)
		pdf.SetFillColor(14, 124, 134)
		pdf.SetTextColor(255, 255, 255)
		x0 := 12.0
		for _, c := range cols {
			pdf.SetXY(x0, pdf.GetY())
			pdf.CellFormat(c.w, 7, c.title, "", 0, "C", true, 0, "")
			x0 += c.w
		}
		pdf.Ln(7)
		pdf.SetFont("Courier", "", 7.5)
		pdf.SetTextColor(15, 23, 42)
	}
	drawTableHead()

	// Rows with manual page breaks (gofpdf has no flowing tables).
	fill := false
	for _, r := range d.rows {
		if pdf.GetY() > 175 {
			pdf.AddPage()
			drawTableHead()
		}
		if fill {
			pdf.SetFillColor(246, 248, 251)
		}
		cells := []string{
			r.openedAt, r.row.Pair, r.row.Direction, truncStr(r.setup, 22),
			fnumStr(r.row.Entry), fnumStr(r.row.StopLoss), fnumStr(r.row.TakeProfit), fnumStr(r.row.ExitPrice),
			fnumStr(r.row.LotSize), fnumStr(r.row.Pnl), fnumStr(r.row.RMultiple), r.row.Status,
		}
		x0 := 12.0
		for i, c := range cols {
			pdf.SetXY(x0, pdf.GetY())
			pdf.CellFormat(c.w, 6, cells[i], "", 0, "C", fill, 0, "")
			x0 += c.w
		}
		pdf.Ln(6)
		fill = !fill
	}
	return pdf.Output(buf)
}

// fnumStr renders a nullable numeric or "" (open trades have no pnl/R yet).
func fnumStr(n pgtype.Numeric) string {
	if !n.Valid {
		return ""
	}
	return strconv.FormatFloat(numericFloat(n), 'f', -1, 64)
}

func truncStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
