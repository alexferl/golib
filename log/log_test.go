package log

import (
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestNew(t *testing.T) {
	tests := []struct {
		config *Config
		fail   bool
	}{
		{DefaultConfig, false},
		{&Config{LogLevel: WarnLevel, LogOutput: OutputStdOut, LogWriter: WriterJSON}, false},
		{&Config{LogLevel: InfoLevel, LogOutput: OutputStdErr, LogWriter: WriterJSON}, false},
		{&Config{LogLevel: InfoLevel, LogOutput: OutputStdOut, LogWriter: WriterConsole}, false},
		{&Config{LogLevel: PanicLevel, LogOutput: OutputStdOut, LogWriter: WriterJSON}, false},
		{&Config{LogLevel: FatalLevel, LogOutput: OutputStdOut, LogWriter: WriterJSON}, false},
		{&Config{LogLevel: ErrorLevel, LogOutput: OutputStdOut, LogWriter: WriterJSON}, false},
		{&Config{LogLevel: WarnLevel, LogOutput: OutputStdOut, LogWriter: WriterJSON}, false},
		{&Config{LogLevel: DebugLevel, LogOutput: OutputStdOut, LogWriter: WriterJSON}, false},
		{&Config{LogLevel: TraceLevel, LogOutput: OutputStdOut, LogWriter: WriterJSON}, false},
		{&Config{LogLevel: Disabled, LogOutput: OutputStdOut, LogWriter: WriterJSON}, false},
		{&Config{LogLevel: "wrong"}, true},
		{&Config{LogLevel: InfoLevel, LogOutput: "wrong"}, true},
		{&Config{LogLevel: InfoLevel, LogOutput: OutputStdOut, LogWriter: "wrong"}, true},
	}

	for _, tc := range tests {
		err := New(tc.config)
		if !tc.fail {
			if err != nil {
				t.Errorf("%v", err)
			}

			level := strings.ToUpper(zerolog.GlobalLevel().String())
			if tc.config.LogLevel != level {
				t.Errorf("got %s expected %s", tc.config.LogLevel, level)
			}
		} else {
			if err == nil {
				t.Error("test did not error")
			}
		}
	}
}
