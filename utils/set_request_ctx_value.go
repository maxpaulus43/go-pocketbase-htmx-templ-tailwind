package utils

import (
	"context"

	"github.com/labstack/echo/v5"
)

func SetRequestContextValue(c echo.Context, key string, val any) {
	c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), key, val)))
}
