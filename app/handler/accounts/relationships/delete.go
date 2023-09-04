package relationships

import (
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handler request for `POST /v1/accounts/username/unfollow`
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err = h.app.Dao.Relationship().Delete(ctx, followingAccount.ID, followerAccount.ID); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
