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
		{&Config{LogLevel: WarnLevel, LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: InfoLevel, LogOutput: "stderr", LogWriter: "json"}, false},
		{&Config{LogLevel: InfoLevel, LogOutput: "stdout", LogWriter: "console"}, false},
		{&Config{LogLevel: PanicLevel, LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: FatalLevel, LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: ErrorLevel, LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: WarnLevel, LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: DebugLevel, LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: TraceLevel, LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: Disabled, LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: "wrong"}, true},
		{&Config{LogLevel: InfoLevel, LogOutput: "wrong"}, true},
		{&Config{LogLevel: InfoLevel, LogOutput: "stdout", LogWriter: "wrong"}, true},
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
