package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/pagination"
	"github.com/britinogn/ctemzjournal/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Trades handles /trades CRUD, visibility toggle and filtered listing.
type Trades struct {
	svc    *service.Trades
	images *service.Images
}

func NewTrades(svc *service.Trades, images *service.Images) *Trades {
	return &Trades{svc: svc, images: images}
}

type tradeCreateRequest struct {
	AccountID     string   `json:"account_id"`
	SetupID       *string  `json:"setup_id"`
	Pair          string   `json:"pair"`
	Direction     string   `json:"direction"`
	Timeframe     *string  `json:"timeframe"`
	OpenedAt      *string  `json:"opened_at"`
	ClosedAt      *string  `json:"closed_at"`
	Entry         *float64 `json:"entry"`
	StopLoss      *float64 `json:"stop_loss"`
	TakeProfit    *float64 `json:"take_profit"`
	ExitPrice     *float64 `json:"exit_price"`
	LotSize       *float64 `json:"lot_size"`
	Commission    float64  `json:"commission"`
	Swap          float64  `json:"swap"`
	FollowedRules *bool    `json:"followed_rules"`
	Emotion       *string  `json:"emotion"`
	Notes         *string  `json:"notes"`
	Status        string   `json:"status"`
	TagIDs        []string `json:"tag_ids"`
}

func (r tradeCreateRequest) toService() (service.TradeCreate, error) {
	accountID, err := uuid.Parse(r.AccountID)
	if err != nil {
		return service.TradeCreate{}, errors.New("invalid account_id")
	}
	out := service.TradeCreate{
		AccountID: accountID, Pair: r.Pair, Direction: r.Direction,
		Timeframe: r.Timeframe, Entry: r.Entry, StopLoss: r.StopLoss,
		TakeProfit: r.TakeProfit, ExitPrice: r.ExitPrice, LotSize: r.LotSize,
		Commission: r.Commission, Swap: r.Swap, FollowedRules: r.FollowedRules,
		Emotion: r.Emotion, Notes: r.Notes, Status: r.Status,
	}
	if r.SetupID != nil {
		id, err := uuid.Parse(*r.SetupID)
		if err != nil {
			return service.TradeCreate{}, errors.New("invalid setup_id")
		}
		out.SetupID = &id
	}
	if r.OpenedAt != nil {
		t, err := time.Parse(time.RFC3339, *r.OpenedAt)
		if err != nil {
			return service.TradeCreate{}, errors.New("invalid opened_at (RFC3339)")
		}
		out.OpenedAt = &t
	}
	if r.ClosedAt != nil {
		t, err := time.Parse(time.RFC3339, *r.ClosedAt)
		if err != nil {
			return service.TradeCreate{}, errors.New("invalid closed_at (RFC3339)")
		}
		out.ClosedAt = &t
	}
	for _, s := range r.TagIDs {
		id, err := uuid.Parse(s)
		if err != nil {
			return service.TradeCreate{}, errors.New("invalid tag_ids")
		}
		out.TagIDs = append(out.TagIDs, id)
	}
	return out, nil
}

func (h *Trades) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	var req tradeCreateRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	in, err := req.toService()
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	trade, err := h.svc.Create(r.Context(), uid, in)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, trade)
}

func (h *Trades) Get(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	trade, err := h.svc.Get(r.Context(), id, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "trade not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "get failed")
		return
	}
	response.JSON(w, http.StatusOK, trade)
}

// List supports ?pair=&setup=&account=&status=&direction=&result=win|loss
// &tag=&from=&to=&limit=&offset=. Tag/result filter in the service layer
// (documented V1 tradeoff: applied after the SQL page fetch).
func (h *Trades) List(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	q := r.URL.Query()
	f := service.TradeFilter{
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
		f.SetupID = &id
	}
	if s := q.Get("account"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid account")
			return
		}
		f.AccountID = &id
	}
	if s := q.Get("tag"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid tag")
			return
		}
		f.TagID = &id
	}
	if s := q.Get("from"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid from (RFC3339)")
			return
		}
		f.From = &t
	}
	if s := q.Get("to"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid to (RFC3339)")
			return
		}
		f.To = &t
	}
	page := pagination.FromRequest(r)
	f.Limit, f.Offset = page.Limit, page.Offset

	trades, err := h.svc.List(r.Context(), uid, f)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, trades)
}

type tradeUpdateRequest struct {
	AccountID     *string   `json:"account_id"`
	SetupID       *string   `json:"setup_id"`
	SetupClear    bool      `json:"setup_clear"`
	Pair          *string   `json:"pair"`
	Direction     *string   `json:"direction"`
	Timeframe     *string   `json:"timeframe"`
	OpenedAt      *string   `json:"opened_at"`
	ClosedAt      *string   `json:"closed_at"`
	Entry         *float64  `json:"entry"`
	StopLoss      *float64  `json:"stop_loss"`
	TakeProfit    *float64  `json:"take_profit"`
	ExitPrice     *float64  `json:"exit_price"`
	LotSize       *float64  `json:"lot_size"`
	Commission    *float64  `json:"commission"`
	Swap          *float64  `json:"swap"`
	FollowedRules *bool     `json:"followed_rules"`
	Emotion       *string   `json:"emotion"`
	Notes         *string   `json:"notes"`
	Status        *string   `json:"status"`
	TagIDs        *[]string `json:"tag_ids"`
}

func (h *Trades) Update(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req tradeUpdateRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	in := service.TradeUpdate{
		SetupClear: req.SetupClear, Pair: req.Pair, Direction: req.Direction,
		Timeframe: req.Timeframe, Entry: req.Entry, StopLoss: req.StopLoss,
		TakeProfit: req.TakeProfit, ExitPrice: req.ExitPrice, LotSize: req.LotSize,
		Commission: req.Commission, Swap: req.Swap, FollowedRules: req.FollowedRules,
		Emotion: req.Emotion, Notes: req.Notes, Status: req.Status,
	}
	if req.AccountID != nil {
		aid, err := uuid.Parse(*req.AccountID)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid account_id")
			return
		}
		in.AccountID = &aid
	}
	if req.SetupID != nil {
		sid, err := uuid.Parse(*req.SetupID)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid setup_id")
			return
		}
		in.SetupID = &sid
	}
	if req.OpenedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.OpenedAt)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid opened_at (RFC3339)")
			return
		}
		in.OpenedAt = &t
	}
	if req.ClosedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ClosedAt)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid closed_at (RFC3339)")
			return
		}
		in.ClosedAt = &t
	}
	if req.TagIDs != nil {
		ids := make([]uuid.UUID, 0, len(*req.TagIDs))
		for _, s := range *req.TagIDs {
			tid, err := uuid.Parse(s)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "invalid tag_ids")
				return
			}
			ids = append(ids, tid)
		}
		in.TagIDs = &ids
	}
	trade, err := h.svc.Update(r.Context(), id, uid, in)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "trade not found")
			return
		}
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, trade)
}

type visibilityRequest struct {
	IsPublic bool `json:"is_public"`
}

// SetVisibility toggles "Make public" (PATCH /trades/{id}/visibility).
func (h *Trades) SetVisibility(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req visibilityRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	trade, err := h.svc.SetVisibility(r.Context(), id, uid, req.IsPublic)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "trade not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "visibility failed")
		return
	}
	response.JSON(w, http.StatusOK, trade)
}

func (h *Trades) Delete(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	// Destroy Cloudinary files first (best-effort inside), then the row.
	// FK cascade removes image rows; doing files first avoids orphan files
	// in Cloudinary. Retrying delete is safe if the row delete fails.
	if h.images != nil {
		_ = h.images.DeleteByTrade(r.Context(), uid, id)
	}
	if err := h.svc.Delete(r.Context(), id, uid); err != nil {
		response.Error(w, http.StatusInternalServerError, "delete failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
