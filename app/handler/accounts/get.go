package accounts

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handler request for `GET /v1/accounts/username`
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username, err := request.UsernameOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if objAccount, err := h.app.Dao.Account().Retrieve(ctx, username); err != nil {
		if err == sql.ErrNoRows {
			httperror.NotFound(w, username)
			return
		}
		httperror.InternalServerError(w, err)
	} else if objAccount != nil {
		if err := json.NewEncoder(w).Encode(objAccount); err != nil {
			httperror.InternalServerError(w, err)
			return
		}
	}
}
