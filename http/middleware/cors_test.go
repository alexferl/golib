package middleware

import (
	"net/http"
	"reflect"
	"testing"
)

func TestCORS_FlagSet(t *testing.T) {
	config := &CORS{
		Enabled:          true,
		AllowOrigins:     []string{"https://example.com"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Custom-Header"},
		MaxAge:           3600,
	}

	fs := config.FlagSet()

	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	enabledFlag := fs.Lookup(CORSEnabled)
	if enabledFlag == nil {
		t.Errorf("Flag %s not found", CORSEnabled)
	} else {
		if enabledFlag.DefValue != "true" {
			t.Errorf("Flag %s default value = %v, want true", CORSEnabled, enabledFlag.DefValue)
		}
	}

	allowOriginsFlag := fs.Lookup(CORSAllowOrigins)
	if allowOriginsFlag == nil {
		t.Errorf("Flag %s not found", CORSAllowOrigins)
	} else {
		if allowOriginsFlag.DefValue != "[https://example.com]" {
			t.Errorf("Flag %s default value = %v, want [https://example.com]", CORSAllowOrigins, allowOriginsFlag.DefValue)
		}
	}

	maxAgeFlag := fs.Lookup(CORSMaxAge)
	if maxAgeFlag == nil {
		t.Errorf("Flag %s not found", CORSMaxAge)
	} else {
		if maxAgeFlag.DefValue != "3600" {
			t.Errorf("Flag %s default value = %v, want 3600", CORSMaxAge, maxAgeFlag.DefValue)
		}
	}
}

func TestCORS_FlagSet_Parse(t *testing.T) {
	config := &CORS{
		Enabled:          false,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet},
		AllowHeaders:     []string{},
		AllowCredentials: false,
		ExposeHeaders:    []string{},
		MaxAge:           0,
	}

	fs := config.FlagSet()

	args := []string{
		"--cors-enabled",
		"--cors-allow-origins", "https://example.com,https://example.org",
		"--cors-allow-methods", "GET,POST,DELETE",
		"--cors-allow-headers", "Authorization,Content-Type",
		"--cors-allow-credentials",
		"--cors-expose-headers", "X-Custom-Header",
		"--cors-max-age", "7200",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	if !config.Enabled {
		t.Errorf("Enabled = %v, want true", config.Enabled)
	}
	expectedOrigins := []string{"https://example.com", "https://example.org"}
	if !reflect.DeepEqual(config.AllowOrigins, expectedOrigins) {
		t.Errorf("AllowOrigins = %v, want %v", config.AllowOrigins, expectedOrigins)
	}
	expectedMethods := []string{"GET", "POST", "DELETE"}
	if !reflect.DeepEqual(config.AllowMethods, expectedMethods) {
		t.Errorf("AllowMethods = %v, want %v", config.AllowMethods, expectedMethods)
	}
	if !config.AllowCredentials {
		t.Errorf("AllowCredentials = %v, want true", config.AllowCredentials)
	}
	if config.MaxAge != 7200 {
		t.Errorf("MaxAge = %v, want 7200", config.MaxAge)
	}
}

func TestDefaultCORS(t *testing.T) {
	if DefaultCORS == nil {
		t.Fatal("DefaultCORS is nil")
	}

	if DefaultCORS.Enabled != false {
		t.Errorf("DefaultCORS.Enabled = %v, want false", DefaultCORS.Enabled)
	}
	expectedOrigins := []string{"*"}
	if !reflect.DeepEqual(DefaultCORS.AllowOrigins, expectedOrigins) {
		t.Errorf("DefaultCORS.AllowOrigins = %v, want %v", DefaultCORS.AllowOrigins, expectedOrigins)
	}
	expectedMethods := []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}
	if !reflect.DeepEqual(DefaultCORS.AllowMethods, expectedMethods) {
		t.Errorf("DefaultCORS.AllowMethods = %v, want %v", DefaultCORS.AllowMethods, expectedMethods)
	}
}

func TestCORS_FlagSet_DefaultValues(t *testing.T) {
	config := &CORS{
		Enabled:          true,
		AllowOrigins:     []string{"https://test.com"},
		AllowMethods:     []string{http.MethodGet},
		AllowHeaders:     []string{"Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Test"},
		MaxAge:           1800,
	}

	fs := config.FlagSet()

	enabledFlag := fs.Lookup(CORSEnabled)
	if enabledFlag == nil {
		t.Fatal("Enabled flag not found")
	}
	if enabledFlag.DefValue != "true" {
		t.Errorf("Enabled flag default = %v, want true", enabledFlag.DefValue)
	}

	maxAgeFlag := fs.Lookup(CORSMaxAge)
	if maxAgeFlag == nil {
		t.Fatal("MaxAge flag not found")
	}
	if maxAgeFlag.DefValue != "1800" {
		t.Errorf("MaxAge flag default = %v, want 1800", maxAgeFlag.DefValue)
	}
}

func TestCORS_FlagSet_DisabledByDefault(t *testing.T) {
	config := DefaultCORS

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

func TestNewCORS(t *testing.T) {
	config := &CORS{
		Enabled:      true,
		AllowOrigins: []string{"https://example.com"},
		AllowMethods: []string{"GET"},
	}

	middleware := NewCORS(config)
	if middleware == nil {
		t.Fatal("NewCORS() returned nil")
	}
}

func TestNewCORS_DefaultConfig(t *testing.T) {
	middleware := NewCORS(DefaultCORS)
	if middleware == nil {
		t.Fatal("NewCORS() with DefaultCORS returned nil")
	}
}
