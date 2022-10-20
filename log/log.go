package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	logLevel := strings.ToLower(config.LogLevel)
	logOutput := strings.ToLower(config.LogOutput)
	logWriter := strings.ToLower(config.LogWriter)

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

	switch strings.ToLower(logWriter) {
	case "console":
		logger = log.Output(zerolog.ConsoleWriter{Out: f})
	case "json":
		break
	default:
		return fmt.Errorf("unknown log-writer '%s'", logWriter)
	}

	log.Logger = logger.With().Timestamp().Caller().Logger()

	switch strings.ToLower(logLevel) {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		return fmt.Errorf("unknown log-level '%s'", logLevel)
	}

	return nil
}
