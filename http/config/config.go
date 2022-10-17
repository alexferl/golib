package config

import (
	"net"
	"net/http"

	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
)

// Config holds all HTTP configuration.
type Config struct {
	BindAddress         net.IP
	BindPort            uint
	CORS                middleware.CORSConfig
	CORSEnabled         bool
	GracefulTimeout     uint
	LogRequestsDisabled bool
}

var DefaultConfig = &Config{
	BindAddress: net.ParseIP("127.0.0.1"),
	BindPort:    1323,
	CORSEnabled: false,
	CORS: middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		},
		AllowHeaders:     []string{},
		AllowCredentials: false,
		ExposeHeaders:    []string{},
		MaxAge:           0,
	},
	GracefulTimeout:     30,
	LogRequestsDisabled: false,
}

// BindFlags adds all the flags from the command line.
func (c *Config) BindFlags(fs *pflag.FlagSet) {
	fs.IPVar(&c.BindAddress, "http-bind-address", c.BindAddress, "The IP address to listen at.")
	fs.UintVar(&c.BindPort, "http-bind-port", c.BindPort, "The port to listen at.")
	fs.UintVar(&c.GracefulTimeout, "http-graceful-timeout", c.GracefulTimeout,
		"Timeout for graceful shutdown.")
	fs.BoolVar(&c.LogRequestsDisabled, "http-log-requests-disabled", c.LogRequestsDisabled,
		"Disable the logging of HTTP requests")
	fs.BoolVar(&c.CORSEnabled, "http-cors-enabled", c.CORSEnabled, "Enable cross-origin resource sharing.")
	fs.StringSliceVar(&c.CORS.AllowOrigins, "http-cors-allow-origins", c.CORS.AllowOrigins,
		"Indicates whether the response can be shared with requesting code from the given origin.")
	fs.StringSliceVar(&c.CORS.AllowMethods, "http-cors-allow-methods", c.CORS.AllowMethods,
		"Indicates which HTTP methods are allowed for cross-origin requests.")
	fs.StringSliceVar(&c.CORS.AllowHeaders, "http-cors-allow-headers", c.CORS.AllowHeaders,
		"Indicate which HTTP headers can be used during an actual request.")
	fs.BoolVar(&c.CORS.AllowCredentials, "http-cors-allow-credentials", c.CORS.AllowCredentials,
		"Tells browsers whether to expose the response to frontend JavaScript code when the request's "+
			"credentials mode (Request.credentials) is 'include'.")
	fs.StringSliceVar(&c.CORS.ExposeHeaders, "http-cors-expose-headers", c.CORS.ExposeHeaders,
		"Indicates which headers can be exposed as part of the response by listing their name.")
	fs.IntVar(&c.CORS.MaxAge, "http-cors-max-age", c.CORS.MaxAge,
		"Indicates how long the results of a preflight request can be cached.")
}
