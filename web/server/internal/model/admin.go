package model

// AdminOverview powers the /admin overview cards (docx §6) plus the
// active/suspended split, public/hidden counts and the signup chart.
type AdminOverview struct {
	UserCount        int64       `json:"user_count"`
	TradeCount       int64       `json:"trade_count"`
	NewUsersThisWeek int64       `json:"new_users_this_week"`
	ActiveUsers      int64       `json:"active_users"`
	SuspendedUsers   int64       `json:"suspended_users"`
	PublicTrades     int64       `json:"public_trades"`
	HiddenJournals   int64       `json:"hidden_journals"`
	NewUsersPrevWeek int64       `json:"new_users_prev_week"`
	SignupsLast7D    []SignupDay `json:"signups_last_7d"`
}

// SignupDay is one day of the signup chart (days without signups absent).
type SignupDay struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// LogoTarget is the reserved site-assets path for a logo/favicon upload.
// The browser uploads to this bucket/path itself (see service docs).
type LogoTarget struct {
	Bucket string `json:"bucket"`
	Path   string `json:"path"`
}
