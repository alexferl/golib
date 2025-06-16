package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
)

// Static holds configuration for static file serving middleware.
type Static struct {
	// Enabled indicates whether static file serving middleware is enabled.
	// Optional. Default value false.
	Enabled bool

	// Root specifies the root directory for static files.
	// Optional. Default value "".
	Root string

	// Index specifies the index file name.
	// Optional. Default value "index.html".
	Index string

	// HTML5 indicates whether to enable HTML5 mode.
	// Optional. Default value false.
	HTML5 bool

	// Browse indicates whether directory browsing is enabled.
	// Optional. Default value false.
	Browse bool

	// IgnoreBase indicates whether to ignore base path.
	// Optional. Default value false.
	IgnoreBase bool
}

// DefaultStatic provides default Static configuration.
var DefaultStatic = &Static{
	Enabled:    false,
	Root:       "",
	Index:      "index.html",
	HTML5:      false,
	Browse:     false,
	IgnoreBase: false,
}

const (
	StaticEnabled    = "static-enabled"
	StaticRoot       = "static-root"
	StaticIndex      = "static-index"
	StaticHTML5      = "static-html5"
	StaticBrowse     = "static-browse"
	StaticIgnoreBase = "static-ignore-base"
)

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (s *Static) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("Static", pflag.ExitOnError)

	fs.BoolVar(&s.Enabled, StaticEnabled, s.Enabled, "Enable static file serving middleware")
	fs.StringVar(&s.Root, StaticRoot, s.Root, "Root directory for static files")
	fs.StringVar(&s.Index, StaticIndex, s.Index, "Index file name")
	fs.BoolVar(&s.HTML5, StaticHTML5, s.HTML5, "Enable HTML5 mode")
	fs.BoolVar(&s.Browse, StaticBrowse, s.Browse, "Enable directory browsing")
	fs.BoolVar(&s.IgnoreBase, StaticIgnoreBase, s.IgnoreBase, "Ignore base path")

	return fs
}

// NewStatic creates a new static file serving middleware with the given configuration.
func NewStatic(config *Static) echo.MiddlewareFunc {
	return middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       config.Root,
		Index:      config.Index,
		HTML5:      config.HTML5,
		Browse:     config.Browse,
		IgnoreBase: config.IgnoreBase,
	})
}
