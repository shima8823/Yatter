package request

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

// Read path parameter `id`
func IDOf(r *http.Request) (uint64, error) {
	ids := chi.URLParam(r, "id")

	if ids == "" {
		return 0, errors.Errorf("id was not presence")
	}

	id, err := strconv.ParseUint(ids, 10, 64)
	if err != nil {
		return 0, errors.Errorf("id was not number")
	}

	return id, nil
}

// Read path parameter `username`
func UsernameOf(r *http.Request) (string, error) {
	username := chi.URLParam(r, "username")

	if username == "" {
		return "", errors.Errorf("username was not presence")
	}
	return username, nil
}

func ParseQueries(r *http.Request) (only_media, max_id, since_id, limit *uint64, err error) {
	only_media, err = ParseQueryPointer(r.URL.Query().Get("only_media"))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	max_id, err = ParseQueryPointer(r.URL.Query().Get("max_id"))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	since_id, err = ParseQueryPointer(r.URL.Query().Get("since_id"))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	limit, err = ParseLimitQuery(r.URL.Query().Get("limit"))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return only_media, max_id, since_id, limit, nil
}

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

func ParseQueryPointer(s string) (*uint64, error) {
	parsed, err := parseQuery(s)
	if err != nil {
		return nil, err
	}

	if parsed.isEmpty {
		return nil, nil
	}
	return &parsed.id, nil
}

func ParseLimitQuery(s string) (*uint64, error) {
	parsed, err := parseQuery(s)
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

func parseQuery(s string) (parsedQuery, error) {
	if s == "" {
		return parsedQuery{isEmpty: true}, nil
	}

	if id, err := strconv.ParseUint(s, 10, 64); err != nil {
		return parsedQuery{}, err
	} else {
		return parsedQuery{id: id}, nil
	}
}
