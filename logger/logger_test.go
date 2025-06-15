package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil config uses default",
			config:  nil,
			wantErr: false,
		},
		{
			name: "valid config",
			config: &Config{
				LogLevel:  "INFO",
				LogFormat: "json",
				LogOutput: "stdout",
			},
			wantErr: false,
		},
		{
			name: "empty fields use defaults",
			config: &Config{
				LogLevel:  "",
				LogFormat: "",
				LogOutput: "",
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			config: &Config{
				LogLevel:  "INVALID",
				LogFormat: "json",
				LogOutput: "stdout",
			},
			wantErr: true,
			errMsg:  "invalid log level",
		},
		{
			name: "invalid log format",
			config: &Config{
				LogLevel:  "INFO",
				LogFormat: "xml",
				LogOutput: "stdout",
			},
			wantErr: true,
			errMsg:  "invalid log format",
		},
		{
			name: "invalid log output",
			config: &Config{
				LogLevel:  "INFO",
				LogFormat: "json",
				LogOutput: "file",
			},
			wantErr: true,
			errMsg:  "invalid log output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Errorf("New() expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("New() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("New() unexpected error = %v", err)
					return
				}
				if logger == nil {
					t.Errorf("New() returned nil logger")
				}
			}
		})
	}
}

func TestCreateZerologLogger(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "text format stdout",
			config: &Config{
				LogLevel:  "INFO",
				LogFormat: "text",
				LogOutput: "stdout",
			},
			wantErr: false,
		},
		{
			name: "json format stderr",
			config: &Config{
				LogLevel:  "DEBUG",
				LogFormat: "json",
				LogOutput: "stderr",
			},
			wantErr: false,
		},
		{
			name: "unknown output",
			config: &Config{
				LogLevel:  "INFO",
				LogFormat: "json",
				LogOutput: "file",
			},
			wantErr: true,
		},
		{
			name: "unknown format",
			config: &Config{
				LogLevel:  "INFO",
				LogFormat: "xml",
				LogOutput: "stdout",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := createZerologLogger(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Errorf("createZerologLogger() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("createZerologLogger() unexpected error = %v", err)
				}
				// Test that logger works by creating a log event
				event := logger.Info()
				if event == nil {
					t.Errorf("createZerologLogger() returned non-functional logger")
				}
			}
		})
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected zerolog.Level
		wantErr  bool
	}{
		{"panic", LevelPanic, zerolog.PanicLevel, false},
		{"fatal", LevelFatal, zerolog.FatalLevel, false},
		{"error", LevelError, zerolog.ErrorLevel, false},
		{"warn", LevelWarn, zerolog.WarnLevel, false},
		{"info", LevelInfo, zerolog.InfoLevel, false},
		{"debug", LevelDebug, zerolog.DebugLevel, false},
		{"trace", LevelTrace, zerolog.TraceLevel, false},
		{"disabled", LevelDisabled, zerolog.Disabled, false},
		{"invalid", "INVALID", zerolog.InfoLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := parseLogLevel(tt.level)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseLogLevel() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("parseLogLevel() unexpected error = %v", err)
				}
				if level != tt.expected {
					t.Errorf("parseLogLevel() = %v, want %v", level, tt.expected)
				}
			}
		})
	}
}

func TestLogger_GetMethods(t *testing.T) {
	config := &Config{
		LogLevel:  "INFO",
		LogFormat: "json",
		LogOutput: "stdout",
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test GetConfig
	gotConfig := logger.GetConfig()
	if gotConfig != config {
		t.Errorf("GetConfig() = %v, want %v", gotConfig, config)
	}

	// Test GetLogger - verify it can be used
	zerologLogger := logger.GetLogger()
	event := zerologLogger.Info()
	if event == nil {
		t.Errorf("GetLogger() returned non-functional logger")
	}
}

func TestLogger_LogMethods(t *testing.T) {
	var buf bytes.Buffer

	// Create a logger that writes to our buffer
	zerologLogger := zerolog.New(&buf).With().Timestamp().Logger()
	logger := &Logger{
		logger: zerologLogger,
		config: DefaultConfig,
	}

	// Test each log method
	tests := []struct {
		name  string
		logFn func() *zerolog.Event
		level string
	}{
		{"panic", logger.Panic, "panic"},
		{"fatal", logger.Fatal, "fatal"},
		{"error", logger.Error, "error"},
		{"warn", logger.Warn, "warn"},
		{"info", logger.Info, "info"},
		{"debug", logger.Debug, "debug"},
		{"trace", logger.Trace, "trace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			// Don't actually send the message for panic/fatal as they would exit
			if tt.level == "panic" || tt.level == "fatal" {
				event := tt.logFn()
				if event == nil {
					t.Errorf("%s() returned nil event", tt.name)
				}
				return
			}

			tt.logFn().Msg("test message")

			if buf.Len() == 0 {
				t.Errorf("%s() did not write to buffer", tt.name)
				return
			}

			var logEntry map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
				t.Errorf("Failed to parse log output: %v", err)
				return
			}

			if logEntry["level"] != tt.level {
				t.Errorf("%s() level = %v, want %v", tt.name, logEntry["level"], tt.level)
			}

			if logEntry["message"] != "test message" {
				t.Errorf("%s() message = %v, want 'test message'", tt.name, logEntry["message"])
			}
		})
	}
}

func TestLogger_WithLevel(t *testing.T) {
	var buf bytes.Buffer
	zerologLogger := zerolog.New(&buf).With().Timestamp().Logger()
	logger := &Logger{
		logger: zerologLogger,
		config: DefaultConfig,
	}

	logger.WithLevel(zerolog.WarnLevel).Msg("test message")

	if buf.Len() == 0 {
		t.Error("WithLevel() did not write to buffer")
		return
	}

	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Errorf("Failed to parse log output: %v", err)
		return
	}

	if logEntry["level"] != "warn" {
		t.Errorf("WithLevel() level = %v, want warn", logEntry["level"])
	}
}

func TestLogger_With(t *testing.T) {
	var buf bytes.Buffer
	zerologLogger := zerolog.New(&buf).With().Timestamp().Logger()
	logger := &Logger{
		logger: zerologLogger,
		config: DefaultConfig,
	}

	contextLogger := logger.With().Str("component", "test").Logger()
	contextLogger.Info().Msg("test message")

	if buf.Len() == 0 {
		t.Error("With() context logger did not write to buffer")
		return
	}

	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Errorf("Failed to parse log output: %v", err)
		return
	}

	if logEntry["component"] != "test" {
		t.Errorf("With() component = %v, want test", logEntry["component"])
	}
}

func TestLogger_Log(t *testing.T) {
	var buf bytes.Buffer
	zerologLogger := zerolog.New(&buf).With().Timestamp().Logger()
	logger := &Logger{
		logger: zerologLogger,
		config: DefaultConfig,
	}

	event := logger.Log()
	if event == nil {
		t.Error("Log() returned nil event")
	}
}
