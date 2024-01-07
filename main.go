package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/controller"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/middleware"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/utils"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/view"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Use(middleware.LoadAuthContextFromCookie(app))
		redirectToLoginRoutes := e.Router.Group("", middleware.RedirectToLogin())
		mustBeLoggedInRoutes := e.Router.Group("", apis.RequireAdminOrRecordAuth())

		e.Router.GET("/", func(c echo.Context) error {
			isLoggedIn := c.Get(apis.ContextAuthRecordKey) != nil
			return utils.Render(view.Index(isLoggedIn), c)
		})

		loginHandler := controller.NewLoginHandler(app)
		e.Router.GET("/login", loginHandler.GetLogin)
		e.Router.POST("/login", loginHandler.PostLogin)
		mustBeLoggedInRoutes.POST("/logout", loginHandler.PostLogout)

		todoHandler := controller.NewTodoHandler(app)
		redirectToLoginRoutes.GET("/todos", todoHandler.GetTodos)
		mustBeLoggedInRoutes.POST("/todos/toggle/:id", todoHandler.PostToggleTodo)

		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./public"), false))
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
