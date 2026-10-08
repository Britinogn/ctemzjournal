package handler

import (
	"net/http"
	"time"

	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/response"
	"github.com/google/uuid"
)

// Export serves the trade downloads (same filters as the trade list):
// GET /trades/export.csv and GET /trades/export.pdf.
type Export struct {
	trades *service.Trades
	pdf    *service.PDFExporter
}

func NewExport(trades *service.Trades, pdf *service.PDFExporter) *Export {
	return &Export{trades: trades, pdf: pdf}
}

// exportFilter parses the shared list filters (no pagination: full export).
func exportFilter(r *http.Request) (service.TradeFilter, bool) {
	q := r.URL.Query()
	filter := service.TradeFilter{
		Pair:      q.Get("pair"),
		Status:    q.Get("status"),
		Direction: q.Get("direction"),
		Result:    q.Get("result"),
	}
	parseUUID := func(key string, dst **uuid.UUID) bool {
		if s := q.Get(key); s != "" {
			id, err := uuid.Parse(s)
			if err != nil {
				return false
			}
			*dst = &id
		}
		return true
	}
	if !parseUUID("setup", &filter.SetupID) {
		return filter, false
	}
	if !parseUUID("account", &filter.AccountID) {
		return filter, false
	}
	if !parseUUID("tag", &filter.TagID) {
		return filter, false
	}
	parseTime := func(key string, dst **time.Time) bool {
		if s := q.Get(key); s != "" {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return false
			}
			*dst = &t
		}
		return true
	}
	if !parseTime("from", &filter.From) {
		return filter, false
	}
	if !parseTime("to", &filter.To) {
		return filter, false
	}
	return filter, true
}

func (h *Export) CSV(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	filter, ok := exportFilter(r)
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid filter")
		return
	}
	filename, data, err := h.trades.ExportCSV(r.Context(), uid, filter)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// PDF serves the designed statement (branded header, summary, table).
func (h *Export) PDF(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	filter, ok := exportFilter(r)
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid filter")
		return
	}
	filename, data, err := h.pdf.Export(r.Context(), uid, filter)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
