package statuses

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handler request for `GET /v1/statuses/id`
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := request.IDOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if objAccount, err := h.app.Dao.Status().Retrieve(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			httperror.NotFound(w, id)
			return
		}
		httperror.InternalServerError(w, err)
		return
	} else if objAccount != nil {
		if err := json.NewEncoder(w).Encode(objAccount); err != nil {
			httperror.InternalServerError(w, err)
			return
		}
	} else {
		httperror.NotFound(w, id)
		return
	}
}
