package accounts

import (
	"net/http"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler/accounts/relationships"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi"
)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/accounts/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	accoutnHandler := &handler{app: app}
	relationshipHandler := relationships.NewHandler(app)
	r.Post("/", accoutnHandler.Create)
	r.Get("/{username}", accoutnHandler.FindUser)

	// Relationship
	r.With(auth.Middleware(app)).Post("/{username}/follow", relationshipHandler.Create)
	r.With(auth.Middleware(app)).Post("/{username}/unfollow", relationshipHandler.Delete)
	r.With(auth.Middleware(app)).Get("/relationships", relationshipHandler.FetchList)

	r.Get("/{username}/following", relationshipHandler.FetchFollowing)
	return r
}
