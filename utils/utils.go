package utils

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/search"
)

func Render(cmp templ.Component, c echo.Context) error {
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func NewRuleQuery(app core.App, c echo.Context, coll models.Collection, rule string) (*dbx.SelectQuery, error) {
	info := apis.RequestInfo(c)
	query := app.Dao().RecordQuery(coll.Name)
	resolver := resolvers.NewRecordFieldResolver(app.Dao(), &coll, info, true)
	expr, err := search.FilterData(rule).BuildExpr(resolver)
	if err != nil {
		return nil, apis.NewBadRequestError("", err)
	}
	resolver.UpdateQuery(query)
	query.AndWhere(expr)
	return query, nil
}

func SetRequestContextValue(c echo.Context, key string, val any) {
	c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), key, val)))
}
