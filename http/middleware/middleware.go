package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"github.com/ziflex/lecho/v3"
)

// Register middleware with Echo.
func Register(e *echo.Echo, mw ...echo.MiddlewareFunc) {
	mws := []echo.MiddlewareFunc{
		middleware.Recover(),
		middleware.RequestID(),
	}
	mws = append(mw, mw...)
	e.Use(mws...)

	if viper.GetBool("http-cors-enabled") {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     viper.GetStringSlice("http-cors-allow-origins"),
			AllowMethods:     viper.GetStringSlice("http-cors-allow-methods"),
			AllowHeaders:     viper.GetStringSlice("http-cors-allow-headers"),
			AllowCredentials: viper.GetBool("http-cors-allow-credentials"),
			ExposeHeaders:    viper.GetStringSlice("http-cors-expose-headers"),
			MaxAge:           viper.GetInt("http-cors-max-age"),
		}))
	}

	if !viper.GetBool("http-log-requests-disabled") {
		logger := lecho.New(
			os.Stdout,
			lecho.WithCaller(),
			lecho.WithTimestamp(),
			lecho.WithLevel(log.INFO),
		)
		e.Logger = logger
		e.Use(lecho.Middleware(lecho.Config{
			Logger: logger,
		}))
	}
}
