package service

import (
	"math/big"
	"strconv"
	"time"

	"github.com/britinogn/ctemzjournal/pkg/calc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// This file holds pgtype <-> Go converters shared by the service layer,
// plus the adapter that turns a stored trade row into calc inputs.
// The server always computes risk/pnl/R — the frontend preview is display only.

// --- pgtype readers ---

func textVal(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

func boolVal(b pgtype.Bool) *bool {
	if !b.Valid {
		return nil
	}
	v := b.Bool
	return &v
}

func timeVal(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	v := t.Time
	return &v
}

// numericFloat converts pgtype.Numeric to float64 (0 when NULL).
func numericFloat(n pgtype.Numeric) float64 {
	if !n.Valid || n.NaN || n.Int == nil {
		return 0
	}
	rat := new(big.Rat).SetInt(n.Int)
	if n.Exp != 0 {
		pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil)
		if n.Exp > 0 {
			rat.Mul(rat, new(big.Rat).SetInt(pow))
		} else {
			rat.Quo(rat, new(big.Rat).SetInt(pow))
		}
	}
	f, _ := rat.Float64()
	return f
}

// --- pgtype writers ---

func textPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func boolPtr(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

func timePtr(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func uuidPtr(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

// numericPtr formats with the given precision (prices need up to 6dp).
func numericPtr(f *float64, prec int) pgtype.Numeric {
	var n pgtype.Numeric
	if f == nil {
		return n
	}
	_ = n.Scan(strconv.FormatFloat(*f, 'f', prec, 64))
	return n
}

// --- calc adapter ---

// priceInput builds calc.TradeInput from plain values.
func priceInput(pair, direction string, entry, stopLoss, exit *float64, lots, commission, swap float64) calc.TradeInput {
	in := calc.TradeInput{
		Pair: pair, Direction: direction,
		LotSize: lots, Commission: commission, Swap: swap,
	}
	if entry != nil {
		in.Entry = *entry
	}
	if stopLoss != nil {
		in.StopLoss = *stopLoss
	}
	if exit != nil {
		in.ExitPrice = *exit
		in.HasExit = true
	}
	return in
}
