package handler

import "github.com/labstack/echo/v4"

func HTTPError(c echo.Context, status int, desc string) error {
	return c.JSON(status, echo.HTTPError{
		Code:     status,
		Message:  desc,
		Internal: nil,
	})
}
