package middleware

import (
	secure "github.com/alexferl/echo-secure"
	"testing"
)

func TestSecure_FlagSet(t *testing.T) {
	config := &Secure{
		Enabled:                         true,
		ContentSecurityPolicy:           "default-src 'self'",
		ContentSecurityPolicyReportOnly: true,
		CrossOriginEmbedderPolicy:       "require-corp",
		CrossOriginOpenerPolicy:         "same-origin",
		CrossOriginResourcePolicy:       "cross-origin",
		PermissionsPolicy:               "geolocation=()",
		ReferrerPolicy:                  "strict-origin",
		Server:                          "MyServer/1.0",
		StrictTransportSecurity: StrictTransportSecurity{
			MaxAge:            31536000,
			ExcludeSubdomains: true,
			PreloadEnabled:    true,
		},
		XContentTypeOptions: "nosniff",
		XFrameOptions:       "DENY",
	}

	fs := config.FlagSet()

	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	enabledFlag := fs.Lookup(SecureEnabled)
	if enabledFlag == nil {
		t.Errorf("Flag %s not found", SecureEnabled)
	} else {
		if enabledFlag.DefValue != "true" {
			t.Errorf("Flag %s default value = %v, want true", SecureEnabled, enabledFlag.DefValue)
		}
	}

	cspFlag := fs.Lookup(SecureContentSecurityPolicy)
	if cspFlag == nil {
		t.Errorf("Flag %s not found", SecureContentSecurityPolicy)
	} else {
		if cspFlag.DefValue != "default-src 'self'" {
			t.Errorf("Flag %s default value = %v, want default-src 'self'", SecureContentSecurityPolicy, cspFlag.DefValue)
		}
	}

	stsMaxAgeFlag := fs.Lookup(SecureSTSMaxAge)
	if stsMaxAgeFlag == nil {
		t.Errorf("Flag %s not found", SecureSTSMaxAge)
	} else {
		if stsMaxAgeFlag.DefValue != "31536000" {
			t.Errorf("Flag %s default value = %v, want 31536000", SecureSTSMaxAge, stsMaxAgeFlag.DefValue)
		}
	}
}

func TestSecure_FlagSet_Parse(t *testing.T) {
	config := &Secure{
		Enabled:                         false,
		ContentSecurityPolicy:           "",
		ContentSecurityPolicyReportOnly: false,
		CrossOriginEmbedderPolicy:       "",
		CrossOriginOpenerPolicy:         "",
		CrossOriginResourcePolicy:       "",
		PermissionsPolicy:               "",
		ReferrerPolicy:                  "",
		Server:                          "",
		StrictTransportSecurity: StrictTransportSecurity{
			MaxAge:            0,
			ExcludeSubdomains: false,
			PreloadEnabled:    false,
		},
		XContentTypeOptions: "",
		XFrameOptions:       "",
	}

	fs := config.FlagSet()

	args := []string{
		"--secure-enabled",
		"--secure-content-security-policy", "default-src 'none'",
		"--secure-content-security-policy-report-only",
		"--secure-cross-origin-embedder-policy", "credentialless",
		"--secure-cross-origin-opener-policy", "same-origin-allow-popups",
		"--secure-cross-origin-resource-policy", "same-site",
		"--secure-permissions-policy", "camera=()",
		"--secure-referrer-policy", "no-referrer",
		"--secure-server", "TestServer/2.0",
		"--secure-sts-max-age", "63072000",
		"--secure-sts-exclude-subdomains",
		"--secure-sts-preload-enabled",
		"--secure-x-content-type-options", "nosniff",
		"--secure-x-frame-options", "SAMEORIGIN",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	if !config.Enabled {
		t.Errorf("Enabled = %v, want true", config.Enabled)
	}
	if config.ContentSecurityPolicy != "default-src 'none'" {
		t.Errorf("ContentSecurityPolicy = %v, want default-src 'none'", config.ContentSecurityPolicy)
	}
	if !config.ContentSecurityPolicyReportOnly {
		t.Errorf("ContentSecurityPolicyReportOnly = %v, want true", config.ContentSecurityPolicyReportOnly)
	}
	if config.CrossOriginEmbedderPolicy != "credentialless" {
		t.Errorf("CrossOriginEmbedderPolicy = %v, want credentialless", config.CrossOriginEmbedderPolicy)
	}
	if config.CrossOriginOpenerPolicy != "same-origin-allow-popups" {
		t.Errorf("CrossOriginOpenerPolicy = %v, want same-origin-allow-popups", config.CrossOriginOpenerPolicy)
	}
	if config.CrossOriginResourcePolicy != "same-site" {
		t.Errorf("CrossOriginResourcePolicy = %v, want same-site", config.CrossOriginResourcePolicy)
	}
	if config.PermissionsPolicy != "camera=()" {
		t.Errorf("PermissionsPolicy = %v, want camera=()", config.PermissionsPolicy)
	}
	if config.ReferrerPolicy != "no-referrer" {
		t.Errorf("ReferrerPolicy = %v, want no-referrer", config.ReferrerPolicy)
	}
	if config.Server != "TestServer/2.0" {
		t.Errorf("Server = %v, want TestServer/2.0", config.Server)
	}
	if config.StrictTransportSecurity.MaxAge != 63072000 {
		t.Errorf("StrictTransportSecurity.MaxAge = %v, want 63072000", config.StrictTransportSecurity.MaxAge)
	}
	if !config.StrictTransportSecurity.ExcludeSubdomains {
		t.Errorf("StrictTransportSecurity.ExcludeSubdomains = %v, want true", config.StrictTransportSecurity.ExcludeSubdomains)
	}
	if !config.StrictTransportSecurity.PreloadEnabled {
		t.Errorf("StrictTransportSecurity.PreloadEnabled = %v, want true", config.StrictTransportSecurity.PreloadEnabled)
	}
	if config.XContentTypeOptions != "nosniff" {
		t.Errorf("XContentTypeOptions = %v, want nosniff", config.XContentTypeOptions)
	}
	if config.XFrameOptions != "SAMEORIGIN" {
		t.Errorf("XFrameOptions = %v, want SAMEORIGIN", config.XFrameOptions)
	}
}

func TestDefaultSecure(t *testing.T) {
	if DefaultSecure == nil {
		t.Fatal("DefaultSecure is nil")
	}

	if DefaultSecure.Enabled != false {
		t.Errorf("DefaultSecure.Enabled = %v, want false", DefaultSecure.Enabled)
	}
	if DefaultSecure.ContentSecurityPolicy != secure.DefaultConfig.ContentSecurityPolicy {
		t.Errorf("DefaultSecure.ContentSecurityPolicy = %v, want %v", DefaultSecure.ContentSecurityPolicy, secure.DefaultConfig.ContentSecurityPolicy)
	}
	if DefaultSecure.ContentSecurityPolicyReportOnly != secure.DefaultConfig.ContentSecurityPolicyReportOnly {
		t.Errorf("DefaultSecure.ContentSecurityPolicyReportOnly = %v, want %v", DefaultSecure.ContentSecurityPolicyReportOnly, secure.DefaultConfig.ContentSecurityPolicyReportOnly)
	}
	if DefaultSecure.StrictTransportSecurity.MaxAge != secure.DefaultConfig.StrictTransportSecurity.MaxAge {
		t.Errorf("DefaultSecure.StrictTransportSecurity.MaxAge = %v, want %v", DefaultSecure.StrictTransportSecurity.MaxAge, secure.DefaultConfig.StrictTransportSecurity.MaxAge)
	}
}

func TestSecure_FlagSet_DefaultValues(t *testing.T) {
	config := &Secure{
		Enabled:               true,
		ContentSecurityPolicy: "test-csp",
		Server:                "TestServer",
		StrictTransportSecurity: StrictTransportSecurity{
			MaxAge:            1800,
			ExcludeSubdomains: true,
			PreloadEnabled:    false,
		},
		XFrameOptions: "DENY",
	}

	fs := config.FlagSet()

	enabledFlag := fs.Lookup(SecureEnabled)
	if enabledFlag == nil {
		t.Fatal("Enabled flag not found")
	}
	if enabledFlag.DefValue != "true" {
		t.Errorf("Enabled flag default = %v, want true", enabledFlag.DefValue)
	}

	cspFlag := fs.Lookup(SecureContentSecurityPolicy)
	if cspFlag == nil {
		t.Fatal("ContentSecurityPolicy flag not found")
	}
	if cspFlag.DefValue != "test-csp" {
		t.Errorf("ContentSecurityPolicy flag default = %v, want test-csp", cspFlag.DefValue)
	}

	serverFlag := fs.Lookup(SecureServer)
	if serverFlag == nil {
		t.Fatal("Server flag not found")
	}
	if serverFlag.DefValue != "TestServer" {
		t.Errorf("Server flag default = %v, want TestServer", serverFlag.DefValue)
	}
}

func TestSecure_FlagSet_DisabledByDefault(t *testing.T) {
	config := DefaultSecure

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

func TestNewSecure(t *testing.T) {
	config := &Secure{
		Enabled:                         true,
		ContentSecurityPolicy:           "default-src 'self'",
		ContentSecurityPolicyReportOnly: true,
		CrossOriginEmbedderPolicy:       "require-corp",
		CrossOriginOpenerPolicy:         "same-origin",
		CrossOriginResourcePolicy:       "cross-origin",
		PermissionsPolicy:               "geolocation=()",
		ReferrerPolicy:                  "strict-origin",
		Server:                          "MyServer/1.0",
		StrictTransportSecurity: StrictTransportSecurity{
			MaxAge:            31536000,
			ExcludeSubdomains: true,
			PreloadEnabled:    true,
		},
		XContentTypeOptions: "nosniff",
		XFrameOptions:       "DENY",
	}

	middleware := NewSecure(config)
	if middleware == nil {
		t.Fatal("NewSecure() returned nil")
	}
}

func TestNewSecure_DefaultConfig(t *testing.T) {
	middleware := NewSecure(DefaultSecure)
	if middleware == nil {
		t.Fatal("NewSecure() with DefaultSecure returned nil")
	}
}
