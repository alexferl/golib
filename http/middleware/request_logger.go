package middleware

import (
	"strconv"
	"time"

	"github.com/alexferl/golib/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
)

// RequestLogger holds configuration for logging middleware.
type RequestLogger struct {
	// Enabled indicates whether logging middleware is enabled.
	// Optional. Default value false.
	Enabled bool

	// Logger instance from the logger submodule.
	// Optional. If nil, a default logger will be created.
	Logger *logger.Logger
}

// DefaultLogger provides default RequestLogger configuration.
var DefaultLogger = &RequestLogger{
	Enabled: false,
	Logger:  nil,
}

const (
	RequestLoggerEnabled = "request-logger-enabled"
)

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (l *RequestLogger) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("Request Logger", pflag.ExitOnError)

	fs.BoolVar(&l.Enabled, RequestLoggerEnabled, l.Enabled, "Enable request logging middleware")

	return fs
}

// NewRequestLogger creates a new request logging middleware with the given configuration.
func NewRequestLogger(config *RequestLogger) echo.MiddlewareFunc {
	var log *logger.Logger
	if config.Logger != nil {
		log = config.Logger
	} else {
		defaultLogger, err := logger.New(logger.DefaultConfig)
		if err != nil {
			panic(err)
		}
		log = defaultLogger
	}

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		HandleError:      true,
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
			log.Info().
				Str("time", time.Now().Format(time.RFC3339Nano)).
				Str("id", v.RequestID).
				Str("remote_ip", v.RemoteIP).
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
}
