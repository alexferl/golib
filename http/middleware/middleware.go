package middleware

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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
		var level log.Lvl

		switch strings.ToUpper(viper.GetString(config.HTTPLogRequestLevel)) {
		case "DEBUG":
			level = log.DEBUG
		case "INFO":
			level = log.INFO
		case "WARN":
			level = log.WARN
		case "ERROR":
			level = log.ERROR
		case "OFF":
			level = log.OFF
		default:
			level = log.INFO
		}

		logger := lecho.New(
			os.Stdout,
			lecho.WithCaller(),
			lecho.WithTimestamp(),
			lecho.WithLevel(level),
		)
		e.Logger = logger
		e.Use(lecho.Middleware(lecho.Config{
			Logger: logger,
		}))
	}
}
