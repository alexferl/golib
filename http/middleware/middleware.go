package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/alexferl/golib/http/config"
)

// Register middleware with Echo.
func Register(e *echo.Echo, mw ...echo.MiddlewareFunc) {
	mws := []echo.MiddlewareFunc{
		middleware.Recover(),
		middleware.RequestID(),
	}

	if viper.GetBool(config.HTTPCORSEnabled) {
		cors := middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     viper.GetStringSlice(config.HTTPCORSAllowOrigins),
			AllowMethods:     viper.GetStringSlice(config.HTTPCORSAllowHeaders),
			AllowHeaders:     viper.GetStringSlice(config.HTTPCORSAllowHeaders),
			AllowCredentials: viper.GetBool(config.HTTPCORSAllowCredentials),
			ExposeHeaders:    viper.GetStringSlice(config.HTTPCORSExposeHeaders),
			MaxAge:           viper.GetInt(config.HTTPCORSMaxAge),
		})
		mws = append(mws, cors)
	}

	if viper.GetBool(config.HTTPLogRequests) {
		logger := middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogRequestID:     true,
			LogRemoteIP:      true,
			LogHost:          true,
			LogMethod:        true,
			LogURI:           true,
			LogUserAgent:     true,
			LogStatus:        true,
			LogError:         true,
			LogLatency:       true,
			LogContentLength: true,
			LogResponseSize:  true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				i, _ := strconv.Atoi(v.ContentLength)
				log.Logger.Info().
					Str("time", time.Now().Format(time.RFC3339Nano)).
					Str("id", v.RequestID).
					Str("remote_id", v.RemoteIP).
					Str("host", v.Host).
					Str("method", v.Method).
					Str("uri", v.URI).
					Str("user_agent", v.UserAgent).
					Int("status", v.Status).
					Err(v.Error).
					Int64("latency", v.Latency.Nanoseconds()).
					Str("latency_human", v.Latency.String()).
					Int64("bytes_in", int64(i)).
					Int64("bytes_out", v.ResponseSize).
					Send()

				return nil
			},
		})

		mws = append(mws, logger)
	}

	mws = append(mws, mw...)

	e.Use(mws...)
}
