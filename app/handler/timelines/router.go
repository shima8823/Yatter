package timelines

import (
	"github.com/go-chi/chi"
	"net/http"
	"yatter-backend-go/app/app"
)

type handler struct {
	app *app.App
}

// Create Handler for `/v1/timelines/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}
	r.Get("/public", h.GetTimeline)
	return r
}
