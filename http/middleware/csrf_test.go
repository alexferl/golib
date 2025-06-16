package middleware

import (
	"net/http"
	"testing"
)

func TestCSRF_FlagSet(t *testing.T) {
	config := &CSRF{
		Enabled:        true,
		TokenLength:    16,
		TokenLookup:    "form:_token",
		ContextKey:     "token",
		CookieName:     "_token",
		CookieDomain:   "example.com",
		CookiePath:     "/api",
		CookieMaxAge:   3600,
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: CSRFSameSiteMode(http.SameSiteStrictMode),
	}

	fs := config.FlagSet()

	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	enabledFlag := fs.Lookup(CSRFEnabled)
	if enabledFlag == nil {
		t.Errorf("Flag %s not found", CSRFEnabled)
	} else {
		if enabledFlag.DefValue != "true" {
			t.Errorf("Flag %s default value = %v, want true", CSRFEnabled, enabledFlag.DefValue)
		}
	}

	tokenLengthFlag := fs.Lookup(CSRFTokenLength)
	if tokenLengthFlag == nil {
		t.Errorf("Flag %s not found", CSRFTokenLength)
	} else {
		if tokenLengthFlag.DefValue != "16" {
			t.Errorf("Flag %s default value = %v, want 16", CSRFTokenLength, tokenLengthFlag.DefValue)
		}
	}

	cookieNameFlag := fs.Lookup(CSRFCookieName)
	if cookieNameFlag == nil {
		t.Errorf("Flag %s not found", CSRFCookieName)
	} else {
		if cookieNameFlag.DefValue != "_token" {
			t.Errorf("Flag %s default value = %v, want _token", CSRFCookieName, cookieNameFlag.DefValue)
		}
	}
}

func TestCSRF_FlagSet_Parse(t *testing.T) {
	config := &CSRF{
		Enabled:        false,
		TokenLength:    32,
		TokenLookup:    "header:X-CSRF-Token",
		ContextKey:     "csrf",
		CookieName:     "_csrf",
		CookieMaxAge:   86400,
		CookieSecure:   false,
		CookieHTTPOnly: false,
		CookieSameSite: CSRFSameSiteMode(http.SameSiteDefaultMode),
	}

	fs := config.FlagSet()

	args := []string{
		"--csrf-enabled",
		"--csrf-token-length", "64",
		"--csrf-token-lookup", "form:csrf_token",
		"--csrf-context-key", "token",
		"--csrf-cookie-name", "_token",
		"--csrf-cookie-domain", "example.com",
		"--csrf-cookie-path", "/secure",
		"--csrf-cookie-max-age", "7200",
		"--csrf-cookie-secure",
		"--csrf-cookie-http-only",
		"--csrf-cookie-same-site", "strict",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	if !config.Enabled {
		t.Errorf("Enabled = %v, want true", config.Enabled)
	}
	if config.TokenLength != 64 {
		t.Errorf("TokenLength = %v, want 64", config.TokenLength)
	}
	if config.TokenLookup != "form:csrf_token" {
		t.Errorf("TokenLookup = %v, want form:csrf_token", config.TokenLookup)
	}
	if config.ContextKey != "token" {
		t.Errorf("ContextKey = %v, want token", config.ContextKey)
	}
	if config.CookieName != "_token" {
		t.Errorf("CookieName = %v, want _token", config.CookieName)
	}
	if config.CookieDomain != "example.com" {
		t.Errorf("CookieDomain = %v, want example.com", config.CookieDomain)
	}
	if config.CookiePath != "/secure" {
		t.Errorf("CookiePath = %v, want /secure", config.CookiePath)
	}
	if config.CookieMaxAge != 7200 {
		t.Errorf("CookieMaxAge = %v, want 7200", config.CookieMaxAge)
	}
	if !config.CookieSecure {
		t.Errorf("CookieSecure = %v, want true", config.CookieSecure)
	}
	if !config.CookieHTTPOnly {
		t.Errorf("CookieHTTPOnly = %v, want true", config.CookieHTTPOnly)
	}
	if config.CookieSameSite != CSRFSameSiteMode(http.SameSiteStrictMode) {
		t.Errorf("CookieSameSite = %v, want %v", config.CookieSameSite, CSRFSameSiteMode(http.SameSiteStrictMode))
	}
}

func TestDefaultCSRF(t *testing.T) {
	if DefaultCSRF == nil {
		t.Fatal("DefaultCSRF is nil")
	}

	if DefaultCSRF.Enabled != false {
		t.Errorf("DefaultCSRF.Enabled = %v, want false", DefaultCSRF.Enabled)
	}
	if DefaultCSRF.TokenLength != 32 {
		t.Errorf("DefaultCSRF.TokenLength = %v, want 32", DefaultCSRF.TokenLength)
	}
	if DefaultCSRF.TokenLookup != "header:X-CSRF-Token" {
		t.Errorf("DefaultCSRF.TokenLookup = %v, want header:X-CSRF-Token", DefaultCSRF.TokenLookup)
	}
	if DefaultCSRF.ContextKey != "csrf" {
		t.Errorf("DefaultCSRF.ContextKey = %v, want csrf", DefaultCSRF.ContextKey)
	}
	if DefaultCSRF.CookieName != "_csrf" {
		t.Errorf("DefaultCSRF.CookieName = %v, want _csrf", DefaultCSRF.CookieName)
	}
	if DefaultCSRF.CookieMaxAge != 86400 {
		t.Errorf("DefaultCSRF.CookieMaxAge = %v, want 86400", DefaultCSRF.CookieMaxAge)
	}
	if DefaultCSRF.CookieSameSite != CSRFSameSiteMode(http.SameSiteDefaultMode) {
		t.Errorf("DefaultCSRF.CookieSameSite = %v, want %v", DefaultCSRF.CookieSameSite, CSRFSameSiteMode(http.SameSiteDefaultMode))
	}
}

func TestCSRF_FlagSet_DefaultValues(t *testing.T) {
	config := &CSRF{
		Enabled:        true,
		TokenLength:    16,
		CookieName:     "_test",
		CookieMaxAge:   1800,
		CookieSecure:   true,
		CookieHTTPOnly: false,
	}

	fs := config.FlagSet()

	enabledFlag := fs.Lookup(CSRFEnabled)
	if enabledFlag == nil {
		t.Fatal("Enabled flag not found")
	}
	if enabledFlag.DefValue != "true" {
		t.Errorf("Enabled flag default = %v, want true", enabledFlag.DefValue)
	}

	tokenLengthFlag := fs.Lookup(CSRFTokenLength)
	if tokenLengthFlag == nil {
		t.Fatal("TokenLength flag not found")
	}
	if tokenLengthFlag.DefValue != "16" {
		t.Errorf("TokenLength flag default = %v, want 16", tokenLengthFlag.DefValue)
	}

	cookieNameFlag := fs.Lookup(CSRFCookieName)
	if cookieNameFlag == nil {
		t.Fatal("CookieName flag not found")
	}
	if cookieNameFlag.DefValue != "_test" {
		t.Errorf("CookieName flag default = %v, want _test", cookieNameFlag.DefValue)
	}
}

func TestCSRF_FlagSet_DisabledByDefault(t *testing.T) {
	config := DefaultCSRF

	fs := config.FlagSet()

	var args []string

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse empty flags: %v", err)
	}

	if config.Enabled {
		t.Errorf("Enabled = %v, want false (default)", config.Enabled)
	}
}

func TestCSRFSameSiteMode_String(t *testing.T) {
	tests := []struct {
		mode CSRFSameSiteMode
		want string
	}{
		{CSRFSameSiteDefaultMode, "default"},
		{CSRFSameSiteLaxMode, "lax"},
		{CSRFSameSiteStrictMode, "strict"},
		{CSRFSameSiteNoneMode, "none"},
		{CSRFSameSiteMode(99), "unknown mode: 99"},
	}

	for _, tt := range tests {
		got := tt.mode.String()
		if got != tt.want {
			t.Errorf("CSRFSameSiteMode(%d).String() = %v, want %v", tt.mode, got, tt.want)
		}
	}
}

func TestCSRFSameSiteMode_Set(t *testing.T) {
	tests := []struct {
		value   string
		want    CSRFSameSiteMode
		wantErr bool
	}{
		{"default", CSRFSameSiteMode(http.SameSiteDefaultMode), false},
		{"lax", CSRFSameSiteMode(http.SameSiteLaxMode), false},
		{"strict", CSRFSameSiteMode(http.SameSiteStrictMode), false},
		{"none", CSRFSameSiteMode(http.SameSiteNoneMode), false},
		{"DEFAULT", CSRFSameSiteMode(http.SameSiteDefaultMode), false},
		{"LAX", CSRFSameSiteMode(http.SameSiteLaxMode), false},
		{"invalid", CSRFSameSiteMode(0), true},
	}

	for _, tt := range tests {
		var mode CSRFSameSiteMode
		err := mode.Set(tt.value)
		if tt.wantErr {
			if err == nil {
				t.Errorf("CSRFSameSiteMode.Set(%q) expected error but got nil", tt.value)
			}
		} else {
			if err != nil {
				t.Errorf("CSRFSameSiteMode.Set(%q) unexpected error: %v", tt.value, err)
			}
			if mode != tt.want {
				t.Errorf("CSRFSameSiteMode.Set(%q) = %v, want %v", tt.value, mode, tt.want)
			}
		}
	}
}

func TestCSRFSameSiteMode_Type(t *testing.T) {
	var mode CSRFSameSiteMode
	if got := mode.Type(); got != "string" {
		t.Errorf("CSRFSameSiteMode.Type() = %v, want string", got)
	}
}

func TestNewCSRF(t *testing.T) {
	config := &CSRF{
		Enabled:     true,
		TokenLength: 32,
		CookieName:  "_csrf",
	}

	middleware := NewCSRF(config)
	if middleware == nil {
		t.Fatal("NewCSRF() returned nil")
	}
}

func TestNewCSRF_DefaultConfig(t *testing.T) {
	middleware := NewCSRF(DefaultCSRF)
	if middleware == nil {
		t.Fatal("NewCSRF() with DefaultCSRF returned nil")
	}
}
