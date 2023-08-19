package timelines

import (
	"encoding/json"
	"net/http"
	"strconv"
	"yatter-backend-go/app/handler/httperror"
)

type (
	parsedQuery struct {
		id      uint64
		isEmpty bool
	}
)

const (
	defaultLimit = 40
	maxLimit     = 80
)

// Handler request for `GET /v1/timelines`
func (h *handler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	only_media, err := parseQueryPointer(r.URL.Query().Get("only_media"))
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	max_id, err := parseQueryPointer(r.URL.Query().Get("max_id"))
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	since_id, err := parseQueryPointer(r.URL.Query().Get("since_id"))
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	limit, err := parseLimitQuery(r.URL.Query().Get("limit"))
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if objStatuses, err := h.app.Dao.Status().PublicTimeline(ctx, only_media, max_id, since_id, limit); err != nil {
		httperror.InternalServerError(w, err)
	} else if objStatuses != nil {
		if err := json.NewEncoder(w).Encode(objStatuses); err != nil {
			httperror.InternalServerError(w, err)
		}
	} else {
		httperror.NotFound(w, "timeline")
	}
}

func parseQueryPointer(s string) (*uint64, error) {
	parsed, err := ParseQuery(s)
	if err != nil {
		return nil, err
	}

	if parsed.isEmpty {
		return nil, nil
	}
	return &parsed.id, nil
}

func parseLimitQuery(s string) (*uint64, error) {
	parsed, err := ParseQuery(s)
	if err != nil {
		return nil, err
	}

	if parsed.isEmpty {
		var limit uint64 = defaultLimit
		return &limit, nil
	} else if parsed.id > maxLimit {
		var limit uint64 = maxLimit
		return &limit, nil
	}
	return &parsed.id, nil
}

func ParseQuery(s string) (parsedQuery, error) {
	if s == "" {
		return parsedQuery{isEmpty: true}, nil
	}

	if id, err := strconv.ParseUint(s, 10, 64); err != nil {
		return parsedQuery{}, err
	} else {
		return parsedQuery{id: id}, nil
	}
}
