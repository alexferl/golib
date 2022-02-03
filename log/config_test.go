package log

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestConfig_AddFlags(t *testing.T) {
	c := DefaultConfig
	fs := &pflag.FlagSet{}
	c.BindFlags(fs)

	level, _ := fs.GetString("log-level")
	out, _ := fs.GetString("log-output")
	writer, _ := fs.GetString("log-writer")

	var tests = []struct {
		input    string
		expected string
	}{
		{DefaultConfig.LogLevel, level},
		{DefaultConfig.LogOutput, out},
		{DefaultConfig.LogWriter, writer},
	}

	for _, tc := range tests {
		if tc.input != tc.expected {
			t.Errorf("got %s expected %s", tc.input, tc.expected)
		}
	}
}
