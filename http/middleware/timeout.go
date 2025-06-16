package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
)

// Timeout holds configuration for request timeout middleware.
type Timeout struct {
	// Enabled indicates whether timeout middleware is enabled.
	// Optional. Default value true.
	Enabled bool

	// ErrorMessage specifies the error message for timeout responses.
	// Optional. Default value "Request timeout".
	ErrorMessage string

	// Duration specifies the timeout duration.
	// Optional. Default value 15 seconds.
	Duration time.Duration
}

// DefaultTimeout provides default Timeout configuration.
var DefaultTimeout = &Timeout{
	Enabled:      true,
	ErrorMessage: "Request timeout",
	Duration:     15 * time.Second,
}

const (
	TimeoutEnabled      = "timeout-enabled"
	TimeoutErrorMessage = "timeout-error-message"
	TimeoutDuration     = "timeout-duration"
)

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (t *Timeout) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("Timeout", pflag.ExitOnError)
	fs.BoolVar(&t.Enabled, TimeoutEnabled, t.Enabled, "Enable request timeout middleware")
	fs.StringVar(&t.ErrorMessage, TimeoutErrorMessage, t.ErrorMessage, "Error message for timeout responses")
	fs.DurationVar(&t.Duration, TimeoutDuration, t.Duration, "Request timeout duration")
	return fs
}

// NewTimeout creates a new request timeout middleware with the given configuration.
func NewTimeout(config *Timeout) echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		ErrorMessage: config.ErrorMessage,
		Timeout:      config.Duration,
	})
}
