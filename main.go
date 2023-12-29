package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/views"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/", func(c echo.Context) error {
			return views.Index().Render(c.Request().Context(), c.Response().Writer)
		})

		// serves static files from the provided public dir
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./public"), false))
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
