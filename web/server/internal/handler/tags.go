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

// Tags handles /tags CRUD. Every route is user-scoped.
type Tags struct {
	svc *service.Tags
}

func NewTags(svc *service.Tags) *Tags { return &Tags{svc: svc} }

type tagCreateRequest struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

func (h *Tags) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	var req tagCreateRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	tag, err := h.svc.Create(r.Context(), uid, service.TagCreate(req))
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, tag)
}

func (h *Tags) List(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDOf(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	tags, err := h.svc.List(r.Context(), uid)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "list failed")
		return
	}
	response.JSON(w, http.StatusOK, tags)
}

func (h *Tags) Get(w http.ResponseWriter, r *http.Request) {
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
	tag, err := h.svc.Get(r.Context(), id, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "tag not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "get failed")
		return
	}
	response.JSON(w, http.StatusOK, tag)
}

type tagUpdateRequest struct {
	Name *string `json:"name"`
	Kind *string `json:"kind"`
}

func (h *Tags) Update(w http.ResponseWriter, r *http.Request) {
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
	var req tagUpdateRequest
	if err := response.Decode(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	tag, err := h.svc.Update(r.Context(), id, uid, service.TagUpdate(req))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "tag not found")
			return
		}
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, tag)
}

func (h *Tags) Delete(w http.ResponseWriter, r *http.Request) {
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
