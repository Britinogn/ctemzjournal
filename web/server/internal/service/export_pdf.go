package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg" // register JPEG decoder for image.Decode
	"image/png"
	"io"
	"math"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/pkg/calc"
	"github.com/go-pdf/fpdf"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	_ "golang.org/x/image/webp" // optional: accept WebP logos (re-encoded to PNG)
)

const (
	logoBucket   = "site-assets"
	maxLogoBytes = 2 << 20
	logoTTL      = 10 * time.Minute
	logoBox      = 18.0 // max logo width/height in mm

	pageWidthUsable = 273.0 // A4 landscape minus 12mm margins
	rowHeight       = 6.0
	rowBreakY       = 188.0 // A4 landscape is 210mm tall; 15mm bottom margin
)

// PDFExporter builds the designed trade-statement PDF: branded header with
// the site logo, summary strip and full trade table.
// It pages through Trades.List, so filters (including tag/result) match
// the trade list and the CSV export exactly.
type PDFExporter struct {
	trades      *Trades
	tags        *repository.TradeTags
	setups      *repository.Setups
	settings    *repository.SiteSettings
	supabaseURL string
	http        *http.Client

	// Small in-memory logo cache so we don't hit Supabase on every export.
	logoMu      sync.Mutex
	logoPath    string
	logoData    []byte
	logoExpires time.Time
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
		http:        &http.Client{Timeout: 5 * time.Second},
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
	// Advance by what we actually received and stop only on an empty page, so
	// a lower server-side limit clamp can never silently truncate the export.
	for offset := int32(0); ; {
		f := filter
		f.Limit, f.Offset = pageSize, offset
		page, err := e.trades.List(ctx, userID, f)
		if err != nil {
			return "", nil, err
		}
		if len(page) == 0 {
			break
		}
		all = append(all, page...)
		offset += int32(len(page))
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
		// NOTE: one query per trade. Swap for a batch lookup if the
		// repository grows a ListByTrades(ctx, ids) method.
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
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].openedAt < rows[j].openedAt })

	// Settings are fetched once and shared with the logo fetch.
	siteName, tagline, logoPath := "Trading Journal", "", ""
	if settings, err := e.settings.Get(ctx); err == nil {
		if settings.SiteName != "" {
			siteName = settings.SiteName
		}
		if settings.Tagline.Valid {
			tagline = settings.Tagline.String
		}
		if settings.LogoPath.Valid {
			logoPath = settings.LogoPath.String
		}
	}

	data := statementData{
		siteName: siteName, tagline: tagline, logo: e.fetchLogo(ctx, logoPath),
		generatedAt: time.Now().UTC(), total: len(rows),
		summary: calc.Summarize(closed), rows: rows,
	}

	var buf bytes.Buffer
	err := renderStatement(&buf, data)
	if err != nil && len(data.logo) > 0 {
		// A logo gofpdf can't embed must never void the whole export.
		buf.Reset()
		data.logo = nil
		err = renderStatement(&buf, data)
	}
	if err != nil {
		return "", nil, err
	}
	return "trades-export-" + time.Now().UTC().Format("20060102-150405") + ".pdf", buf.Bytes(), nil
}

// fetchLogo downloads the site logo and returns it normalized to a plain
// 8-bit PNG (so interlaced, 16-bit, mislabeled or WebP files are all safe to
// embed), or nil. It never fails the export.
func (e *PDFExporter) fetchLogo(ctx context.Context, logoPath string) []byte {
	if e.supabaseURL == "" || logoPath == "" {
		return nil
	}

	e.logoMu.Lock()
	if e.logoPath == logoPath && time.Now().Before(e.logoExpires) {
		d := e.logoData
		e.logoMu.Unlock()
		return d
	}
	e.logoMu.Unlock()

	// Clean the path (blocks "..") and escape each segment.
	clean := strings.TrimPrefix(path.Clean("/"+logoPath), "/")
	if clean == "" {
		return nil
	}
	segs := strings.Split(clean, "/")
	for i, s := range segs {
		segs[i] = url.PathEscape(s)
	}
	u := e.supabaseURL + "/storage/v1/object/public/" + logoBucket + "/" + strings.Join(segs, "/")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil
	}
	resp, err := e.http.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	// Read limit+1 so an oversized file is rejected rather than truncated.
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxLogoBytes+1))
	if err != nil || len(raw) == 0 || len(raw) > maxLogoBytes {
		return nil
	}

	img, _, err := image.Decode(bytes.NewReader(raw)) // sniffs real format, ignores extension
	if err != nil {
		return nil
	}
	b := img.Bounds()
	dst := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy())) // forces 8-bit
	draw.Draw(dst, dst.Bounds(), img, b.Min, draw.Src)

	var out bytes.Buffer
	if err := png.Encode(&out, dst); err != nil {
		return nil
	}

	e.logoMu.Lock()
	e.logoPath, e.logoData, e.logoExpires = logoPath, out.Bytes(), time.Now().Add(logoTTL)
	e.logoMu.Unlock()
	return out.Bytes()
}

type statementData struct {
	siteName    string
	tagline     string
	logo        []byte // normalized PNG, or nil
	generatedAt time.Time
	total       int
	summary     calc.Summary
	rows        []pdfTrade
}

func renderStatement(buf *bytes.Buffer, d statementData) error {
	pdf := fpdf.New("L", "mm", "A4", "")
	// Built-in fonts are cp1252; translate every string we draw.
	tr := pdf.UnicodeTranslatorFromDescriptor("")

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

	// Logo: fit inside a logoBox x logoBox square, keeping aspect ratio.
	textX, logoBottom := 12.0, top
	if len(d.logo) > 0 {
		if cfg, err := png.DecodeConfig(bytes.NewReader(d.logo)); err == nil && cfg.Width > 0 && cfg.Height > 0 {
			w, h := logoBox, logoBox
			if ratio := float64(cfg.Width) / float64(cfg.Height); ratio >= 1 {
				h = logoBox / ratio
			} else {
				w = logoBox * ratio
			}
			pdf.RegisterImageOptionsReader("sitelogo", fpdf.ImageOptions{ImageType: "png"}, bytes.NewReader(d.logo))
			pdf.Image("sitelogo", 12, top, w, h, false, "", 0, "")
			textX = 12 + w + 4
			logoBottom = top + h
		}
	}

	pdf.SetXY(textX, top)
	pdf.SetFont("Helvetica", "B", 17)
	pdf.SetTextColor(15, 23, 42)
	pdf.CellFormat(0, 8, tr(d.siteName), "", 1, "L", false, 0, "")
	pdf.SetX(textX)
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(100, 116, 139)
	pdf.CellFormat(0, 5, tr("Trade statement  ·  "+d.generatedAt.Format("2006-01-02 15:04 UTC")), "", 1, "L", false, 0, "")
	if d.tagline != "" {
		pdf.SetX(textX)
		pdf.SetFont("Helvetica", "I", 9)
		pdf.CellFormat(0, 5, tr(d.tagline), "", 1, "L", false, 0, "")
	}
	// Start the summary below whichever is taller: the logo or the text block.
	pdf.SetY(math.Max(pdf.GetY(), logoBottom))

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
	colW := pageWidthUsable / float64(len(stats))
	y0 := pdf.GetY()
	for i, s := range stats {
		pdf.SetXY(12+colW*float64(i), y0)
		pdf.CellFormat(colW, 9, s, "", 0, "C", true, 0, "")
	}
	pdf.SetY(y0 + 11)

	// Widths sum to exactly pageWidthUsable (273mm).
	cols := []struct {
		title string
		w     float64
	}{
		{"Opened", 20}, {"Pair", 22}, {"Dir", 12}, {"Setup", 34}, {"Tags", 36},
		{"Entry", 22}, {"SL", 22}, {"TP", 22}, {"Exit", 22},
		{"Lots", 12}, {"P&L", 22}, {"R", 13}, {"Status", 14},
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

	// Rows with manual page breaks (fpdf has no flowing tables).
	fill := false
	for _, r := range d.rows {
		if pdf.GetY() > rowBreakY {
			pdf.AddPage()
			drawTableHead()
		}
		if fill {
			pdf.SetFillColor(246, 248, 251)
		}
		cells := []string{
			r.openedAt, r.row.Pair, r.row.Direction, truncStr(r.setup, 18), truncStr(r.tags, 20),
			fnumStr(r.row.Entry), fnumStr(r.row.StopLoss), fnumStr(r.row.TakeProfit), fnumStr(r.row.ExitPrice),
			fnumStr(r.row.LotSize), fnumStr(r.row.Pnl), fnumStr(r.row.RMultiple), r.row.Status,
		}
		x0 := 12.0
		for i, c := range cols {
			pdf.SetXY(x0, pdf.GetY())
			switch c.title {
			case "P&L":
				setSignColor(pdf, r.row.Pnl)
			case "R":
				setSignColor(pdf, r.row.RMultiple)
			default:
				pdf.SetTextColor(15, 23, 42)
			}
			pdf.CellFormat(c.w, rowHeight, tr(cells[i]), "", 0, "C", fill, 0, "")
			x0 += c.w
		}
		pdf.Ln(rowHeight)
		fill = !fill
	}
	return pdf.Output(buf)
}

// setSignColor paints positive values green, negative red, else default ink.
func setSignColor(pdf *fpdf.Fpdf, n pgtype.Numeric) {
	if !n.Valid {
		pdf.SetTextColor(15, 23, 42)
		return
	}
	switch v := numericFloat(n); {
	case v > 0:
		pdf.SetTextColor(22, 163, 74)
	case v < 0:
		pdf.SetTextColor(220, 38, 38)
	default:
		pdf.SetTextColor(15, 23, 42)
	}
}

// fnumStr renders a nullable numeric or "" (open trades have no pnl/R yet).
func fnumStr(n pgtype.Numeric) string {
	if !n.Valid {
		return ""
	}
	return strconv.FormatFloat(numericFloat(n), 'f', -1, 64)
}

// truncStr shortens by runes (not bytes) so multibyte characters are never split.
func truncStr(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}