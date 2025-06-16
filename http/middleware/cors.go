package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
)

// CORS holds configuration for Cross-Origin Resource Sharing.
type CORS struct {
	// Enabled indicates whether CORS middleware is enabled.
	// Optional. Default value false.
	Enabled bool

	// AllowOrigins specifies the allowed origins for CORS requests.
	// Optional. Default value []string{"*"}.
	AllowOrigins []string

	// AllowMethods specifies the allowed HTTP methods for CORS requests.
	// Optional. Default value []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}.
	AllowMethods []string

	// AllowHeaders specifies the allowed headers for CORS requests.
	// Optional. Default value []string{}.
	AllowHeaders []string

	// AllowCredentials indicates whether credentials are allowed in CORS requests.
	// Optional. Default value false.
	AllowCredentials bool

	// ExposeHeaders specifies the headers exposed to the client.
	// Optional. Default value []string{}.
	ExposeHeaders []string

	// MaxAge specifies the maximum age for preflight requests in seconds.
	// Optional. Default value 0.
	MaxAge int
}

// DefaultCORS provides default CORS configuration.
var DefaultCORS = &CORS{
	Enabled:      false,
	AllowOrigins: []string{"*"},
	AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
}

const (
	CORSEnabled          = "cors-enabled"
	CORSAllowOrigins     = "cors-allow-origins"
	CORSAllowMethods     = "cors-allow-methods"
	CORSAllowHeaders     = "cors-allow-headers"
	CORSAllowCredentials = "cors-allow-credentials"
	CORSExposeHeaders    = "cors-expose-headers"
	CORSMaxAge           = "cors-max-age"
)

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (c *CORS) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("CORS", pflag.ExitOnError)
	fs.BoolVar(&c.Enabled, CORSEnabled, c.Enabled, "Enable CORS middleware")
	fs.StringSliceVar(&c.AllowOrigins, CORSAllowOrigins, c.AllowOrigins, "Allowed origins for CORS requests")
	fs.StringSliceVar(&c.AllowMethods, CORSAllowMethods, c.AllowMethods, "Allowed HTTP methods for CORS requests")
	fs.StringSliceVar(&c.AllowHeaders, CORSAllowHeaders, c.AllowHeaders, "Allowed headers for CORS requests")
	fs.BoolVar(&c.AllowCredentials, CORSAllowCredentials, c.AllowCredentials, "Allow credentials in CORS requests")
	fs.StringSliceVar(&c.ExposeHeaders, CORSExposeHeaders, c.ExposeHeaders, "Headers exposed to the client")
	fs.IntVar(&c.MaxAge, CORSMaxAge, c.MaxAge, "Maximum age for preflight requests in seconds")
	return fs
}

// NewCORS creates a new CORS middleware with the given configuration.
func NewCORS(config *CORS) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.AllowOrigins,
		AllowMethods:     config.AllowMethods,
		AllowHeaders:     config.AllowHeaders,
		AllowCredentials: config.AllowCredentials,
		ExposeHeaders:    config.ExposeHeaders,
		MaxAge:           config.MaxAge,
	})
}
