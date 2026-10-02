package service

import (
	"testing"
	"time"

	"github.com/britinogn/ctemzjournal/pkg/calc"
)

// TestGroupCalendarTimezone proves calendar grouping uses the display
// timezone: two trades on the same UTC day land on different Lagos days
// when they straddle midnight WAT (UTC+1).
func TestGroupCalendarTimezone(t *testing.T) {
	lagos, err := time.LoadLocation("Africa/Lagos")
	if err != nil {
		t.Skip("tzdata unavailable")
	}
	utc := time.UTC
	trades := []calc.ClosedTrade{
		{Net: 100, ClosedAt: time.Date(2026, 9, 1, 22, 30, 0, 0, utc)}, // 23:30 WAT Sep 1
		{Net: -50, ClosedAt: time.Date(2026, 9, 1, 23, 30, 0, 0, utc)}, // 00:30 WAT Sep 2
		{Net: 200, ClosedAt: time.Date(2026, 9, 2, 0, 30, 0, 0, utc)},  // 01:30 WAT Sep 2
	}
	days := groupCalendar(trades, lagos)
	if len(days) != 2 {
		t.Fatalf("days = %+v, want 2", days)
	}
	if days[0].Date != "2026-09-01" || days[0].Trades != 1 || days[0].Pnl != 100 || days[0].Wins != 1 {
		t.Errorf("day 1: %+v", days[0])
	}
	if days[1].Date != "2026-09-02" || days[1].Trades != 2 || days[1].Pnl != 150 {
		t.Errorf("day 2: %+v", days[1])
	}
	if days[1].Wins != 1 || days[1].Losses != 1 {
		t.Errorf("day 2 wins/losses: %+v", days[1])
	}
}

func TestGroupCalendarEmpty(t *testing.T) {
	if days := groupCalendar(nil, time.UTC); len(days) != 0 {
		t.Fatalf("empty: %+v", days)
	}
}
