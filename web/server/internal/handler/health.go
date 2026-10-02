package handler

import (
	"net/http"

	"github.com/britinogn/ctemzjournal/pkg/response"
)

// Health answers GET /healthz (no auth).
func Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
