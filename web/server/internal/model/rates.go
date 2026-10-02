package model

// RatesResponse is GET /public/rates: cached prices with updated_at
// and a stale flag (docx §6).
type RatesResponse struct {
	Rates     []Rate  `json:"rates"`
	UpdatedAt *string `json:"updated_at"`
	Stale     bool    `json:"stale"`
}

// Rate is one cached quote for the home-page strip.
type Rate struct {
	Pair      string  `json:"pair"`
	Price     float64 `json:"price"`
	UpdatedAt string  `json:"updated_at"`
	Stale     bool    `json:"stale"`
}
