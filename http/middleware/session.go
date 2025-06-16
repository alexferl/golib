package middleware

import (
	"fmt"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/spf13/pflag"
)

// SessionCookieStore holds configuration for cookie-based session storage.
type SessionCookieStore struct {
	// Secret specifies the secret key for cookie sessions.
	// Optional. Default value "changeme".
	Secret string
}

type SessionStore string

const (
	sessionStoreCookie = "cookie"
)

const (
	SessionStoreCookie SessionStore = sessionStoreCookie
)

var SessionStores = []string{sessionStoreCookie}

func (s *SessionStore) String() string {
	switch *s {
	case SessionStoreCookie:
		return sessionStoreCookie
	default:
		return fmt.Sprintf("unknown store: %s", *s)
	}
}

func (s *SessionStore) Set(value string) error {
	switch strings.ToLower(value) {
	case sessionStoreCookie:
		*s = SessionStoreCookie
		return nil
	default:
		return fmt.Errorf("invalid session store: %s (must be one of: %s)", value, strings.Join(SessionStores, ", "))
	}
}

func (s *SessionStore) Type() string {
	return "string"
}

// Session holds configuration for session middleware.
type Session struct {
	// Enabled indicates whether session middleware is enabled.
	// Optional. Default value false.
	Enabled bool

	// Store specifies the session store type.
	// Optional. Default value "cookie".
	Store SessionStore

	// Cookie holds cookie store configuration.
	// Optional. Default value with secret "changeme".
	Cookie SessionCookieStore
}

// DefaultSession provides default Session configuration.
var DefaultSession = &Session{
	Enabled: false,
	Store:   sessionStoreCookie,
	Cookie: SessionCookieStore{
		Secret: "changeme",
	},
}

const (
	SessionEnabled      = "session-enabled"
	SessionStoreType    = "session-store" // Changed from SessionStore
	SessionCookieSecret = "session-cookie-secret"
)

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (s *Session) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("Session", pflag.ExitOnError)
	fs.BoolVar(&s.Enabled, SessionEnabled, s.Enabled, "Enable session middleware")
	fs.Var(&s.Store, SessionStoreType, fmt.Sprintf("Session store type\nValues: %s", strings.Join(SessionStores, ", ")))
	fs.StringVar(&s.Cookie.Secret, SessionCookieSecret, s.Cookie.Secret, "Secret key for cookie sessions")
	return fs
}

// NewSession creates a new session middleware with the given configuration.
func NewSession(config *Session) echo.MiddlewareFunc {
	switch config.Store {
	case SessionStoreCookie:
		return session.Middleware(sessions.NewCookieStore([]byte(config.Cookie.Secret)))
	}

	return nil
}
