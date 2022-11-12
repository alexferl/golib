package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/ziflex/lecho/v3"

	"github.com/alexferl/golib/http/config"
)

// Register middleware with Echo.
func Register(e *echo.Echo, mw ...echo.MiddlewareFunc) {
	mws := []echo.MiddlewareFunc{
		middleware.Recover(),
		middleware.RequestID(),
	}
	mws = append(mw, mw...)
	e.Use(mws...)

	if viper.GetBool(config.HTTPCORSEnabled) {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     viper.GetStringSlice(config.HTTPCORSAllowOrigins),
			AllowMethods:     viper.GetStringSlice(config.HTTPCORSAllowHeaders),
			AllowHeaders:     viper.GetStringSlice(config.HTTPCORSAllowHeaders),
			AllowCredentials: viper.GetBool(config.HTTPCORSAllowCredentials),
			ExposeHeaders:    viper.GetStringSlice(config.HTTPCORSExposeHeaders),
			MaxAge:           viper.GetInt(config.HTTPCORSMaxAge),
		}))
	}

	if viper.GetBool(config.HTTPLogRequests) {
		logger := lecho.From(log.Logger)
		e.Logger = logger
		e.Use(lecho.Middleware(lecho.Config{
			Logger: logger,
		}))
	}
}
