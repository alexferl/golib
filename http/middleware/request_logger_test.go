package middleware

import (
	"testing"

	"github.com/alexferl/golib/logger"
)

func TestRequestLogger_FlagSet(t *testing.T) {
	config := &RequestLogger{
		Enabled: true,
		Logger:  nil,
	}

	fs := config.FlagSet()

	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	enabledFlag := fs.Lookup(RequestLoggerEnabled)
	if enabledFlag == nil {
		t.Errorf("Flag %s not found", RequestLoggerEnabled)
	} else {
		if enabledFlag.DefValue != "true" {
			t.Errorf("Flag %s default value = %v, want true", RequestLoggerEnabled, enabledFlag.DefValue)
		}
	}
}

func TestRequestLogger_FlagSet_Parse(t *testing.T) {
	config := &RequestLogger{
		Enabled: false,
		Logger:  nil,
	}

	fs := config.FlagSet()

	args := []string{
		"--request-logger-enabled",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	if !config.Enabled {
		t.Errorf("Enabled = %v, want true", config.Enabled)
	}
}

func TestDefaultLogger(t *testing.T) {
	if DefaultLogger == nil {
		t.Fatal("DefaultLogger is nil")
	}

	if DefaultLogger.Enabled != false {
		t.Errorf("DefaultLogger.Enabled = %v, want false", DefaultLogger.Enabled)
	}
	if DefaultLogger.Logger != nil {
		t.Errorf("DefaultLogger.Logger = %v, want nil", DefaultLogger.Logger)
	}
}

func TestRequestLogger_FlagSet_DefaultValues(t *testing.T) {
	config := &RequestLogger{
		Enabled: true,
		Logger:  nil,
	}

	fs := config.FlagSet()

	enabledFlag := fs.Lookup(RequestLoggerEnabled)
	if enabledFlag == nil {
		t.Fatal("Enabled flag not found")
	}
	if enabledFlag.DefValue != "true" {
		t.Errorf("Enabled flag default = %v, want true", enabledFlag.DefValue)
	}
}

func TestRequestLogger_FlagSet_DisabledByDefault(t *testing.T) {
	config := DefaultLogger

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

func TestNewRequestLogger(t *testing.T) {
	config := &RequestLogger{
		Enabled: true,
		Logger:  nil,
	}

	middleware := NewRequestLogger(config)
	if middleware == nil {
		t.Fatal("NewRequestLogger() returned nil")
	}
}

func TestNewRequestLogger_WithCustomLogger(t *testing.T) {
	customLogger, err := logger.New(&logger.Config{
		LogLevel:  logger.LevelInfo,
		LogFormat: logger.FormatJSON,
		LogOutput: logger.OutputStdOut,
	})
	if err != nil {
		t.Fatalf("Failed to create custom logger: %v", err)
	}

	config := &RequestLogger{
		Enabled: true,
		Logger:  customLogger,
	}

	middleware := NewRequestLogger(config)
	if middleware == nil {
		t.Fatal("NewRequestLogger() with custom logger returned nil")
	}
}

func TestNewRequestLogger_DefaultConfig(t *testing.T) {
	middleware := NewRequestLogger(DefaultLogger)
	if middleware == nil {
		t.Fatal("NewRequestLogger() with DefaultLogger returned nil")
	}
}
