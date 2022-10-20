package middleware

import (
	"context"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/labstack/echo/v4"

	"github.com/alexferl/golib/http/handler"
)

type OpenAPIConfig struct {
	File       string
	ContextKey string
}

var DefaultOpenAPIConfig = OpenAPIConfig{
	ContextKey: "validator",
}

func OpenAPI(file string) echo.MiddlewareFunc {
	c := DefaultOpenAPIConfig
	c.File = file
	return OpenAPIWithConfig(c)
}

func OpenAPIWithConfig(config OpenAPIConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.Background()
			loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
			schema, err := loader.LoadFromFile(config.File)
			if err != nil {
				c.Logger().Errorf("error loading schema file: %v", err)
				return err
			}

			err = schema.Validate(ctx)
			if err != nil {
				c.Logger().Errorf("error validating schema: %v", err)
				return err
			}

			r, err := gorillamux.NewRouter(schema)
			if err != nil {
				c.Logger().Errorf("error creating router: %v", err)
				return err
			}

			route, pathParams, err := r.FindRoute(c.Request())
			if err != nil {
				c.Logger().Debugf("error finding route for %s: %v", c.Request().URL.String(), err)

				if err == routers.ErrPathNotFound {
					return handler.HTTPError(c, http.StatusNotFound, "route not found")
				}

				if err == routers.ErrMethodNotAllowed {
					return handler.HTTPError(c, http.StatusMethodNotAllowed, "method not allowed")
				}
				return err
			}

			requestValidationInput := &openapi3filter.RequestValidationInput{
				Request:    c.Request(),
				PathParams: pathParams,
				Route:      route,
			}
			err = openapi3filter.ValidateRequest(ctx, requestValidationInput)
			if err != nil {
				c.Logger().Debugf("error validating request: %v", err)
				return err
			}

			c.Set(config.ContextKey, requestValidationInput)

			return next(c)
		}
	}
}
