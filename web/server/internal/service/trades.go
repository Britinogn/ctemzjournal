package service

import (
	"context"
	"errors"
	"time"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/model"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/pkg/calc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Trades owns trade rules: account/setup ownership, server-side risk/pnl/R,
// visibility defaults, and list filters.
type Trades struct {
	trades   *repository.Trades
	tags     *repository.TradeTags
	accounts *repository.Accounts
	setups   *repository.Setups
	tagRepo  *repository.Tags
	queries  *sqlc.Queries
}

func NewTrades(
	trades *repository.Trades,
	tags *repository.TradeTags,
	accounts *repository.Accounts,
	setups *repository.Setups,
	tagRepo *repository.Tags,
	queries *sqlc.Queries,
) *Trades {
	return &Trades{trades: trades, tags: tags, accounts: accounts, setups: setups, tagRepo: tagRepo, queries: queries}
}

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
	FollowedRules *bool
	Emotion       *string
	Notes         *string
	Status        string
	TagIDs        []uuid.UUID
}

func (s *Trades) Create(ctx context.Context, userID uuid.UUID, in TradeCreate) (sqlc.Trade, error) {
	if in.Pair == "" {
		return sqlc.Trade{}, errors.New("pair is required")
	}
	if in.Direction != model.DirectionLong && in.Direction != model.DirectionShort {
		return sqlc.Trade{}, errors.New("direction must be long or short")
	}
	if in.Status == "" {
		in.Status = model.TradeOpen
	}
	if in.Status != model.TradeOpen && in.Status != model.TradeClosed {
		return sqlc.Trade{}, errors.New("status must be open or closed")
	}
	if _, err := s.accounts.Get(ctx, in.AccountID, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.Trade{}, errors.New("account not found")
		}
		return sqlc.Trade{}, err
	}
	if in.SetupID != nil {
		if _, err := s.setups.Get(ctx, *in.SetupID, userID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return sqlc.Trade{}, errors.New("setup not found")
			}
			return sqlc.Trade{}, err
		}
	}
	if err := s.checkTags(ctx, userID, in.TagIDs); err != nil {
		return sqlc.Trade{}, err
	}
	if in.Status == model.TradeClosed && in.ExitPrice == nil {
		return sqlc.Trade{}, errors.New("exit_price is required to close a trade")
	}

	risk, pnl, rMult := deriveMetrics(in.Pair, in.Direction, in.Entry, in.StopLoss, in.ExitPrice, in.LotSize, in.Status, in.Commission, in.Swap)

	trade, err := s.trades.Create(ctx, userID, repository.TradeCreate{
		AccountID: in.AccountID, SetupID: in.SetupID, Pair: in.Pair, Direction: in.Direction,
		Timeframe: in.Timeframe, OpenedAt: in.OpenedAt, ClosedAt: closedAtOr(in.ClosedAt, in.Status),
		Entry: in.Entry, StopLoss: in.StopLoss, TakeProfit: in.TakeProfit, ExitPrice: in.ExitPrice,
		LotSize: in.LotSize, Commission: in.Commission, Swap: in.Swap,
		RiskAmount: risk, Pnl: pnl, RMultiple: rMult,
		FollowedRules: in.FollowedRules, Emotion: in.Emotion, Notes: in.Notes, Status: in.Status,
	})
	if err != nil {
		return sqlc.Trade{}, err
	}
	if len(in.TagIDs) > 0 {
		if err := s.tags.SetTags(ctx, trade.ID, in.TagIDs); err != nil {
			return sqlc.Trade{}, err
		}
	}
	return trade, nil
}

func (s *Trades) Get(ctx context.Context, id, userID uuid.UUID) (sqlc.Trade, error) {
	return s.trades.Get(ctx, id, userID)
}

// TradeFilter extends the repository filter with tag + result (win/loss).
type TradeFilter struct {
	Pair      string
	SetupID   *uuid.UUID
	AccountID *uuid.UUID
	Status    string
	Direction string
	Result    string // "win" | "loss" | ""
	TagID     *uuid.UUID
	From      *time.Time
	To        *time.Time
	Limit     int32
	Offset    int32
}

func (s *Trades) List(ctx context.Context, userID uuid.UUID, f TradeFilter) ([]sqlc.Trade, error) {
	if f.Result != "" && f.Result != "win" && f.Result != "loss" {
		return nil, errors.New("result must be win or loss")
	}
	allowed := map[uuid.UUID]bool{}
	if f.TagID != nil {
		if _, err := s.tagRepo.Get(ctx, *f.TagID, userID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []sqlc.Trade{}, nil
			}
			return nil, err
		}
		ids, err := s.queries.ListTradeIDsByTag(ctx, *f.TagID)
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			allowed[id] = true
		}
	}
	rows, err := s.trades.List(ctx, userID, repository.TradeFilter{
		Pair: f.Pair, SetupID: f.SetupID, Status: f.Status, AccountID: f.AccountID,
		Direction: f.Direction, From: f.From, To: f.To, Limit: f.Limit, Offset: f.Offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]sqlc.Trade, 0, len(rows))
	for _, t := range rows {
		if f.TagID != nil && !allowed[t.ID] {
			continue
		}
		if f.Result != "" {
			if !t.Pnl.Valid {
				continue
			}
			pnl := numericFloat(t.Pnl)
			if f.Result == "win" && pnl <= 0 {
				continue
			}
			if f.Result == "loss" && pnl >= 0 {
				continue
			}
		}
		out = append(out, t)
	}
	return out, nil
}

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
	FollowedRules *bool
	Emotion       *string
	Notes         *string
	Status        *string
	TagIDs        *[]uuid.UUID // nil = unchanged, else replace
}

func (s *Trades) Update(ctx context.Context, id, userID uuid.UUID, in TradeUpdate) (sqlc.Trade, error) {
	current, err := s.trades.Get(ctx, id, userID)
	if err != nil {
		return sqlc.Trade{}, err
	}
	if in.AccountID != nil {
		if _, err := s.accounts.Get(ctx, *in.AccountID, userID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return sqlc.Trade{}, errors.New("account not found")
			}
			return sqlc.Trade{}, err
		}
	}
	if in.SetupID != nil {
		if _, err := s.setups.Get(ctx, *in.SetupID, userID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return sqlc.Trade{}, errors.New("setup not found")
			}
			return sqlc.Trade{}, err
		}
	}
	if in.TagIDs != nil {
		if err := s.checkTags(ctx, userID, *in.TagIDs); err != nil {
			return sqlc.Trade{}, err
		}
	}

	// Merge current row with the patch to recompute metrics.
	merged := TradeCreate{
		Pair:      current.Pair,
		Direction: current.Direction,
		Status:    current.Status,
	}
	if v := textVal(current.Timeframe); v != nil {
		merged.Timeframe = v
	}
	if v := timeVal(current.OpenedAt); v != nil {
		merged.OpenedAt = v
	}
	if v := timeVal(current.ClosedAt); v != nil {
		merged.ClosedAt = v
	}
	if current.Entry.Valid {
		merged.Entry = ptr(numericFloat(current.Entry))
	}
	if current.StopLoss.Valid {
		merged.StopLoss = ptr(numericFloat(current.StopLoss))
	}
	if current.TakeProfit.Valid {
		merged.TakeProfit = ptr(numericFloat(current.TakeProfit))
	}
	if current.ExitPrice.Valid {
		merged.ExitPrice = ptr(numericFloat(current.ExitPrice))
	}
	if current.LotSize.Valid {
		merged.LotSize = ptr(numericFloat(current.LotSize))
	}
	merged.Commission = numericFloat(current.Commission)
	merged.Swap = numericFloat(current.Swap)
	if v := boolVal(current.FollowedRules); v != nil {
		merged.FollowedRules = v
	}
	if v := textVal(current.Emotion); v != nil {
		merged.Emotion = v
	}
	if v := textVal(current.Notes); v != nil {
		merged.Notes = v
	}
	if current.SetupID.Valid {
		sid := uuid.UUID(current.SetupID.Bytes)
		merged.SetupID = &sid
	}

	// Apply the patch.
	if in.Pair != nil {
		merged.Pair = *in.Pair
	}
	if in.Direction != nil {
		if *in.Direction != model.DirectionLong && *in.Direction != model.DirectionShort {
			return sqlc.Trade{}, errors.New("direction must be long or short")
		}
		merged.Direction = *in.Direction
	}
	if in.Status != nil {
		if *in.Status != model.TradeOpen && *in.Status != model.TradeClosed {
			return sqlc.Trade{}, errors.New("status must be open or closed")
		}
		merged.Status = *in.Status
	}
	if in.Timeframe != nil {
		merged.Timeframe = in.Timeframe
	}
	if in.OpenedAt != nil {
		merged.OpenedAt = in.OpenedAt
	}
	if in.ClosedAt != nil {
		merged.ClosedAt = in.ClosedAt
	}
	if in.Entry != nil {
		merged.Entry = in.Entry
	}
	if in.StopLoss != nil {
		merged.StopLoss = in.StopLoss
	}
	if in.TakeProfit != nil {
		merged.TakeProfit = in.TakeProfit
	}
	if in.ExitPrice != nil {
		merged.ExitPrice = in.ExitPrice
	}
	if in.LotSize != nil {
		merged.LotSize = in.LotSize
	}
	if in.Commission != nil {
		merged.Commission = *in.Commission
	}
	if in.Swap != nil {
		merged.Swap = *in.Swap
	}
	if in.FollowedRules != nil {
		merged.FollowedRules = in.FollowedRules
	}
	if in.Emotion != nil {
		merged.Emotion = in.Emotion
	}
	if in.Notes != nil {
		merged.Notes = in.Notes
	}
	if merged.Pair == "" {
		return sqlc.Trade{}, errors.New("pair is required")
	}
	if merged.Status == model.TradeClosed && merged.ExitPrice == nil {
		return sqlc.Trade{}, errors.New("exit_price is required to close a trade")
	}

	risk, pnl, rMult := deriveMetrics(merged.Pair, merged.Direction, merged.Entry, merged.StopLoss, merged.ExitPrice, merged.LotSize, merged.Status, merged.Commission, merged.Swap)

	updated, err := s.trades.Update(ctx, id, userID, repository.TradeUpdate{
		AccountID: in.AccountID, SetupID: in.SetupID, SetupClear: in.SetupClear,
		Pair: in.Pair, Direction: in.Direction, Timeframe: in.Timeframe,
		OpenedAt: in.OpenedAt, ClosedAt: closedAtOr(in.ClosedAt, merged.Status),
		Entry: in.Entry, StopLoss: in.StopLoss, TakeProfit: in.TakeProfit, ExitPrice: in.ExitPrice,
		LotSize: in.LotSize, Commission: in.Commission, Swap: in.Swap,
		RiskAmount: risk, Pnl: pnl, RMultiple: rMult,
		RiskClear: risk == nil, PnlClear: pnl == nil,
		FollowedRules: in.FollowedRules, Emotion: in.Emotion, Notes: in.Notes, Status: in.Status,
	})
	if err != nil {
		return sqlc.Trade{}, err
	}
	if merged.Status == model.TradeOpen && current.Status == model.TradeClosed {
		// Reopened: drop exit/closed_at/pnl/R so the open trade shows no
		// realized numbers (patch fields were applied above).
		updated, err = s.queries.ReopenTrade(ctx, sqlc.ReopenTradeParams{ID: id, UserID: userID})
		if err != nil {
			return sqlc.Trade{}, err
		}
	}
	if in.TagIDs != nil {
		if err := s.tags.SetTags(ctx, id, *in.TagIDs); err != nil {
			return sqlc.Trade{}, err
		}
	}
	return updated, nil
}

// SetVisibility toggles is_public (default false; explicit user action).
func (s *Trades) SetVisibility(ctx context.Context, id, userID uuid.UUID, isPublic bool) (sqlc.Trade, error) {
	return s.trades.SetVisibility(ctx, id, userID, isPublic)
}

func (s *Trades) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.trades.Delete(ctx, id, userID)
}

func (s *Trades) checkTags(ctx context.Context, userID uuid.UUID, tagIDs []uuid.UUID) error {
	for _, tagID := range tagIDs {
		if _, err := s.tagRepo.Get(ctx, tagID, userID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.New("tag not found")
			}
			return err
		}
	}
	return nil
}

// deriveMetrics computes server-side risk/pnl/R. Risk needs entry+stop+lots;
// pnl/R additionally need a closed trade with an exit. Nil means "unknown"
// (open trade) and is stored as NULL.
func deriveMetrics(pair, direction string, entry, stopLoss, exit, lots *float64, status string, commission, swap float64) (risk, pnl, rMult *float64) {
	lotsVal := 0.0
	if lots != nil {
		lotsVal = *lots
	}
	if entry != nil && stopLoss != nil && lots != nil && *lots > 0 {
		m := calc.MetricsOf(priceInput(pair, direction, entry, stopLoss, nil, lotsVal, 0, 0))
		risk = ptr(m.Risk)
	}
	if status == model.TradeClosed && exit != nil {
		m := calc.MetricsOf(priceInput(pair, direction, entry, stopLoss, exit, lotsVal, commission, swap))
		pnl = ptr(m.Net)
		rMult = ptr(m.R)
	}
	return risk, pnl, rMult
}

// closedAtOr defaults closed_at to now when a trade is closed without one.
func closedAtOr(t *time.Time, status string) *time.Time {
	if t != nil {
		return t
	}
	if status == model.TradeClosed {
		now := time.Now().UTC()
		return &now
	}
	return nil
}

func ptr[T any](v T) *T { return &v }
