package log

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestNew(t *testing.T) {
	tests := []struct {
		config *Config
		fail   bool
	}{
		{DefaultConfig, false},
		{&Config{LogLevel: "warn", LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: "info", LogOutput: "stderr", LogWriter: "json"}, false},
		{&Config{LogLevel: "info", LogOutput: "stdout", LogWriter: "console"}, false},
		{&Config{LogLevel: "panic", LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: "fatal", LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: "error", LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: "warn", LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: "debug", LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: "trace", LogOutput: "stdout", LogWriter: "json"}, false},
		{&Config{LogLevel: "wrong"}, true},
		{&Config{LogLevel: "info", LogOutput: "wrong"}, true},
		{&Config{LogLevel: "info", LogOutput: "stdout", LogWriter: "wrong"}, true},
	}

	for _, tc := range tests {
		err := New(tc.config)
		if !tc.fail {
			if err != nil {
				t.Errorf("%v", err)
			}

			if tc.config.LogLevel != zerolog.GlobalLevel().String() {
				t.Errorf("got %s expected %s", tc.config.LogLevel, zerolog.GlobalLevel().String())
			}
		} else {
			if err == nil {
				t.Error("test did not error")
			}
		}
	}
}
