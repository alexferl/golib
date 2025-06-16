package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
)

// BodyLimit holds configuration for limiting request body size.
type BodyLimit struct {
	// Enabled indicates whether body limit middleware is enabled.
	// Optional. Default value false.
	Enabled bool

	// MaxSize specifies the maximum request body size (e.g., "1MB", "2KB").
	// Optional. Default value "1MB".
	MaxSize string
}

// DefaultBodyLimit provides default BodyLimit configuration.
var DefaultBodyLimit = &BodyLimit{
	Enabled: false,
	MaxSize: "1MB",
}

const (
	BodyLimitEnabled = "body-limit-enabled"
	BodyLimitMaxSize = "body-limit-max-size"
)

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (b *BodyLimit) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("Body Limit", pflag.ExitOnError)
	fs.BoolVar(&b.Enabled, BodyLimitEnabled, b.Enabled, "Enable request body size limiting")
	fs.StringVar(&b.MaxSize, BodyLimitMaxSize, b.MaxSize, "Maximum request body size (e.g., 1MB, 2KB)")
	return fs
}

// NewBodyLimit creates a new body limit middleware with the given configuration.
func NewBodyLimit(config *BodyLimit) echo.MiddlewareFunc {
	return middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Limit: config.MaxSize,
	})
}
