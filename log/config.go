package log

import (
	"fmt"

	"github.com/spf13/pflag"
)

// Config holds all log configuration.
type Config struct {
	LogLevel  string
	LogOutput string
	LogWriter string
}

var DefaultConfig = &Config{
	LogLevel:  InfoLevel,
	LogOutput: "stdout",
	LogWriter: "console",
}

const (
	LogOutput = "log-output"
	LogWriter = "log-writer"
	LogLevel  = "log-level"
)

// BindFlags adds all the flags from the command line.
func (c *Config) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.LogOutput, LogOutput, c.LogOutput, "The output to write to. "+
		"'stdout' means log to stdout, 'stderr' means log to stderr.")
	fs.StringVar(&c.LogWriter, LogWriter, c.LogWriter,
		"The log writer. Valid writers are: 'console' and 'json'.")
	fs.StringVar(&c.LogLevel, LogLevel, c.LogLevel, fmt.Sprintf("The granularity of log outputs. "+
		"Valid levels: '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'",
		PanicLevel, FatalLevel, ErrorLevel, WarnLevel, InfoLevel, DebugLevel, TraceLevel, Disabled))
}
