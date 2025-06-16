package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
)

// RequestID holds configuration for request ID middleware.
type RequestID struct {
	// Enabled indicates whether request ID middleware is enabled.
	// Optional. Default value false.
	Enabled bool

	// TargetHeader specifies the header name for the request ID.
	// Optional. Default value "X-Request-ID".
	TargetHeader string
}

// DefaultRequestID provides default RequestID configuration.
var DefaultRequestID = &RequestID{
	Enabled:      false,
	TargetHeader: "X-Request-ID",
}

const (
	RequestIDEnabled      = "request-id-enabled"
	RequestIDTargetHeader = "request-id-target-header"
)

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (r *RequestID) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("Request ID", pflag.ExitOnError)
	fs.BoolVar(&r.Enabled, RequestIDEnabled, r.Enabled, "Enable request ID middleware")
	fs.StringVar(&r.TargetHeader, RequestIDTargetHeader, r.TargetHeader, "Header name for the request ID")
	return fs
}

// NewRequestID creates a new request ID middleware with the given configuration.
func NewRequestID(config *RequestID) echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		TargetHeader: config.TargetHeader,
	})
}
