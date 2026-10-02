package handler

import (
	"net/http"
	"time"

	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/response"
	"github.com/google/uuid"
)

// Export serves GET /trades/export.csv (same filters as the trade list).
type Export struct {
	trades *service.Trades
}

func NewExport(trades *service.Trades) *Export { return &Export{trades: trades} }

func (h *Export) CSV(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	q := r.URL.Query()
	filter := service.TradeFilter{
		Pair:      q.Get("pair"),
		Status:    q.Get("status"),
		Direction: q.Get("direction"),
		Result:    q.Get("result"),
	}
	if s := q.Get("setup"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid setup")
			return
		}
		filter.SetupID = &id
	}
	if s := q.Get("account"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid account")
			return
		}
		filter.AccountID = &id
	}
	if s := q.Get("tag"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid tag")
			return
		}
		filter.TagID = &id
	}
	if s := q.Get("from"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid from (RFC3339)")
			return
		}
		filter.From = &t
	}
	if s := q.Get("to"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid to (RFC3339)")
			return
		}
		filter.To = &t
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
