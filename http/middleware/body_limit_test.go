package middleware

import (
	"testing"
)

func TestBodyLimit_FlagSet(t *testing.T) {
	config := &BodyLimit{
		Enabled: true,
		MaxSize: "2MB",
	}

	fs := config.FlagSet()

	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	enabledFlag := fs.Lookup(BodyLimitEnabled)
	if enabledFlag == nil {
		t.Errorf("Flag %s not found", BodyLimitEnabled)
	} else {
		if enabledFlag.DefValue != "true" {
			t.Errorf("Flag %s default value = %v, want true", BodyLimitEnabled, enabledFlag.DefValue)
		}
	}

	maxSizeFlag := fs.Lookup(BodyLimitMaxSize)
	if maxSizeFlag == nil {
		t.Errorf("Flag %s not found", BodyLimitMaxSize)
	} else {
		if maxSizeFlag.DefValue != "2MB" {
			t.Errorf("Flag %s default value = %v, want 2MB", BodyLimitMaxSize, maxSizeFlag.DefValue)
		}
	}
}

func TestBodyLimit_FlagSet_Parse(t *testing.T) {
	config := &BodyLimit{
		Enabled: false,
		MaxSize: "1MB",
	}

	fs := config.FlagSet()

	args := []string{
		"--body-limit-enabled",
		"--body-limit-max-size", "5MB",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	if !config.Enabled {
		t.Errorf("Enabled = %v, want true", config.Enabled)
	}
	if config.MaxSize != "5MB" {
		t.Errorf("MaxSize = %v, want 5MB", config.MaxSize)
	}
}

func TestDefaultBodyLimit(t *testing.T) {
	if DefaultBodyLimit == nil {
		t.Fatal("DefaultBodyLimit is nil")
	}

	if DefaultBodyLimit.Enabled != false {
		t.Errorf("DefaultBodyLimit.Enabled = %v, want false", DefaultBodyLimit.Enabled)
	}
	if DefaultBodyLimit.MaxSize != "1MB" {
		t.Errorf("DefaultBodyLimit.MaxSize = %v, want 1MB", DefaultBodyLimit.MaxSize)
	}
}

func TestBodyLimit_FlagSet_DefaultValues(t *testing.T) {
	config := &BodyLimit{
		Enabled: true,
		MaxSize: "10MB",
	}

	fs := config.FlagSet()

	enabledFlag := fs.Lookup(BodyLimitEnabled)
	if enabledFlag == nil {
		t.Fatal("Enabled flag not found")
	}
	if enabledFlag.DefValue != "true" {
		t.Errorf("Enabled flag default = %v, want true", enabledFlag.DefValue)
	}

	maxSizeFlag := fs.Lookup(BodyLimitMaxSize)
	if maxSizeFlag == nil {
		t.Fatal("MaxSize flag not found")
	}
	if maxSizeFlag.DefValue != "10MB" {
		t.Errorf("MaxSize flag default = %v, want 10MB", maxSizeFlag.DefValue)
	}
}

func TestBodyLimit_FlagSet_DisabledByDefault(t *testing.T) {
	config := DefaultBodyLimit

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

func TestNewBodyLimit(t *testing.T) {
	config := &BodyLimit{
		Enabled: true,
		MaxSize: "1KB",
	}

	middleware := NewBodyLimit(config)
	if middleware == nil {
		t.Fatal("NewBodyLimit() returned nil")
	}
}

func TestNewBodyLimit_DefaultConfig(t *testing.T) {
	middleware := NewBodyLimit(DefaultBodyLimit)
	if middleware == nil {
		t.Fatal("NewBodyLimit() with DefaultBodyLimit returned nil")
	}
}
