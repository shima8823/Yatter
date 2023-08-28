package relationships

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handler request for `GET /v1/accounts/username/following`
func (h *handler) FetchFollowing(w http.ResponseWriter, r *http.Request) {
	var username string
	var account *object.Account
	var err error

	username, err = request.UsernameOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	ctx := r.Context()

	repo := h.app.Dao.Relationship()
	account, err = repo.FindAccountByUsername(ctx, username)
	if err != nil {
		httperror.NotFound(w, err)
		return
	}

	limit, err := request.ParseLimitQuery(r.URL.Query().Get("limit"))
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	accounts, err := repo.FeatchFollowing(ctx, account.ID, limit)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
