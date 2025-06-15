package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	LevelPanic    = "PANIC"
	LevelFatal    = "FATAL"
	LevelError    = "ERROR"
	LevelWarn     = "WARN"
	LevelInfo     = "INFO"
	LevelDebug    = "DEBUG"
	LevelTrace    = "TRACE"
	LevelDisabled = "DISABLED"
)

const (
	FormatText = "text"
	FormatJSON = "json"
)

const (
	OutputStdOut = "stdout"
	OutputStdErr = "stderr"
)

var levels = []string{LevelPanic, LevelFatal, LevelError, LevelWarn, LevelInfo, LevelDebug, LevelTrace, LevelDisabled}
var formats = []string{FormatText, FormatJSON}
var outputs = []string{OutputStdOut, OutputStdErr}

// Logger wraps zerolog.Logger with configuration.
type Logger struct {
	logger zerolog.Logger
	config *Config
}

// New creates a new Logger instance with the given config.
// Uses DefaultConfig if config is nil.
func New(config *Config) (*Logger, error) {
	if config == nil {
		config = DefaultConfig
	}

	if config.LogLevel == "" {
		config.LogLevel = DefaultConfig.LogLevel
	}
	if config.LogOutput == "" {
		config.LogOutput = DefaultConfig.LogOutput
	}
	if config.LogFormat == "" {
		config.LogFormat = DefaultConfig.LogFormat
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	logger, err := createZerologLogger(config)
	if err != nil {
		return nil, err
	}

	return &Logger{
		logger: logger,
		config: config,
	}, nil
}

// createZerologLogger creates and configures a zerolog.Logger.
func createZerologLogger(config *Config) (zerolog.Logger, error) {
	logOutput := strings.ToLower(config.LogOutput)
	logFormat := strings.ToLower(config.LogFormat)
	logLevel := strings.ToUpper(config.LogLevel)

	var output io.Writer
	switch logOutput {
	case OutputStdOut:
		output = os.Stdout
	case OutputStdErr:
		output = os.Stderr
	default:
		return zerolog.Logger{}, fmt.Errorf("unknown log output '%s'", logOutput)
	}

	var logger zerolog.Logger
	switch logFormat {
	case FormatText:
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339Nano,
		})
	case FormatJSON:
		logger = zerolog.New(output)
	default:
		return zerolog.Logger{}, fmt.Errorf("unknown log format '%s'", logFormat)
	}

	logger = logger.With().Timestamp().Caller().Logger()

	level, err := parseLogLevel(logLevel)
	if err != nil {
		return zerolog.Logger{}, err
	}
	logger = logger.Level(level)

	return logger, nil
}

// parseLogLevel converts string log level to zerolog.Level.
func parseLogLevel(level string) (zerolog.Level, error) {
	switch level {
	case LevelPanic:
		return zerolog.PanicLevel, nil
	case LevelFatal:
		return zerolog.FatalLevel, nil
	case LevelError:
		return zerolog.ErrorLevel, nil
	case LevelWarn:
		return zerolog.WarnLevel, nil
	case LevelInfo:
		return zerolog.InfoLevel, nil
	case LevelDebug:
		return zerolog.DebugLevel, nil
	case LevelTrace:
		return zerolog.TraceLevel, nil
	case LevelDisabled:
		return zerolog.Disabled, nil
	default:
		return zerolog.InfoLevel, fmt.Errorf("unknown log level '%s'", level)
	}
}

// GetLogger returns the underlying zerolog.Logger.
func (l *Logger) GetLogger() zerolog.Logger {
	return l.logger
}

// GetConfig returns the logger configuration.
func (l *Logger) GetConfig() *Config {
	return l.config
}

// Panic creates a panic level log event.
func (l *Logger) Panic() *zerolog.Event {
	return l.logger.Panic()
}

// Fatal creates a fatal level log event.
func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

// Error creates an error level log event.
func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

// Warn creates a warning level log event.
func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

// Info creates an info level log event.
func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

// Debug creates a debug level log event.
func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

// Trace creates a trace level log event.
func (l *Logger) Trace() *zerolog.Event {
	return l.logger.Trace()
}

// Log creates a log event with no specific level.
func (l *Logger) Log() *zerolog.Event {
	return l.logger.Log()
}

// WithLevel creates a log event with the specified level.
func (l *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
	return l.logger.WithLevel(level)
}

// With creates a child logger with additional context.
func (l *Logger) With() zerolog.Context {
	return l.logger.With()
}
