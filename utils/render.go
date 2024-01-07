package utils

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

func Render(cmp templ.Component, c echo.Context) error {
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
