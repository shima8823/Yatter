package statuses

import (
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"yatter-backend-go/app/handler/httperror"
)

// Handler request for `DELETE /v1/statuses/id`
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	strId := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if _, err := h.app.Dao.Status().FindByID(ctx, id); err != nil {
		httperror.NotFound(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := h.app.Dao.Status().DeleteByID(ctx, id); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
