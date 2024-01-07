package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/utils"
	"github.com/maxpaulus43/go-pocketbase-htmx-templ-tailwind/view/login"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tokens"
)

type LoginHandler struct{ Handler }

func NewLoginHandler(app *pocketbase.PocketBase) LoginHandler {
	return LoginHandler{NewHandler(app)}
}

func (h LoginHandler) GetLogin(c echo.Context) error {
	return utils.Render(login.Login(), c)
}

func (h LoginHandler) PostLogin(c echo.Context) error {
	// https://github.com/pocketbase/pocketbase/discussions/3447
	users, err := h.app.Dao().FindCollectionByNameOrId("users")
	// TODO allow admins to login too
	if err != nil {
		return utils.Render(login.LoginWithError(fmt.Sprint(err)), c)
	}
	form := forms.NewRecordPasswordLogin(h.app, users)
	c.Bind(form)
	authRecord, err := form.Submit()
	if err != nil {
		return utils.Render(login.LoginWithError(fmt.Sprint(err)), c)
	}
	token, err := tokens.NewRecordAuthToken(h.app, authRecord)
	if err != nil {
		fmt.Println(err)
		return utils.Render(login.LoginWithError(fmt.Sprint(err)), c)
	}
	c.SetCookie(&http.Cookie{
		Name:     "pb_auth",
		Value:    token,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		MaxAge:   int(h.app.Settings().RecordAuthToken.Duration),
		Path:     "/",
	})

	return c.Redirect(http.StatusFound, "/")

}

func (h LoginHandler) PostLogout(c echo.Context) error {
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
}
