package controller

import (
	"github.com/pocketbase/pocketbase"
)

type Handler struct {
	app *pocketbase.PocketBase
}

func NewHandler(app *pocketbase.PocketBase) Handler {
	return Handler{app: app}
}
