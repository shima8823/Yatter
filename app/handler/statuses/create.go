package statuses

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

type AddRequest struct {
	Status string
}

// Handle request for `POST /v1/statuses`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	account := auth.AccountOf(r)
	if account == nil {
		httperror.Error(w, http.StatusUnauthorized)
		return
	}
	ctx := r.Context()

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	status := new(object.Status)
	status.AccountId = account.ID
	status.Content = req.Status

	repo := h.app.Dao.Status()
	if err := repo.Create(ctx, status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
