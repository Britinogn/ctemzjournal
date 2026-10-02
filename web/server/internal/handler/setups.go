package handler

import (
	"errors"
	"net/http"

	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Setups handles /setups CRUD. Every route is user-scoped.
type Setups struct {
	svc *service.Setups
}

func NewSetups(svc *service.Setups) *Setups { return &Setups{svc: svc} }

type setupCreateRequest struct {
	Name         string  `json:"name"`
	Rules        *string `json:"rules"`
	Invalidation *string `json:"invalidation"`
}

func (h *Setups) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	var req setupCreateRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	setup, err := h.svc.Create(r.Context(), uid, service.SetupCreate(req))
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, setup)
}

func (h *Setups) List(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	setups, err := h.svc.List(r.Context(), uid)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "list failed")
		return
	}
	response.JSON(w, http.StatusOK, setups)
}

func (h *Setups) Get(w http.ResponseWriter, r *http.Request) {
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
	setup, err := h.svc.Get(r.Context(), id, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "setup not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "get failed")
		return
	}
	response.JSON(w, http.StatusOK, setup)
}

type setupUpdateRequest struct {
	Name         *string `json:"name"`
	Rules        *string `json:"rules"`
	Invalidation *string `json:"invalidation"`
}

func (h *Setups) Update(w http.ResponseWriter, r *http.Request) {
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
	var req setupUpdateRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	setup, err := h.svc.Update(r.Context(), id, uid, service.SetupUpdate(req))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "setup not found")
			return
		}
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, setup)
}

func (h *Setups) Delete(w http.ResponseWriter, r *http.Request) {
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
