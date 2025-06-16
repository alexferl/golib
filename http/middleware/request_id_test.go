package middleware

import (
	"testing"
)

func TestRequestID_FlagSet(t *testing.T) {
	config := &RequestID{
		Enabled:      true,
		TargetHeader: "X-Custom-Request-ID",
	}

	fs := config.FlagSet()

	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	enabledFlag := fs.Lookup(RequestIDEnabled)
	if enabledFlag == nil {
		t.Errorf("Flag %s not found", RequestIDEnabled)
	} else {
		if enabledFlag.DefValue != "true" {
			t.Errorf("Flag %s default value = %v, want true", RequestIDEnabled, enabledFlag.DefValue)
		}
	}

	targetHeaderFlag := fs.Lookup(RequestIDTargetHeader)
	if targetHeaderFlag == nil {
		t.Errorf("Flag %s not found", RequestIDTargetHeader)
	} else {
		if targetHeaderFlag.DefValue != "X-Custom-Request-ID" {
			t.Errorf("Flag %s default value = %v, want X-Custom-Request-ID", RequestIDTargetHeader, targetHeaderFlag.DefValue)
		}
	}
}

func TestRequestID_FlagSet_Parse(t *testing.T) {
	config := &RequestID{
		Enabled:      false,
		TargetHeader: "X-Request-ID",
	}

	fs := config.FlagSet()

	args := []string{
		"--request-id-enabled",
		"--request-id-target-header", "X-Trace-ID",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	if !config.Enabled {
		t.Errorf("Enabled = %v, want true", config.Enabled)
	}
	if config.TargetHeader != "X-Trace-ID" {
		t.Errorf("TargetHeader = %v, want X-Trace-ID", config.TargetHeader)
	}
}

func TestDefaultRequestID(t *testing.T) {
	if DefaultRequestID == nil {
		t.Fatal("DefaultRequestID is nil")
	}

	if DefaultRequestID.Enabled != false {
		t.Errorf("DefaultRequestID.Enabled = %v, want false", DefaultRequestID.Enabled)
	}
	if DefaultRequestID.TargetHeader != "X-Request-ID" {
		t.Errorf("DefaultRequestID.TargetHeader = %v, want X-Request-ID", DefaultRequestID.TargetHeader)
	}
}

func TestRequestID_FlagSet_DefaultValues(t *testing.T) {
	config := &RequestID{
		Enabled:      true,
		TargetHeader: "X-Test-Request-ID",
	}

	fs := config.FlagSet()

	enabledFlag := fs.Lookup(RequestIDEnabled)
	if enabledFlag == nil {
		t.Fatal("Enabled flag not found")
	}
	if enabledFlag.DefValue != "true" {
		t.Errorf("Enabled flag default = %v, want true", enabledFlag.DefValue)
	}

	targetHeaderFlag := fs.Lookup(RequestIDTargetHeader)
	if targetHeaderFlag == nil {
		t.Fatal("TargetHeader flag not found")
	}
	if targetHeaderFlag.DefValue != "X-Test-Request-ID" {
		t.Errorf("TargetHeader flag default = %v, want X-Test-Request-ID", targetHeaderFlag.DefValue)
	}
}

func TestRequestID_FlagSet_DisabledByDefault(t *testing.T) {
	config := DefaultRequestID

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

func TestNewRequestID(t *testing.T) {
	config := &RequestID{
		Enabled:      true,
		TargetHeader: "X-Custom-ID",
	}

	middleware := NewRequestID(config)
	if middleware == nil {
		t.Fatal("NewRequestID() returned nil")
	}
}

func TestNewRequestID_DefaultConfig(t *testing.T) {
	middleware := NewRequestID(DefaultRequestID)
	if middleware == nil {
		t.Fatal("NewRequestID() with DefaultRequestID returned nil")
	}
}
