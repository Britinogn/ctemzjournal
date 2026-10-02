package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/britinogn/ctemzjournal/internal/middleware"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/response"
)

// Accounts handles /accounts CRUD. Every route is user-scoped.
type Accounts struct {
	svc *service.Accounts
}

func NewAccounts(svc *service.Accounts) *Accounts { return &Accounts{svc: svc} }

func userIDOf(r *http.Request) (uuid.UUID, bool) {
	u := middleware.FromContext(r.Context())
	if u == nil {
		return uuid.Nil, false
	}
	uid, err := uuid.Parse(u.ID)
	if err != nil {
		return uuid.Nil, false
	}
	return uid, true
}

type accountCreateRequest struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Currency        string  `json:"currency"`
	StartingBalance float64 `json:"starting_balance"`
}

func (h *Accounts) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	var req accountCreateRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	account, err := h.svc.Create(r.Context(), uid, service.AccountCreate(req))
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, account)
}

func (h *Accounts) List(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	accounts, err := h.svc.List(r.Context(), uid)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "list failed")
		return
	}
	response.JSON(w, http.StatusOK, accounts)
}

func (h *Accounts) Get(w http.ResponseWriter, r *http.Request) {
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
	account, err := h.svc.Get(r.Context(), id, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "account not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "get failed")
		return
	}
	response.JSON(w, http.StatusOK, account)
}

type accountUpdateRequest struct {
	Name            *string  `json:"name"`
	Type            *string  `json:"type"`
	Currency        *string  `json:"currency"`
	StartingBalance *float64 `json:"starting_balance"`
}

func (h *Accounts) Update(w http.ResponseWriter, r *http.Request) {
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
	var req accountUpdateRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	account, err := h.svc.Update(r.Context(), id, uid, service.AccountUpdate(req))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "account not found")
			return
		}
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, account)
}

func (h *Accounts) Delete(w http.ResponseWriter, r *http.Request) {
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
	if err := h.svc.Delete(r.Context(), id, uid); err != nil {
		response.Error(w, http.StatusInternalServerError, "delete failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
