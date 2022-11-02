package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/alexferl/golib/http/router"
)

type Handler interface {
	GetRoutes() []*router.Route
}

func JSONError(c echo.Context, status int, msg string) error {
	return c.JSON(status, echo.HTTPError{
		Code:    status,
		Message: msg,
	})
}
