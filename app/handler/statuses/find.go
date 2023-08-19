package statuses

import (
	"encoding/json"
	"net/http"
	"strconv"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

// Handler request for `GET /v1/statuses/id`
func (h *handler) FindStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	strId := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if objAccount, err := h.app.Dao.Status().FindByID(ctx, id); err != nil {
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
