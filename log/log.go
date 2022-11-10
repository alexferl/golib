package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	PanicLevel = "PANIC"
	FatalLevel = "FATAL"
	ErrorLevel = "ERROR"
	WarnLevel  = "WARN"
	InfoLevel  = "INFO"
	DebugLevel = "DEBUG"
	TraceLevel = "TRACE"
	Disabled   = "DISABLED"
)

// New initializes the logger based on the passed Config,
// defaults to DefaultConfig if `config` is nil.
func New(config *Config) error {
	if config == nil {
		config = DefaultConfig
	}

	if config.LogLevel == "" {
		config.LogLevel = DefaultConfig.LogLevel
	}

	if config.LogOutput == "" {
		config.LogOutput = DefaultConfig.LogOutput
	}

	if config.LogWriter == "" {
		config.LogWriter = DefaultConfig.LogWriter
	}

	logOutput := strings.ToLower(config.LogOutput)
	logWriter := strings.ToLower(config.LogWriter)
	logLevel := strings.ToUpper(config.LogLevel)

	var f *os.File
	switch logOutput {
	case "stdout":
		f = os.Stdout
	case "stderr":
		f = os.Stderr
	default:
		return fmt.Errorf("unknown log-output '%s'", logOutput)
	}

	logger := zerolog.New(f)

	switch logWriter {
	case "console":
		logger = log.Output(zerolog.ConsoleWriter{
			Out: f, TimeFormat: time.RFC3339Nano,
		})
	case "json":
		break
	default:
		return fmt.Errorf("unknown log-writer '%s'", logWriter)
	}

	log.Logger = logger.With().Timestamp().Caller().Logger()

	switch logLevel {
	case PanicLevel:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case FatalLevel:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case ErrorLevel:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case WarnLevel:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case InfoLevel:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case DebugLevel:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case TraceLevel:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case Disabled:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	default:
		return fmt.Errorf("unknown log-level '%s'", logLevel)
	}

	return nil
}
