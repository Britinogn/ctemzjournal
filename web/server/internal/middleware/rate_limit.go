package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

// LoginLimiter throttles login-adjacent routes: 10 req/min per IP.
func LoginLimiter(next http.Handler) http.Handler {
	return httprate.LimitByIP(10, time.Minute)(next)
}

// UploadLimiter throttles image signature/upload routes: 30 req/min per IP.
func UploadLimiter(next http.Handler) http.Handler {
	return httprate.LimitByIP(30, time.Minute)(next)
}

// ExportLimiter throttles CSV exports (heavy: full pagination + tag joins).
func ExportLimiter(next http.Handler) http.Handler {
	return httprate.LimitByIP(10, time.Minute)(next)
}
