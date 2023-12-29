package main

import (
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/models"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/views"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

const (
	TODOS = "todos"
)

func render(cmp templ.Component, c echo.Context) error {
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/", func(c echo.Context) error {
			return render(views.Index(), c)
		})

		e.Router.GET("/todos", func(c echo.Context) error {
			query := app.Dao().RecordQuery(TODOS).
				AndWhere(dbx.HashExp{"is_complete": false}).
				Limit(10)
			todos := []models.Todo{}
			if err := query.All(&todos); err != nil {
				return err
			}
			return render(views.Todos(todos), c)
		})

		e.Router.POST("/check/:id", func(c echo.Context) error {
			record, err := app.Dao().FindRecordById(TODOS, c.PathParam("id"))
			if err != nil {
				return err
			}
			var isComplete bool = record.Get("is_complete").(bool)
			record.Set("is_complete", !isComplete)
			if err := app.Dao().SaveRecord(record); err != nil {
				return err
			}
			return c.Redirect(http.StatusFound, "/todos")
		})

		// serves static files from the provided public dir
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./public"), false))
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
