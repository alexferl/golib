package logger

import (
	"strings"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
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
			name: "valid config case insensitive",
			config: &Config{
				LogLevel:  "debug",
				LogFormat: "TEXT",
				LogOutput: "STDERR",
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
			errMsg:  "invalid log level 'INVALID'",
		},
		{
			name: "invalid log format",
			config: &Config{
				LogLevel:  "INFO",
				LogFormat: "xml",
				LogOutput: "stdout",
			},
			wantErr: true,
			errMsg:  "invalid log format 'xml'",
		},
		{
			name: "invalid log output",
			config: &Config{
				LogLevel:  "INFO",
				LogFormat: "json",
				LogOutput: "file",
			},
			wantErr: true,
			errMsg:  "invalid log output 'file'",
		},
		{
			name: "empty values",
			config: &Config{
				LogLevel:  "",
				LogFormat: "",
				LogOutput: "",
			},
			wantErr: true,
			errMsg:  "invalid log level ''",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Config.Validate() expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Config.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Config.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestConfig_FlagSet(t *testing.T) {
	config := &Config{
		LogLevel:  "DEBUG",
		LogFormat: "text",
		LogOutput: "stderr",
	}

	fs := config.FlagSet()

	// Test that flagset is created
	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	// Test that flags are defined
	logLevelFlag := fs.Lookup(LogLevel)
	if logLevelFlag == nil {
		t.Errorf("Flag %s not found", LogLevel)
	} else {
		if logLevelFlag.DefValue != "DEBUG" {
			t.Errorf("Flag %s default value = %v, want DEBUG", LogLevel, logLevelFlag.DefValue)
		}
	}

	logFormatFlag := fs.Lookup(LogFormat)
	if logFormatFlag == nil {
		t.Errorf("Flag %s not found", LogFormat)
	} else {
		if logFormatFlag.DefValue != "text" {
			t.Errorf("Flag %s default value = %v, want text", LogFormat, logFormatFlag.DefValue)
		}
	}

	logOutputFlag := fs.Lookup(LogOutput)
	if logOutputFlag == nil {
		t.Errorf("Flag %s not found", LogOutput)
	} else {
		if logOutputFlag.DefValue != "stderr" {
			t.Errorf("Flag %s default value = %v, want stderr", LogOutput, logOutputFlag.DefValue)
		}
	}
}

func TestConfig_FlagSet_Parse(t *testing.T) {
	config := &Config{
		LogLevel:  "INFO",
		LogFormat: "json",
		LogOutput: "stdout",
	}

	fs := config.FlagSet()

	// Test parsing flags
	args := []string{
		"--log-level", "ERROR",
		"--log-format", "text",
		"--log-output", "stderr",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Check that config values were updated
	if config.LogLevel != "ERROR" {
		t.Errorf("LogLevel = %v, want ERROR", config.LogLevel)
	}
	if config.LogFormat != "text" {
		t.Errorf("LogFormat = %v, want text", config.LogFormat)
	}
	if config.LogOutput != "stderr" {
		t.Errorf("LogOutput = %v, want stderr", config.LogOutput)
	}
}

func TestDefaultConfig(t *testing.T) {
	if DefaultConfig == nil {
		t.Fatal("DefaultConfig is nil")
	}

	// Test default values
	if DefaultConfig.LogLevel != LevelInfo {
		t.Errorf("DefaultConfig.LogLevel = %v, want %v", DefaultConfig.LogLevel, LevelInfo)
	}
	if DefaultConfig.LogFormat != FormatText {
		t.Errorf("DefaultConfig.LogFormat = %v, want %v", DefaultConfig.LogFormat, FormatText)
	}
	if DefaultConfig.LogOutput != OutputStdOut {
		t.Errorf("DefaultConfig.LogOutput = %v, want %v", DefaultConfig.LogOutput, OutputStdOut)
	}

	// Test that default config is valid
	err := DefaultConfig.Validate()
	if err != nil {
		t.Errorf("DefaultConfig.Validate() error = %v", err)
	}
}

func TestConfig_FlagSet_Usage(t *testing.T) {
	config := &Config{}
	fs := config.FlagSet()

	// Test that usage strings contain expected values
	logLevelFlag := fs.Lookup(LogLevel)
	if logLevelFlag != nil {
		usage := logLevelFlag.Usage
		for _, level := range levels {
			if !strings.Contains(usage, level) {
				t.Errorf("Flag %s usage does not contain level %s", LogLevel, level)
			}
		}
	}

	logFormatFlag := fs.Lookup(LogFormat)
	if logFormatFlag != nil {
		usage := logFormatFlag.Usage
		for _, format := range formats {
			if !strings.Contains(usage, format) {
				t.Errorf("Flag %s usage does not contain format %s", LogFormat, format)
			}
		}
	}

	logOutputFlag := fs.Lookup(LogOutput)
	if logOutputFlag != nil {
		usage := logOutputFlag.Usage
		for _, output := range outputs {
			if !strings.Contains(usage, output) {
				t.Errorf("Flag %s usage does not contain output %s", LogOutput, output)
			}
		}
	}
}
