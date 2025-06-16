package middleware

import (
	secure "github.com/alexferl/echo-secure"
	"github.com/labstack/echo/v4"
	"github.com/spf13/pflag"
)

// StrictTransportSecurity holds configuration for HSTS.
type StrictTransportSecurity struct {
	// MaxAge specifies the max age for HSTS in seconds.
	// Optional. Default value from secure.DefaultConfig.
	MaxAge int

	// ExcludeSubdomains indicates whether to exclude subdomains from HSTS.
	// Optional. Default value from secure.DefaultConfig.
	ExcludeSubdomains bool

	// PreloadEnabled indicates whether HSTS preload is enabled.
	// Optional. Default value from secure.DefaultConfig.
	PreloadEnabled bool
}

// Secure holds configuration for security middleware.
type Secure struct {
	// Enabled indicates whether security middleware is enabled.
	// Optional. Default value false.
	Enabled bool

	// ContentSecurityPolicy specifies the CSP header value.
	// Optional. Default value from secure.DefaultConfig.
	ContentSecurityPolicy string

	// ContentSecurityPolicyReportOnly indicates whether CSP is in report-only mode.
	// Optional. Default value from secure.DefaultConfig.
	ContentSecurityPolicyReportOnly bool

	// CrossOriginEmbedderPolicy specifies the COEP header value.
	// Optional. Default value from secure.DefaultConfig.
	CrossOriginEmbedderPolicy string

	// CrossOriginOpenerPolicy specifies the COOP header value.
	// Optional. Default value from secure.DefaultConfig.
	CrossOriginOpenerPolicy string

	// CrossOriginResourcePolicy specifies the CORP header value.
	// Optional. Default value from secure.DefaultConfig.
	CrossOriginResourcePolicy string

	// PermissionsPolicy specifies the Permissions-Policy header value.
	// Optional. Default value from secure.DefaultConfig.
	PermissionsPolicy string

	// ReferrerPolicy specifies the Referrer-Policy header value.
	// Optional. Default value from secure.DefaultConfig.
	ReferrerPolicy string

	// Server specifies the Server header value.
	// Optional. Default value from secure.DefaultConfig.
	Server string

	// StrictTransportSecurity holds HSTS configuration.
	// Optional. Default value from secure.DefaultConfig.
	StrictTransportSecurity StrictTransportSecurity

	// XContentTypeOptions specifies the X-Content-Type-Options header value.
	// Optional. Default value from secure.DefaultConfig.
	XContentTypeOptions string

	// XFrameOptions specifies the X-Frame-Options header value.
	// Optional. Default value from secure.DefaultConfig.
	XFrameOptions string
}

// DefaultSecure provides default Secure configuration.
var DefaultSecure = &Secure{
	Enabled:                         false,
	ContentSecurityPolicy:           secure.DefaultConfig.ContentSecurityPolicy,
	ContentSecurityPolicyReportOnly: secure.DefaultConfig.ContentSecurityPolicyReportOnly,
	CrossOriginEmbedderPolicy:       secure.DefaultConfig.CrossOriginEmbedderPolicy,
	CrossOriginOpenerPolicy:         secure.DefaultConfig.CrossOriginOpenerPolicy,
	CrossOriginResourcePolicy:       secure.DefaultConfig.CrossOriginResourcePolicy,
	PermissionsPolicy:               secure.DefaultConfig.PermissionsPolicy,
	ReferrerPolicy:                  secure.DefaultConfig.ReferrerPolicy,
	Server:                          secure.DefaultConfig.Server,
	StrictTransportSecurity: StrictTransportSecurity{
		MaxAge:            secure.DefaultConfig.StrictTransportSecurity.MaxAge,
		ExcludeSubdomains: secure.DefaultConfig.StrictTransportSecurity.ExcludeSubdomains,
		PreloadEnabled:    secure.DefaultConfig.StrictTransportSecurity.PreloadEnabled,
	},
	XContentTypeOptions: secure.DefaultConfig.XContentTypeOptions,
	XFrameOptions:       secure.DefaultConfig.XFrameOptions,
}

const (
	SecureEnabled                         = "secure-enabled"
	SecureContentSecurityPolicy           = "secure-content-security-policy"
	SecureContentSecurityPolicyReportOnly = "secure-content-security-policy-report-only"
	SecureCrossOriginEmbedderPolicy       = "secure-cross-origin-embedder-policy"
	SecureCrossOriginOpenerPolicy         = "secure-cross-origin-opener-policy"
	SecureCrossOriginResourcePolicy       = "secure-cross-origin-resource-policy"
	SecurePermissionsPolicy               = "secure-permissions-policy"
	SecureReferrerPolicy                  = "secure-referrer-policy"
	SecureServer                          = "secure-server"
	SecureSTSMaxAge                       = "secure-sts-max-age"
	SecureSTSExcludeSubdomains            = "secure-sts-exclude-subdomains"
	SecureSTSPreloadEnabled               = "secure-sts-preload-enabled"
	SecureXContentTypeOptions             = "secure-x-content-type-options"
	SecureXFrameOptions                   = "secure-x-frame-options"
)

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (s *Secure) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("Secure", pflag.ExitOnError)

	fs.BoolVar(&s.Enabled, SecureEnabled, s.Enabled, "Enable security middleware")
	fs.StringVar(&s.ContentSecurityPolicy, SecureContentSecurityPolicy, s.ContentSecurityPolicy, "Content Security Policy header value")
	fs.BoolVar(&s.ContentSecurityPolicyReportOnly, SecureContentSecurityPolicyReportOnly, s.ContentSecurityPolicyReportOnly, "Enable CSP report-only mode")
	fs.StringVar(&s.CrossOriginEmbedderPolicy, SecureCrossOriginEmbedderPolicy, s.CrossOriginEmbedderPolicy, "Cross-Origin-Embedder-Policy header value")
	fs.StringVar(&s.CrossOriginOpenerPolicy, SecureCrossOriginOpenerPolicy, s.CrossOriginOpenerPolicy, "Cross-Origin-Opener-Policy header value")
	fs.StringVar(&s.CrossOriginResourcePolicy, SecureCrossOriginResourcePolicy, s.CrossOriginResourcePolicy, "Cross-Origin-Resource-Policy header value")
	fs.StringVar(&s.PermissionsPolicy, SecurePermissionsPolicy, s.PermissionsPolicy, "Permissions-Policy header value")
	fs.StringVar(&s.ReferrerPolicy, SecureReferrerPolicy, s.ReferrerPolicy, "Referrer-Policy header value")
	fs.StringVar(&s.Server, SecureServer, s.Server, "Server header value")
	fs.IntVar(&s.StrictTransportSecurity.MaxAge, SecureSTSMaxAge, s.StrictTransportSecurity.MaxAge, "HSTS max age in seconds")
	fs.BoolVar(&s.StrictTransportSecurity.ExcludeSubdomains, SecureSTSExcludeSubdomains, s.StrictTransportSecurity.ExcludeSubdomains, "Exclude subdomains from HSTS")
	fs.BoolVar(&s.StrictTransportSecurity.PreloadEnabled, SecureSTSPreloadEnabled, s.StrictTransportSecurity.PreloadEnabled, "Enable HSTS preload")
	fs.StringVar(&s.XContentTypeOptions, SecureXContentTypeOptions, s.XContentTypeOptions, "X-Content-Type-Options header value")
	fs.StringVar(&s.XFrameOptions, SecureXFrameOptions, s.XFrameOptions, "X-Frame-Options header value")

	return fs
}

// NewSecure creates a new security middleware with the given configuration.
func NewSecure(config *Secure) echo.MiddlewareFunc {
	return secure.New(secure.Config{
		ContentSecurityPolicy:           config.ContentSecurityPolicy,
		ContentSecurityPolicyReportOnly: config.ContentSecurityPolicyReportOnly,
		CrossOriginEmbedderPolicy:       config.CrossOriginEmbedderPolicy,
		CrossOriginOpenerPolicy:         config.CrossOriginOpenerPolicy,
		CrossOriginResourcePolicy:       config.CrossOriginResourcePolicy,
		PermissionsPolicy:               config.PermissionsPolicy,
		ReferrerPolicy:                  config.ReferrerPolicy,
		Server:                          config.Server,
		StrictTransportSecurity: secure.StrictTransportSecurity{
			MaxAge:            config.StrictTransportSecurity.MaxAge,
			ExcludeSubdomains: config.StrictTransportSecurity.ExcludeSubdomains,
			PreloadEnabled:    config.StrictTransportSecurity.PreloadEnabled,
		},
		XContentTypeOptions: config.XContentTypeOptions,
		XFrameOptions:       config.XFrameOptions,
	})
}
