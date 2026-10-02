package pagination

import (
	"net/http"
	"strconv"
)

// Params carries limit/offset for list endpoints.
type Params struct {
	Limit  int32
	Offset int32
}

// FromRequest parses ?limit=&offset= (defaults 20/0, clamps limit 1..100).
func FromRequest(r *http.Request) Params {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	return Params{Limit: int32(limit), Offset: int32(offset)}
}
