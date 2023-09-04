package relationships

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handler request for `GET /v1/accounts/username/followers`
func (h *handler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	var username string
	var account *object.Account
	var err error

	username, err = request.UsernameOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	ctx := r.Context()

	account, err = h.app.Dao.Account().Retrieve(ctx, username)
	if err != nil {
		httperror.NotFound(w, err)
		return
	}

	only_media, max_id, since_id, limit, err := request.ParseQueries(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	accounts, err := h.app.Dao.Relationship().RetrieveFollowers(ctx, account.ID, only_media, max_id, since_id, limit)
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
