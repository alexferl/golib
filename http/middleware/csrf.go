package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
)

type CSRFSameSiteMode http.SameSite

const (
	csrfSameSiteDefaultMode = "default"
	csrfSameSiteLaxMode     = "lax"
	csrfSameSiteStrictMode  = "strict"
	csrfSameSiteNoneMode    = "none"
)

const (
	CSRFSameSiteDefaultMode CSRFSameSiteMode = iota + 1
	CSRFSameSiteLaxMode
	CSRFSameSiteStrictMode
	CSRFSameSiteNoneMode
)

var CSRFSameSiteModes = []string{csrfSameSiteDefaultMode, csrfSameSiteLaxMode, csrfSameSiteStrictMode, csrfSameSiteNoneMode}

func (m *CSRFSameSiteMode) String() string {
	switch *m {
	case CSRFSameSiteDefaultMode:
		return csrfSameSiteDefaultMode
	case CSRFSameSiteLaxMode:
		return csrfSameSiteLaxMode
	case CSRFSameSiteStrictMode:
		return csrfSameSiteStrictMode
	case CSRFSameSiteNoneMode:
		return csrfSameSiteNoneMode
	default:
		return fmt.Sprintf("unknown mode: %d", *m)
	}
}

func (m *CSRFSameSiteMode) Set(value string) error {
	switch strings.ToLower(value) {
	case csrfSameSiteDefaultMode:
		*m = CSRFSameSiteMode(http.SameSiteDefaultMode)
		return nil
	case csrfSameSiteLaxMode:
		*m = CSRFSameSiteMode(http.SameSiteLaxMode)
		return nil
	case csrfSameSiteStrictMode:
		*m = CSRFSameSiteMode(http.SameSiteStrictMode)
		return nil
	case csrfSameSiteNoneMode:
		*m = CSRFSameSiteMode(http.SameSiteNoneMode)
		return nil
	default:
		return fmt.Errorf("invalid same site mode: %s (must be one of: %s)", value, strings.Join(CSRFSameSiteModes, ", "))
	}
}

func (m *CSRFSameSiteMode) Type() string {
	return "string"
}

// CSRF holds configuration for Cross-Site Request Forgery protection.
type CSRF struct {
	// Enabled indicates whether CSRF middleware is enabled.
	// Optional. Default value false.
	Enabled bool

	// TokenLength specifies the length of the CSRF token.
	// Optional. Default value 32.
	TokenLength uint8

	// TokenLookup is a string in the form of "<source>:<name>" or "<source>:<name>,<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:X-CSRF-Token".
	// Possible values:
	// - "header:<name>" or "header:<name>:<cut-prefix>"
	// - "query:<name>"
	// - "form:<name>"
	// Multiple sources example:
	// - "header:X-CSRF-Token,query:csrf"
	TokenLookup string

	// ContextKey specifies the key used to store CSRF token in context.
	// Optional. Default value "csrf".
	ContextKey string

	// CookieName specifies the name of the CSRF cookie.
	// Optional. Default value "_csrf".
	CookieName string

	// CookieDomain specifies the domain for the CSRF cookie.
	// Optional. Default value "".
	CookieDomain string

	// CookiePath specifies the path for the CSRF cookie.
	// Optional. Default value "".
	CookiePath string

	// CookieMaxAge specifies the max age for the CSRF cookie in seconds.
	// Optional. Default value 86400.
	CookieMaxAge int

	// CookieSecure indicates whether the CSRF cookie should be secure.
	// Optional. Default value false.
	CookieSecure bool

	// CookieHTTPOnly indicates whether the CSRF cookie should be HTTP only.
	// Optional. Default value false.
	CookieHTTPOnly bool

	// CookieSameSite specifies the SameSite attribute for the CSRF cookie.
	// Optional. Default value SameSiteDefaultMode.
	CookieSameSite CSRFSameSiteMode
}

// DefaultCSRF provides default CSRF configuration.
var DefaultCSRF = &CSRF{
	Enabled:        false,
	TokenLength:    32,
	TokenLookup:    "header:X-CSRF-Token",
	ContextKey:     "csrf",
	CookieName:     "_csrf",
	CookieMaxAge:   86400,
	CookieSameSite: CSRFSameSiteMode(http.SameSiteDefaultMode),
}

const (
	CSRFEnabled        = "csrf-enabled"
	CSRFTokenLength    = "csrf-token-length"
	CSRFTokenLookup    = "csrf-token-lookup"
	CSRFContextKey     = "csrf-context-key"
	CSRFCookieName     = "csrf-cookie-name"
	CSRFCookieDomain   = "csrf-cookie-domain"
	CSRFCookiePath     = "csrf-cookie-path"
	CSRFCookieMaxAge   = "csrf-cookie-max-age"
	CSRFCookieSecure   = "csrf-cookie-secure"
	CSRFCookieHTTPOnly = "csrf-cookie-http-only"
	CSRFCookieSameSite = "csrf-cookie-same-site"
)

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (c *CSRF) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("CSRF", pflag.ExitOnError)

	fs.BoolVar(&c.Enabled, CSRFEnabled, c.Enabled, "Enable CSRF protection middleware")
	fs.Uint8Var(&c.TokenLength, CSRFTokenLength, c.TokenLength, "Length of the CSRF token")
	fs.StringVar(&c.TokenLookup, CSRFTokenLookup, c.TokenLookup, "Where to look for the CSRF token")
	fs.StringVar(&c.ContextKey, CSRFContextKey, c.ContextKey, "Key used to store CSRF token in context")
	fs.StringVar(&c.CookieName, CSRFCookieName, c.CookieName, "Name of the CSRF cookie")
	fs.StringVar(&c.CookieDomain, CSRFCookieDomain, c.CookieDomain, "Domain for the CSRF cookie")
	fs.StringVar(&c.CookiePath, CSRFCookiePath, c.CookiePath, "Path for the CSRF cookie")
	fs.IntVar(&c.CookieMaxAge, CSRFCookieMaxAge, c.CookieMaxAge, "Max age for the CSRF cookie in seconds")
	fs.BoolVar(&c.CookieSecure, CSRFCookieSecure, c.CookieSecure, "Whether the CSRF cookie should be secure")
	fs.BoolVar(&c.CookieHTTPOnly, CSRFCookieHTTPOnly, c.CookieHTTPOnly, "Whether the CSRF cookie should be HTTP only")
	fs.Var(&c.CookieSameSite, CSRFCookieSameSite, fmt.Sprintf("SameSite attribute for CSRF cookie\nValues: %s", strings.Join(CSRFSameSiteModes, ", ")))

	return fs
}

// NewCSRF creates a new CSRF middleware with the given configuration.
func NewCSRF(config *CSRF) echo.MiddlewareFunc {
	return middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLength:    config.TokenLength,
		TokenLookup:    config.TokenLookup,
		ContextKey:     config.ContextKey,
		CookieName:     config.CookieName,
		CookieDomain:   config.CookieDomain,
		CookiePath:     config.CookiePath,
		CookieMaxAge:   config.CookieMaxAge,
		CookieSecure:   config.CookieSecure,
		CookieHTTPOnly: config.CookieHTTPOnly,
		CookieSameSite: http.SameSite(config.CookieSameSite),
	})
}
