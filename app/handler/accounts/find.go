package accounts

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
)

// Handler request for `GET /v1/accounts/username`
func (h *handler) FindUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := chi.URLParam(r, "username")

	w.Header().Set("Content-Type", "application/json")
	if objAccount, err := h.app.Dao.Account().FindByUsername(ctx, username); err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if objAccount != nil {
		if err := json.NewEncoder(w).Encode(objAccount); err != nil {
			httperror.InternalServerError(w, err)
			return
		}
	} else {
		httperror.NotFound(w, username)
		return
	}
}
