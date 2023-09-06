package statuses

import (
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handler request for `DELETE /v1/statuses/id`
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := request.IDOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if _, err := h.app.Dao.Status().Retrieve(ctx, id); err != nil {
		httperror.NotFound(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := h.app.Dao.Status().Delete(ctx, id); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
