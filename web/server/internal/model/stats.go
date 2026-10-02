package model

// Stats DTOs for the /stats/* endpoints (snake_case for the frontend).

type StatsSummary struct {
	Total       int     `json:"total"` // closed trades
	Wins        int     `json:"wins"`
	Losses      int     `json:"losses"`
	WinRate     float64 `json:"win_rate"`
	AvgR        float64 `json:"avg_r"`
	Expectancy  float64 `json:"expectancy"`
	TotalPnl    float64 `json:"total_pnl"`
	MaxDrawdown float64 `json:"max_drawdown"`
	RuleRate    float64 `json:"rule_rate"`
	OpenTrades  int64   `json:"open_trades"`
	TotalTrades int64   `json:"total_trades"` // open + closed
}

type EquityPoint struct {
	Date   string  `json:"date"`
	Equity float64 `json:"equity"`
}

type CalendarDay struct {
	Date   string  `json:"date"`
	Trades int     `json:"trades"`
	Wins   int     `json:"wins"`
	Losses int     `json:"losses"`
	Pnl    float64 `json:"pnl"`
}

type SetupStat struct {
	SetupID   *string `json:"setup_id"`
	SetupName string  `json:"setup_name"`
	Trades    int     `json:"trades"`
	Pnl       float64 `json:"pnl"`
	AvgR      float64 `json:"avg_r"`
}

type PairStat struct {
	Pair   string  `json:"pair"`
	Trades int     `json:"trades"`
	Pnl    float64 `json:"pnl"`
	AvgR   float64 `json:"avg_r"`
}
