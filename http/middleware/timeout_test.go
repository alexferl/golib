package middleware

import (
	"testing"
	"time"
)

func TestTimeout_FlagSet(t *testing.T) {
	config := &Timeout{
		Enabled:      false,
		ErrorMessage: "Custom timeout message",
		Duration:     30 * time.Second,
	}

	fs := config.FlagSet()

	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	enabledFlag := fs.Lookup(TimeoutEnabled)
	if enabledFlag == nil {
		t.Errorf("Flag %s not found", TimeoutEnabled)
	} else {
		if enabledFlag.DefValue != "false" {
			t.Errorf("Flag %s default value = %v, want false", TimeoutEnabled, enabledFlag.DefValue)
		}
	}

	errorMessageFlag := fs.Lookup(TimeoutErrorMessage)
	if errorMessageFlag == nil {
		t.Errorf("Flag %s not found", TimeoutErrorMessage)
	} else {
		if errorMessageFlag.DefValue != "Custom timeout message" {
			t.Errorf("Flag %s default value = %v, want Custom timeout message", TimeoutErrorMessage, errorMessageFlag.DefValue)
		}
	}

	durationFlag := fs.Lookup(TimeoutDuration)
	if durationFlag == nil {
		t.Errorf("Flag %s not found", TimeoutDuration)
	} else {
		if durationFlag.DefValue != "30s" {
			t.Errorf("Flag %s default value = %v, want 30s", TimeoutDuration, durationFlag.DefValue)
		}
	}
}

func TestTimeout_FlagSet_Parse(t *testing.T) {
	config := &Timeout{
		Enabled:      true,
		ErrorMessage: "Request timeout",
		Duration:     15 * time.Second,
	}

	fs := config.FlagSet()

	args := []string{
		"--timeout-enabled=false",
		"--timeout-error-message", "Service unavailable",
		"--timeout-duration", "45s",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	if config.Enabled {
		t.Errorf("Enabled = %v, want false", config.Enabled)
	}
	if config.ErrorMessage != "Service unavailable" {
		t.Errorf("ErrorMessage = %v, want Service unavailable", config.ErrorMessage)
	}
	if config.Duration != 45*time.Second {
		t.Errorf("Duration = %v, want 45s", config.Duration)
	}
}

func TestDefaultTimeout(t *testing.T) {
	if DefaultTimeout == nil {
		t.Fatal("DefaultTimeout is nil")
	}

	if DefaultTimeout.Enabled != true {
		t.Errorf("DefaultTimeout.Enabled = %v, want true", DefaultTimeout.Enabled)
	}
	if DefaultTimeout.ErrorMessage != "Request timeout" {
		t.Errorf("DefaultTimeout.ErrorMessage = %v, want Request timeout", DefaultTimeout.ErrorMessage)
	}
	if DefaultTimeout.Duration != 15*time.Second {
		t.Errorf("DefaultTimeout.Duration = %v, want 15s", DefaultTimeout.Duration)
	}
}

func TestTimeout_FlagSet_DefaultValues(t *testing.T) {
	config := &Timeout{
		Enabled:      false,
		ErrorMessage: "Test timeout",
		Duration:     5 * time.Minute,
	}

	fs := config.FlagSet()

	enabledFlag := fs.Lookup(TimeoutEnabled)
	if enabledFlag == nil {
		t.Fatal("Enabled flag not found")
	}
	if enabledFlag.DefValue != "false" {
		t.Errorf("Enabled flag default = %v, want false", enabledFlag.DefValue)
	}

	errorMessageFlag := fs.Lookup(TimeoutErrorMessage)
	if errorMessageFlag == nil {
		t.Fatal("ErrorMessage flag not found")
	}
	if errorMessageFlag.DefValue != "Test timeout" {
		t.Errorf("ErrorMessage flag default = %v, want Test timeout", errorMessageFlag.DefValue)
	}

	durationFlag := fs.Lookup(TimeoutDuration)
	if durationFlag == nil {
		t.Fatal("Duration flag not found")
	}
	if durationFlag.DefValue != "5m0s" {
		t.Errorf("Duration flag default = %v, want 5m0s", durationFlag.DefValue)
	}
}

func TestTimeout_FlagSet_EnabledByDefault(t *testing.T) {
	config := DefaultTimeout

	fs := config.FlagSet()

	var args []string

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse empty flags: %v", err)
	}

	if !config.Enabled {
		t.Errorf("Enabled = %v, want true (default)", config.Enabled)
	}
}

func TestNewTimeout(t *testing.T) {
	config := &Timeout{
		Enabled:      true,
		ErrorMessage: "Custom timeout",
		Duration:     30 * time.Second,
	}

	middleware := NewTimeout(config)
	if middleware == nil {
		t.Fatal("NewTimeout() returned nil")
	}
}

func TestNewTimeout_DefaultConfig(t *testing.T) {
	middleware := NewTimeout(DefaultTimeout)
	if middleware == nil {
		t.Fatal("NewTimeout() with DefaultTimeout returned nil")
	}
}
