package timelines

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handler request for `GET /v1/timelines/home`
func (h *handler) GetHome(w http.ResponseWriter, r *http.Request) {
	account := auth.AccountOf(r)
	if account == nil {
		httperror.Error(w, http.StatusUnauthorized)
		return
	}
	ctx := r.Context()

	only_media, max_id, since_id, limit, err := request.ParseQueries(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if objStatuses, err := h.app.Dao.Status().HomeTimeline(ctx, account.ID, only_media, max_id, since_id, limit); err != nil {
		httperror.InternalServerError(w, err)
	} else if objStatuses != nil {
		if err := json.NewEncoder(w).Encode(objStatuses); err != nil {
			httperror.InternalServerError(w, err)
		}
	} else {
		httperror.NotFound(w, "timeline")
	}
}
