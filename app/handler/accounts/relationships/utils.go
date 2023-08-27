package relationships

import "yatter-backend-go/app/app"

type handler struct {
	app *app.App
}

func NewHandler(app *app.App) *handler {
	return &handler{app: app}
}
