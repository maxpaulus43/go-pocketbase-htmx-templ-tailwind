package main

import (
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/middleware"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/model"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/view"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/view/login"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/view/todos"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tokens"
)

const (
	TODOS = "todos"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Use(middleware.LoadAuthContextFromCookie(app))

		e.Router.GET("/", func(c echo.Context) error {
			isLoggedIn := c.Get(apis.ContextAuthRecordKey) != nil
			return render(view.Index(isLoggedIn), c)
		})

		e.Router.GET("/login", func(c echo.Context) error {
			return render(login.Login(), c)
		})

		e.Router.POST("/login", func(c echo.Context) error {
			// https://github.com/pocketbase/pocketbase/discussions/3447
			users, err := app.Dao().FindCollectionByNameOrId("users")
			if err != nil {
				return err
			}
			form := forms.NewRecordPasswordLogin(app, users)
			c.Bind(form)
			authRecord, err := form.Submit()
			if err != nil {
				return err
			}
			token, err := tokens.NewRecordAuthToken(app, authRecord)
			if err != nil {
				return err
			}
			c.SetCookie(&http.Cookie{
				Name:     "pb_auth",
				Value:    token,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				HttpOnly: true,
				MaxAge:   int(app.Settings().RecordAuthToken.Duration),
				Path:     "/",
			})

			return c.Redirect(http.StatusFound, "/")
		})

		redirectToLoginRoutes := e.Router.Group("", middleware.RedirectToLogin())
		redirectToLoginRoutes.GET("/todos", func(c echo.Context) error {
			query := app.Dao().RecordQuery(TODOS).Limit(10)
			todosList := []model.Todo{}
			if err := query.All(&todosList); err != nil {
				return err
			}
			return render(todos.Todos(todosList), c)
		})

		mustBeLoggedInRoutes := e.Router.Group("", apis.RequireAdminOrRecordAuth())
		mustBeLoggedInRoutes.POST("/logout", func(c echo.Context) error {
			c.SetCookie(&http.Cookie{
				Name:     "pb_auth",
				Value:    "",
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				HttpOnly: true,
				MaxAge:   0,
				Path:     "/",
			})
			return c.Redirect(http.StatusSeeOther, "/")
		})
		mustBeLoggedInRoutes.POST("/todos/toggle/:id", func(c echo.Context) error {
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

func render(cmp templ.Component, c echo.Context) error {
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
