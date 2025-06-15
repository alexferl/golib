package logger

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/pflag"
)

// Config holds logger configuration options.
type Config struct {
	LogLevel  string
	LogFormat string
	LogOutput string
}

// DefaultConfig provides sensible default values.
var DefaultConfig = &Config{
	LogLevel:  LevelInfo,
	LogFormat: FormatText,
	LogOutput: OutputStdOut,
}

const (
	LogLevel  = "log-level"
	LogFormat = "log-format"
	LogOutput = "log-output"
)

// Validate checks if configuration values are valid.
func (c *Config) Validate() error {
	logLevel := strings.ToUpper(c.LogLevel)
	logFormat := strings.ToLower(c.LogFormat)
	logOutput := strings.ToLower(c.LogOutput)

	if !slices.Contains(levels, logLevel) {
		return fmt.Errorf("invalid log level '%s', must be one of: %s",
			c.LogLevel, strings.Join(levels, ", "))
	}

	if !slices.Contains(formats, logFormat) {
		return fmt.Errorf("invalid log format '%s', must be one of: %s",
			c.LogFormat, strings.Join(formats, ", "))
	}

	if !slices.Contains(outputs, logOutput) {
		return fmt.Errorf("invalid log output '%s', must be one of: %s",
			c.LogOutput, strings.Join(outputs, ", "))
	}

	return nil
}

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (c *Config) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("Logger", pflag.ExitOnError)
	fs.StringVar(&c.LogLevel, LogLevel, c.LogLevel,
		fmt.Sprintf("Log granularity\nValues: %s", strings.Join(levels, ", ")),
	)
	fs.StringVar(&c.LogFormat, LogFormat, c.LogFormat,
		fmt.Sprintf("Log format\nValues: %s", strings.Join(formats, ", ")),
	)
	fs.StringVar(&c.LogOutput, LogOutput, c.LogOutput,
		fmt.Sprintf("Output destination\nValues: %s", strings.Join(outputs, ", ")),
	)

	return fs
}
