package relationships

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

// Handler request for `GET /v1/accounts/relationships`
func (h *handler) FetchList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	account := auth.AccountOf(r)
	if account == nil {
		httperror.Error(w, http.StatusUnauthorized)
		return
	}

	repo := h.app.Dao.Relationship()
	relationships, err := repo.FeatchRelationships(ctx, account.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relationships); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
