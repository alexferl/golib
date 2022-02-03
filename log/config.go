package log

import "github.com/spf13/pflag"

// Config holds all log configuration
type Config struct {
	LogLevel  string
	LogOutput string
	LogWriter string
}

var (
	DefaultConfig = &Config{
		LogLevel:  "info",
		LogOutput: "stdout",
		LogWriter: "console",
	}
)

// BindFlags adds all the flags from the command line
func (c *Config) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.LogOutput, "log-output", c.LogOutput, "The output to write to. "+
		"'stdout' means log to stdout, 'stderr' means log to stderr.")
	fs.StringVar(&c.LogWriter, "log-writer", c.LogWriter,
		"The log writer. Valid writers are: 'console' and 'json'.")
	fs.StringVar(&c.LogLevel, "log-level", c.LogLevel, "The granularity of log outputs. "+
		"Valid log levels: 'panic', 'fatal', 'error', 'warn', 'info', 'debug' and 'trace'.")
}
