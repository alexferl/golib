package middleware

import (
	"context"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type OpenAPIConfig struct {
	File       string
	ContextKey string
}

var DefaultOpenAPIConfig = OpenAPIConfig{
	ContextKey: "validator",
}

type Error struct {
	Message string `json:"error" xml:"error"`
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
				log.Error().Msgf("error loading schema file: %v", err)
				return err
			}

			err = schema.Validate(ctx)
			if err != nil {
				log.Error().Msgf("error validating schema: %v", err)
				return err
			}

			r, err := gorillamux.NewRouter(schema)
			if err != nil {
				log.Error().Msgf("error creating router: %v", err)
				return err
			}

			route, pathParams, err := r.FindRoute(c.Request())
			if err != nil {
				log.Debug().Msgf("error finding route: %v", err)

				if err == routers.ErrPathNotFound {
					m := &Error{Message: "route not found"}
					return c.JSON(http.StatusNotFound, m)
				}

				if err == routers.ErrMethodNotAllowed {
					m := &Error{Message: "method not allowed"}
					return c.JSON(http.StatusMethodNotAllowed, m)
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
				log.Debug().Msgf("error validating request: %v", err)
				return err
			}

			c.Set(config.ContextKey, requestValidationInput)

			return next(c)
		}
	}
}
