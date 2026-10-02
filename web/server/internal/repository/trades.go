package repository

import (
	"context"
	"strconv"
	"time"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Trades wraps the sqlc trade queries. Every method is scoped by user_id,
// and money math is computed by the service layer (pkg/calc) before insert.
type Trades struct {
	q *sqlc.Queries
}

func NewTrades(q *sqlc.Queries) *Trades { return &Trades{q: q} }

// TradeCreate carries Go-native values; nil numerics become NULL.
type TradeCreate struct {
	AccountID     uuid.UUID
	SetupID       *uuid.UUID
	Pair          string
	Direction     string
	Timeframe     *string
	OpenedAt      *time.Time
	ClosedAt      *time.Time
	Entry         *float64
	StopLoss      *float64
	TakeProfit    *float64
	ExitPrice     *float64
	LotSize       *float64
	Commission    float64
	Swap          float64
	RiskAmount    *float64
	Pnl           *float64
	RMultiple     *float64
	FollowedRules *bool
	Emotion       *string
	Notes         *string
	Status        string
}

func (t *Trades) Create(ctx context.Context, userID uuid.UUID, in TradeCreate) (sqlc.Trade, error) {
	commission := numPtr(&in.Commission, 2)
	swap := numPtr(&in.Swap, 2)
	return t.q.CreateTrade(ctx, sqlc.CreateTradeParams{
		UserID:        userID,
		AccountID:     in.AccountID,
		SetupID:       uuidOrNull(in.SetupID),
		Pair:          in.Pair,
		Direction:     in.Direction,
		Timeframe:     textArg(in.Timeframe),
		OpenedAt:      timeOrNull(in.OpenedAt),
		ClosedAt:      timeOrNull(in.ClosedAt),
		Entry:         numPtr(in.Entry, 6),
		StopLoss:      numPtr(in.StopLoss, 6),
		TakeProfit:    numPtr(in.TakeProfit, 6),
		ExitPrice:     numPtr(in.ExitPrice, 6),
		LotSize:       numPtr(in.LotSize, 4),
		Commission:    commission,
		Swap:          swap,
		RiskAmount:    numPtr(in.RiskAmount, 2),
		Pnl:           numPtr(in.Pnl, 2),
		RMultiple:     numPtr(in.RMultiple, 4),
		FollowedRules: boolOrNull(in.FollowedRules),
		Emotion:       textArg(in.Emotion),
		Notes:         textArg(in.Notes),
		Status:        in.Status,
		IsPublic:      false, // docx: is_public defaults to false, toggled explicitly
	})
}

func (t *Trades) Get(ctx context.Context, id, userID uuid.UUID) (sqlc.Trade, error) {
	return t.q.GetTrade(ctx, sqlc.GetTradeParams{ID: id, UserID: userID})
}

// TradeFilter mirrors the list filters (pair, setup, account, status,
// direction, date range). Tag and result filters apply in the service layer.
type TradeFilter struct {
	Pair      string
	SetupID   *uuid.UUID
	Status    string
	AccountID *uuid.UUID
	Direction string
	From      *time.Time
	To        *time.Time
	Limit     int32
	Offset    int32
}

func (t *Trades) List(ctx context.Context, userID uuid.UUID, f TradeFilter) ([]sqlc.Trade, error) {
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}
	return t.q.ListTrades(ctx, sqlc.ListTradesParams{
		UserID:     userID,
		Pair:       f.Pair,
		SetupID:    uuidOrNull(f.SetupID),
		Status:     f.Status,
		AccountID:  uuidOrNull(f.AccountID),
		Direction:  f.Direction,
		OpenedFrom: timeOrNull(f.From),
		OpenedTo:   timeOrNull(f.To),
		PageOffset: f.Offset,
		PageLimit:  f.Limit,
	})
}

// TradeUpdate mirrors TradeCreate but every field is nullable for PATCH.
type TradeUpdate struct {
	AccountID     *uuid.UUID
	SetupID       *uuid.UUID
	SetupClear    bool
	Pair          *string
	Direction     *string
	Timeframe     *string
	OpenedAt      *time.Time
	ClosedAt      *time.Time
	Entry         *float64
	StopLoss      *float64
	TakeProfit    *float64
	ExitPrice     *float64
	LotSize       *float64
	Commission    *float64
	Swap          *float64
	RiskAmount    *float64
	Pnl           *float64
	RMultiple     *float64
	RiskClear     bool
	PnlClear      bool
	FollowedRules *bool
	Emotion       *string
	Notes         *string
	Status        *string
}

func (t *Trades) Update(ctx context.Context, id, userID uuid.UUID, in TradeUpdate) (sqlc.Trade, error) {
	setup := uuidOrNull(in.SetupID)
	if in.SetupClear {
		setup = pgtype.UUID{Valid: false}
	}
	risk := numPtr(in.RiskAmount, 2)
	pnl := numPtr(in.Pnl, 2)
	rMult := numPtr(in.RMultiple, 4)
	if in.RiskClear {
		risk = pgtype.Numeric{Valid: false}
	}
	if in.PnlClear {
		pnl = pgtype.Numeric{Valid: false}
		rMult = pgtype.Numeric{Valid: false}
	}
	return t.q.UpdateTrade(ctx, sqlc.UpdateTradeParams{
		AccountID:     uuidOrNull(in.AccountID),
		SetupID:       setup,
		Pair:          textArg(in.Pair),
		Direction:     textArg(in.Direction),
		Timeframe:     textArg(in.Timeframe),
		OpenedAt:      timeOrNull(in.OpenedAt),
		ClosedAt:      timeOrNull(in.ClosedAt),
		Entry:         numPtr(in.Entry, 6),
		StopLoss:      numPtr(in.StopLoss, 6),
		TakeProfit:    numPtr(in.TakeProfit, 6),
		ExitPrice:     numPtr(in.ExitPrice, 6),
		LotSize:       numPtr(in.LotSize, 4),
		Commission:    numPtr(in.Commission, 2),
		Swap:          numPtr(in.Swap, 2),
		RiskAmount:    risk,
		Pnl:           pnl,
		RMultiple:     rMult,
		FollowedRules: boolOrNull(in.FollowedRules),
		Emotion:       textArg(in.Emotion),
		Notes:         textArg(in.Notes),
		Status:        textArg(in.Status),
		ID:            id,
		UserID:        userID,
	})
}

func (t *Trades) SetVisibility(ctx context.Context, id, userID uuid.UUID, isPublic bool) (sqlc.Trade, error) {
	return t.q.UpdateTradeVisibility(ctx, sqlc.UpdateTradeVisibilityParams{
		ID: id, UserID: userID, IsPublic: isPublic,
	})
}

func (t *Trades) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return t.q.DeleteTrade(ctx, sqlc.DeleteTradeParams{ID: id, UserID: userID})
}

func (t *Trades) ExportAll(ctx context.Context, userID uuid.UUID) ([]sqlc.Trade, error) {
	return t.q.ExportTradesByUser(ctx, userID)
}

// CountByUser returns all trades taken by a user (dashboard total).
func (t *Trades) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	return t.q.CountTrades(ctx, userID)
}

// CountOpenByUser returns the user's currently open trades.
func (t *Trades) CountOpenByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	return t.q.CountOpenTrades(ctx, userID)
}

// CountOpenByAccount returns open trades for one account.
func (t *Trades) CountOpenByAccount(ctx context.Context, userID, accountID uuid.UUID) (int64, error) {
	return t.q.CountOpenTradesByAccount(ctx, sqlc.CountOpenTradesByAccountParams{UserID: userID, AccountID: accountID})
}

// CountAll returns every trade in the system (admin overview).
func (t *Trades) CountAll(ctx context.Context) (int64, error) {
	return t.q.CountAllTrades(ctx)
}

// Hide marks a trade hidden (or restores it). Admin-only; the public
// endpoint never returns hidden journals.
func (t *Trades) Hide(ctx context.Context, id uuid.UUID, hidden bool) (sqlc.Trade, error) {
	return t.q.HideJournal(ctx, sqlc.HideJournalParams{ID: id, HiddenByAdmin: hidden})
}

// --- small converters (textArg/numericArg live in accounts.go, same package) ---

// numPtr formats a float with precision (prices 6dp, money 2dp, lots 4dp).
func numPtr(f *float64, prec int) pgtype.Numeric {
	var n pgtype.Numeric
	if f == nil {
		return n
	}
	_ = n.Scan(strconv.FormatFloat(*f, 'f', prec, 64))
	return n
}

func uuidOrNull(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

func timeOrNull(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func boolOrNull(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}
