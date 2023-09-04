package relationships

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

type AddRequest struct {
	followingID string
	followerID  string
}

// Handler request for `POST /v1/accounts/{username}/follow`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	followingAccount := auth.AccountOf(r)
	if followingAccount == nil {
		httperror.Error(w, http.StatusUnauthorized)
		return
	}
	var followerUsername string
	var followerAccount *object.Account
	var err error

	followerUsername, err = request.UsernameOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	ctx := r.Context()

	followerAccount, err = h.app.Dao.Account().Retrieve(ctx, followerUsername)
	if err != nil {
		httperror.NotFound(w, err)
		return
	}

	relationship := new(object.Relationship)
	relationship.FollowingId = followingAccount.ID
	relationship.FollowerId = followerAccount.ID

	if followingAccount.ID == followerAccount.ID {
		httperror.Error(w, http.StatusBadRequest)
		return
	}

	if err = h.app.Dao.Relationship().Create(ctx, followingAccount.ID, followerAccount.ID); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relationship); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
