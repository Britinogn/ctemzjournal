package model

// AdminOverview powers the /admin overview cards (docx §6):
// user count, trade count, new users this week.
type AdminOverview struct {
	UserCount        int64 `json:"user_count"`
	TradeCount       int64 `json:"trade_count"`
	NewUsersThisWeek int64 `json:"new_users_this_week"`
}
