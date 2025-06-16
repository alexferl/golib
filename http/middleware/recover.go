package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
)

// Recover holds configuration for panic recovery middleware.
type Recover struct {
	// Enabled indicates whether recovery middleware is enabled.
	// Optional. Default value true.
	Enabled bool

	// Size of the stack to be printed.
	// Optional. Default value 4KB.
	StackSize int

	// DisableStackAll disables formatting stack traces of all other goroutines
	// into buffer after the trace for the current goroutine.
	// Optional. Default value false.
	DisableStackAll bool

	// DisablePrintStack disables printing stack trace.
	// Optional. Default value as false.
	DisablePrintStack bool

	// DisableErrorHandler disables the call to centralized HTTPErrorHandler.
	// The recovered error is then passed back to upstream middleware, instead of swallowing the error.
	// Optional. Default value false.
	DisableErrorHandler bool
}

// DefaultRecover provides default Recover configuration.
var DefaultRecover = &Recover{
	Enabled:             true,
	StackSize:           4 << 10, // 4 KB
	DisableStackAll:     false,
	DisablePrintStack:   false,
	DisableErrorHandler: false,
}

const (
	RecoverEnabled             = "recover-enabled"
	RecoverStackSize           = "recover-stack-size"
	RecoverDisableStackAll     = "recover-disable-stack-all"
	RecoverDisablePrintStack   = "recover-disable-print-stack"
	RecoverDisableErrorHandler = "recover-disable-error-handler"
)

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (r *Recover) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("Recover", pflag.ExitOnError)

	fs.BoolVar(&r.Enabled, RecoverEnabled, r.Enabled, "Enable panic recovery middleware")
	fs.IntVar(&r.StackSize, RecoverStackSize, r.StackSize, "Stack size for recovery in bytes")
	fs.BoolVar(&r.DisableStackAll, RecoverDisableStackAll, r.DisableStackAll, "Disable stack trace for all errors")
	fs.BoolVar(&r.DisablePrintStack, RecoverDisablePrintStack, r.DisablePrintStack, "Disable printing stack trace")
	fs.BoolVar(&r.DisableErrorHandler, RecoverDisableErrorHandler, r.DisableErrorHandler, "Disable custom error handler")

	return fs
}

// NewRecover creates a new panic recovery middleware with the given configuration.
func NewRecover(config *Recover) echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:           config.StackSize,
		DisableStackAll:     config.DisableStackAll,
		DisablePrintStack:   config.DisablePrintStack,
		DisableErrorHandler: config.DisableErrorHandler,
	})
}
