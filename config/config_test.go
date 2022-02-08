package config

import (
	"testing"
)

func TestNew(t *testing.T) {
	c := New()

	tests := []struct {
		input    string
		expected string
	}{
		{DefaultConfig.AppName, c.AppName},
		{DefaultConfig.EnvName, c.EnvName},
		{DefaultConfig.EnvVarPrefix, c.EnvVarPrefix},
	}

	for _, tc := range tests {
		if tc.input != tc.expected {
			t.Errorf("got %s expected %s", tc.input, tc.expected)
		}
	}
}
