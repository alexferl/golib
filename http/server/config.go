package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/pflag"
)

// Config holds configuration for the HTTP server.
type Config struct {
	// Name specifies the application name.
	// Optional. Default value "app".
	Name string

	// Version specifies the application version.
	// Optional. Default value "1.0.0".
	Version string

	// GracefulTimeout specifies the duration for graceful shutdown.
	// Optional. Default value 30 seconds.
	GracefulTimeout time.Duration

	// HTTP holds HTTP server configuration.
	// Optional. Default value with localhost:8080 bind address.
	HTTP HTTPConfig

	// TLS holds TLS/HTTPS server configuration.
	// Optional. Default value with TLS disabled.
	TLS TLSConfig

	// Compress holds compression configuration.
	// Optional. Default value with compression disabled.
	Compress CompressConfig

	// Redirect holds redirect configuration.
	// Optional. Default value with HTTPS redirect disabled.
	Redirect RedirectConfig

	// Healthcheck holds healthcheck configuration.
	// Optional. Default value with standard endpoints enabled.
	Healthcheck HealthcheckConfig

	// Prometheus holds Prometheus metrics configuration.
	// Optional. Default value with metrics disabled.
	Prometheus PrometheusConfig
}

// HTTPConfig holds HTTP server configuration.
type HTTPConfig struct {
	// BindAddr specifies the HTTP bind address.
	// Optional. Default value "localhost:8080".
	BindAddr string

	// IdleTimeout specifies the HTTP idle timeout.
	// Optional. Default value 60 seconds.
	IdleTimeout time.Duration

	// ReadTimeout specifies the HTTP read timeout.
	// Optional. Default value 10 seconds.
	ReadTimeout time.Duration

	// ReadHeaderTimeout specifies the HTTP read header timeout.
	// Optional. Default value 5 seconds.
	ReadHeaderTimeout time.Duration

	// WriteTimeout specifies the HTTP write timeout.
	// Optional. Default value 10 seconds.
	WriteTimeout time.Duration

	// MaxHeaderBytes specifies the maximum header bytes.
	// Optional. Default value 1MB.
	MaxHeaderBytes int
}

// TLSConfig holds TLS/HTTPS server configuration.
type TLSConfig struct {
	// Enabled indicates whether TLS/HTTPS is enabled.
	// Optional. Default value false.
	Enabled bool

	// BindAddr specifies the TLS bind address.
	// Optional. Default value "localhost:8443".
	BindAddr string

	// CertFile specifies the TLS certificate file path.
	// Optional. Default value "".
	CertFile string

	// KeyFile specifies the TLS key file path.
	// Optional. Default value "".
	KeyFile string

	// ACME holds ACME/Let's Encrypt configuration.
	// Optional. Default value with ACME disabled.
	ACME ACMEConfig
}

// ACMEConfig holds ACME/Let's Encrypt configuration.
type ACMEConfig struct {
	// Enabled indicates whether ACME/Let's Encrypt is enabled.
	// Optional. Default value false.
	Enabled bool

	// Email specifies the ACME email address.
	// Optional. Default value "".
	Email string

	// HostWhitelist specifies the ACME host whitelist.
	// Optional. Default value empty slice.
	HostWhitelist []string

	// CachePath specifies the ACME cache path.
	// Optional. Default value "./certs".
	CachePath string

	// DirectoryURL specifies the ACME directory URL.
	// Optional. Default value "".
	DirectoryURL string
}

// CompressConfig holds compression configuration.
type CompressConfig struct {
	// Enabled indicates whether compression is enabled.
	// Optional. Default value false.
	Enabled bool

	// Level specifies the compression level.
	// Optional. Default value 6.
	Level int

	// MinLength specifies the minimum length for compression.
	// Optional. Default value 1024.
	MinLength int
}

// RedirectConfig holds redirect configuration.
type RedirectConfig struct {
	// HTTPS indicates whether to redirect HTTP to HTTPS.
	// Optional. Default value false.
	HTTPS bool

	// Code specifies the redirect status code.
	// Optional. Default value 301 (Moved Permanently).
	Code int
}

// HealthcheckConfig holds healthcheck endpoint configuration.
type HealthcheckConfig struct {
	// LivenessEndpoint specifies the liveness check endpoint.
	// Optional. Default value "/livez".
	LivenessEndpoint string

	// LivenessHandler specifies the liveness check handler.
	// Optional. Default value returns 200 OK.
	LivenessHandler echo.HandlerFunc

	// ReadinessEndpoint specifies the readiness check endpoint.
	// Optional. Default value "/readyz".
	ReadinessEndpoint string

	// ReadinessHandler specifies the readiness check handler.
	// Optional. Default value returns 200 OK.
	ReadinessHandler echo.HandlerFunc

	// StartupEndpoint specifies the startup check endpoint.
	// Optional. Default value "/startupz".
	StartupEndpoint string

	// StartupHandler specifies the startup check handler.
	// Optional. Default value returns 200 OK.
	StartupHandler echo.HandlerFunc
}

func defaultHealthcheckHandler(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

// PrometheusConfig holds Prometheus metrics configuration.
type PrometheusConfig struct {
	// Enabled indicates whether Prometheus metrics are enabled.
	// Optional. Default value false.
	Enabled bool

	// Path specifies the HTTP path for Prometheus metrics.
	// Optional. Default value "/metrics".
	Path string
}

// DefaultConfig provides default server configuration.
var DefaultConfig = &Config{
	Name:            "app",
	Version:         "1.0.0",
	GracefulTimeout: 30 * time.Second,
	HTTP: HTTPConfig{
		BindAddr:          "localhost:8080",
		MaxHeaderBytes:    1 << 20, // 1MB
		IdleTimeout:       60 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	},
	TLS: TLSConfig{
		Enabled:  false,
		BindAddr: "localhost:8443",
		CertFile: "",
		KeyFile:  "",
		ACME: ACMEConfig{
			Enabled:       false,
			Email:         "",
			CachePath:     "./certs",
			HostWhitelist: []string{},
			DirectoryURL:  "",
		},
	},
	Compress: CompressConfig{
		Enabled:   false,
		Level:     6,
		MinLength: 1024,
	},
	Redirect: RedirectConfig{
		HTTPS: false,
		Code:  http.StatusMovedPermanently,
	},
	Healthcheck: HealthcheckConfig{
		LivenessEndpoint:  "/livez",
		LivenessHandler:   defaultHealthcheckHandler,
		ReadinessEndpoint: "/readyz",
		ReadinessHandler:  defaultHealthcheckHandler,
		StartupEndpoint:   "/startupz",
		StartupHandler:    defaultHealthcheckHandler,
	},
	Prometheus: PrometheusConfig{
		Enabled: false,
		Path:    "/metrics",
	},
}

const (
	ServerName                         = "server-name"
	ServerVersion                      = "server-version"
	ServerGracefulTimeout              = "server-graceful-timeout"
	ServerHTTPBindAddr                 = "server-http-bind-addr"
	ServerHTTPIdleTimeout              = "server-http-idle-timeout"
	ServerHTTPReadTimeout              = "server-http-read-timeout"
	ServerHTTPReadHeaderTimeout        = "server-http-read-header-timeout"
	ServerHTTPWriteTimeout             = "server-http-write-timeout"
	ServerHTTPMaxHeaderBytes           = "server-http-max-header-bytes"
	ServerTLSEnabled                   = "server-tls-enabled"
	ServerTLSBindAddr                  = "server-tls-bind-addr"
	ServerTLSCertFile                  = "server-tls-cert-file"
	ServerTLSKeyFile                   = "server-tls-key-file"
	ServerTLSACMEEnabled               = "server-tls-acme-enabled"
	ServerTLSACMEEmail                 = "server-tls-acme-email"
	ServerTLSACMEHostWhitelist         = "server-tls-acme-host-whitelist"
	ServerTLSACMECachePath             = "server-tls-acme-cache-path"
	ServerTLSACMEDirectoryURL          = "server-tls-acme-directory-url"
	ServerCompressEnabled              = "server-compress-enabled"
	ServerCompressLevel                = "server-compress-level"
	ServerCompressMinLength            = "server-compress-min-length"
	ServerRedirectHTTPS                = "server-redirect-https"
	ServerRedirectCode                 = "server-redirect-code"
	ServerHealthcheckLivenessEndpoint  = "server-healthcheck-liveness-endpoint"
	ServerHealthcheckReadinessEndpoint = "server-healthcheck-readiness-endpoint"
	ServerHealthcheckStartupEndpoint   = "server-healthcheck-startup-endpoint"
	ServerPrometheusEnabled            = "server-prometheus-enabled"
	ServerPrometheusPath               = "server-prometheus-path"
)

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (c *Config) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("Server", pflag.ExitOnError)

	fs.StringVar(&c.Name, ServerName, c.Name, "Application name")
	fs.StringVar(&c.Version, ServerVersion, c.Version, "Application version")
	fs.DurationVar(&c.GracefulTimeout, ServerGracefulTimeout, c.GracefulTimeout, "Graceful shutdown timeout")

	// HTTP config
	fs.StringVar(&c.HTTP.BindAddr, ServerHTTPBindAddr, c.HTTP.BindAddr, "HTTP bind address")
	fs.DurationVar(&c.HTTP.IdleTimeout, ServerHTTPIdleTimeout, c.HTTP.IdleTimeout, "HTTP idle timeout")
	fs.DurationVar(&c.HTTP.ReadTimeout, ServerHTTPReadTimeout, c.HTTP.ReadTimeout, "HTTP read timeout")
	fs.DurationVar(&c.HTTP.ReadHeaderTimeout, ServerHTTPReadHeaderTimeout, c.HTTP.ReadHeaderTimeout, "HTTP read header timeout")
	fs.DurationVar(&c.HTTP.WriteTimeout, ServerHTTPWriteTimeout, c.HTTP.WriteTimeout, "HTTP write timeout")
	fs.IntVar(&c.HTTP.MaxHeaderBytes, ServerHTTPMaxHeaderBytes, c.HTTP.MaxHeaderBytes, "HTTP max header bytes")

	// TLS config
	fs.BoolVar(&c.TLS.Enabled, ServerTLSEnabled, c.TLS.Enabled, "Enable TLS/HTTPS")
	fs.StringVar(&c.TLS.BindAddr, ServerTLSBindAddr, c.TLS.BindAddr, "TLS bind address")
	fs.StringVar(&c.TLS.CertFile, ServerTLSCertFile, c.TLS.CertFile, "TLS certificate file")
	fs.StringVar(&c.TLS.KeyFile, ServerTLSKeyFile, c.TLS.KeyFile, "TLS key file")

	// ACME config
	fs.BoolVar(&c.TLS.ACME.Enabled, ServerTLSACMEEnabled, c.TLS.ACME.Enabled, "Enable ACME/Let's Encrypt")
	fs.StringVar(&c.TLS.ACME.Email, ServerTLSACMEEmail, c.TLS.ACME.Email, "ACME email address")
	fs.StringSliceVar(&c.TLS.ACME.HostWhitelist, ServerTLSACMEHostWhitelist, c.TLS.ACME.HostWhitelist, "ACME host whitelist")
	fs.StringVar(&c.TLS.ACME.CachePath, ServerTLSACMECachePath, c.TLS.ACME.CachePath, "ACME cache path")
	fs.StringVar(&c.TLS.ACME.DirectoryURL, ServerTLSACMEDirectoryURL, c.TLS.ACME.DirectoryURL, "ACME directory URL")

	// Compression config
	fs.BoolVar(&c.Compress.Enabled, ServerCompressEnabled, c.Compress.Enabled, "Enable compression")
	fs.IntVar(&c.Compress.Level, ServerCompressLevel, c.Compress.Level, "Compression level")
	fs.IntVar(&c.Compress.MinLength, ServerCompressMinLength, c.Compress.MinLength, "Minimum length for compression")

	// Redirect config
	fs.BoolVar(&c.Redirect.HTTPS, ServerRedirectHTTPS, c.Redirect.HTTPS, "Redirect HTTP to HTTPS")
	fs.IntVar(&c.Redirect.Code, ServerRedirectCode, c.Redirect.Code, "Redirect status code")

	// Healthcheck config
	fs.StringVar(&c.Healthcheck.LivenessEndpoint, ServerHealthcheckLivenessEndpoint, c.Healthcheck.LivenessEndpoint, "Liveness check endpoint")
	fs.StringVar(&c.Healthcheck.ReadinessEndpoint, ServerHealthcheckReadinessEndpoint, c.Healthcheck.ReadinessEndpoint, "Readiness check endpoint")
	fs.StringVar(&c.Healthcheck.StartupEndpoint, ServerHealthcheckStartupEndpoint, c.Healthcheck.StartupEndpoint, "Startup check endpoint")

	// Prometheus config
	fs.BoolVar(&c.Prometheus.Enabled, ServerPrometheusEnabled, c.Prometheus.Enabled, "Enable Prometheus metrics")
	fs.StringVar(&c.Prometheus.Path, ServerPrometheusPath, c.Prometheus.Path, "Prometheus metrics endpoint path")

	return fs
}
