package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"strconv"
	"strings"
	"time"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ExportCSV dumps the user's trades (honouring the list filters) as CSV.
// Paging is internal (100/page); tags are joined into one column.
func (s *Trades) ExportCSV(ctx context.Context, userID uuid.UUID, filter TradeFilter) (filename string, data []byte, err error) {
	const pageSize = 100
	var all []sqlc.Trade
	for offset := int32(0); ; offset += pageSize {
		filter.Limit, filter.Offset = pageSize, offset
		page, err := s.List(ctx, userID, filter)
		if err != nil {
			return "", nil, err
		}
		all = append(all, page...)
		if len(page) < pageSize {
			break
		}
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	header := []string{
		"id", "account_id", "setup_id", "pair", "direction", "timeframe",
		"opened_at", "closed_at", "entry", "stop_loss", "take_profit", "exit_price",
		"lot_size", "commission", "swap", "risk_amount", "pnl", "r_multiple",
		"followed_rules", "emotion", "notes", "status", "is_public", "tags", "created_at",
	}
	if err := w.Write(header); err != nil {
		return "", nil, err
	}
	for _, t := range all {
		var tagNames []string
		if tags, err := s.tags.ListByTrade(ctx, t.ID); err == nil {
			for _, tag := range tags {
				tagNames = append(tagNames, tag.Name)
			}
		}
		row := []string{
			t.ID.String(), t.AccountID.String(), uuidText(t.SetupID), safe(t.Pair), safe(t.Direction), safe(textStr(t.Timeframe)),
			timeStr(t.OpenedAt), timeStr(t.ClosedAt), numStr(t.Entry), numStr(t.StopLoss),
			numStr(t.TakeProfit), numStr(t.ExitPrice), numStr(t.LotSize), numStr(t.Commission),
			numStr(t.Swap), numStr(t.RiskAmount), numStr(t.Pnl), numStr(t.RMultiple),
			boolStr(t.FollowedRules), safe(textStr(t.Emotion)), safe(textStr(t.Notes)), t.Status,
			strconv.FormatBool(t.IsPublic), safe(strings.Join(tagNames, ";")), timeStr(t.CreatedAt),
		}
		if err := w.Write(row); err != nil {
			return "", nil, err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", nil, err
	}
	filename = "trades-export-" + time.Now().UTC().Format("20060102-150405") + ".csv"
	return filename, buf.Bytes(), nil
}

func uuidText(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	return uuid.UUID(id.Bytes).String()
}

// safe neutralizes CSV formula injection: fields starting with = + - @
// (or tab) would execute as spreadsheet formulas on open. Prefixing with
// a single quote keeps the visible value while defusing execution.
func safe(s string) string {
	if s == "" {
		return s
	}
	switch s[0] {
	case '=', '+', '-', '@', '\t':
		return "'" + s
	}
	return s
}

func textStr(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func numStr(n pgtype.Numeric) string {
	if !n.Valid {
		return ""
	}
	return strconv.FormatFloat(numericFloat(n), 'f', -1, 64)
}

func boolStr(b pgtype.Bool) string {
	if !b.Valid {
		return ""
	}
	return strconv.FormatBool(b.Bool)
}

func timeStr(t pgtype.Timestamptz) string {
	if !t.Valid {
		return ""
	}
	return t.Time.UTC().Format(time.RFC3339)
}
