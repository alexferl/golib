package config

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func TestNewConfigLoader(t *testing.T) {
	tests := []struct {
		name     string
		options  []Option
		expected *ConfigLoader
	}{
		{
			name:    "default config loader",
			options: nil,
			expected: &ConfigLoader{
				envVarPrefix: "",
				configPaths:  []string{"./configs", "/configs"},
				configType:   "toml",
			},
		},
		{
			name: "with env prefix",
			options: []Option{
				WithEnvPrefix("MYAPP"),
			},
			expected: &ConfigLoader{
				envVarPrefix: "MYAPP",
				configPaths:  []string{"./configs", "/configs"},
				configType:   "toml",
			},
		},
		{
			name: "with custom config paths",
			options: []Option{
				WithConfigPaths("./custom", "/etc/app"),
			},
			expected: &ConfigLoader{
				envVarPrefix: "",
				configPaths:  []string{"./custom", "/etc/app"},
				configType:   "toml",
			},
		},
		{
			name: "with config type",
			options: []Option{
				WithConfigType("yaml"),
			},
			expected: &ConfigLoader{
				envVarPrefix: "",
				configPaths:  []string{"./configs", "/configs"},
				configType:   "yaml",
			},
		},
		{
			name: "with all options",
			options: []Option{
				WithEnvPrefix("TEST"),
				WithConfigPaths("./test"),
				WithConfigType("json"),
			},
			expected: &ConfigLoader{
				envVarPrefix: "TEST",
				configPaths:  []string{"./test"},
				configType:   "json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewConfigLoader(tt.options...)

			if loader.envVarPrefix != tt.expected.envVarPrefix {
				t.Errorf("envVarPrefix = %v, want %v", loader.envVarPrefix, tt.expected.envVarPrefix)
			}
			if len(loader.configPaths) != len(tt.expected.configPaths) {
				t.Errorf("configPaths length = %v, want %v", len(loader.configPaths), len(tt.expected.configPaths))
			}
			for i, path := range loader.configPaths {
				if path != tt.expected.configPaths[i] {
					t.Errorf("configPaths[%d] = %v, want %v", i, path, tt.expected.configPaths[i])
				}
			}
			if loader.configType != tt.expected.configType {
				t.Errorf("configType = %v, want %v", loader.configType, tt.expected.configType)
			}
		})
	}
}

func TestConfigLoader_LoadConfig(t *testing.T) {
	// Reset viper for each test
	defer func() {
		viper.Reset()
	}()

	// Save original args and restore after test
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	tests := []struct {
		name       string
		args       []string
		loadOpts   []LoadOption
		flagSetup  func(fs *pflag.FlagSet)
		envVars    map[string]string
		wantErr    bool
		errContain string
		validate   func(*Config) error
	}{
		{
			name:     "load without config file - default values",
			args:     []string{"testapp"},
			loadOpts: nil,
			wantErr:  false,
			validate: func(c *Config) error {
				if c.AppName != DefaultAppName {
					t.Errorf("AppName = %v, want %v", c.AppName, DefaultAppName)
				}
				if c.EnvName != DefaultEnvName {
					t.Errorf("EnvName = %v, want %v", c.EnvName, DefaultEnvName)
				}
				return nil
			},
		},
		{
			name:     "load with command line flags",
			args:     []string{"testapp", "--app-name", "myapp", "--env-name", "prod"},
			loadOpts: nil,
			wantErr:  false,
			validate: func(c *Config) error {
				if c.AppName != "myapp" {
					t.Errorf("AppName = %v, want myapp", c.AppName)
				}
				if c.EnvName != "prod" {
					t.Errorf("EnvName = %v, want prod", c.EnvName)
				}
				return nil
			},
		},
		{
			name:     "load with environment variables",
			args:     []string{"testapp"},
			loadOpts: nil,
			envVars: map[string]string{
				"TEST_APP_NAME": "envapp",
				"TEST_ENV_NAME": "staging",
			},
			wantErr: false,
			validate: func(c *Config) error {
				if c.AppName != "envapp" {
					t.Errorf("AppName = %v, want envapp", c.AppName)
				}
				if c.EnvName != "staging" {
					t.Errorf("EnvName = %v, want staging", c.EnvName)
				}
				return nil
			},
		},
		{
			name:     "load with custom flag",
			args:     []string{"testapp", "--custom-flag", "value"},
			loadOpts: nil,
			flagSetup: func(fs *pflag.FlagSet) {
				fs.String("custom-flag", "", "custom flag")
			},
			wantErr: false,
		},
		{
			name:       "load with config file enabled",
			args:       []string{"testapp"},
			loadOpts:   []LoadOption{WithConfigFile(true)},
			wantErr:    true,
			errContain: "Not Found", // Updated to match actual Viper error
		},
		{
			name:       "load with custom config file",
			args:       []string{"testapp"},
			loadOpts:   []LoadOption{WithCustomConfigFile("custom-config")},
			wantErr:    true,
			errContain: "Not Found", // Updated to match actual Viper error
		},
		{
			name:       "validation error - empty app name",
			args:       []string{"testapp", "--app-name", ""},
			loadOpts:   nil,
			wantErr:    true,
			errContain: "application name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper for each test
			viper.Reset()

			for key, value := range tt.envVars {
				err := os.Setenv(key, value)
				if err != nil {
					t.Fatalf("failed to set env var %s: %v", key, err)
				}
				defer func(k string) {
					err := os.Unsetenv(k)
					if err != nil {
						t.Errorf("failed to unset env var %s: %v", k, err)
					}
				}(key)
			}

			// Set up command line args
			os.Args = tt.args

			loader := NewConfigLoader(WithEnvPrefix("TEST"))

			var flagSets []func(fs *pflag.FlagSet)
			if tt.flagSetup != nil {
				flagSets = append(flagSets, tt.flagSetup)
			}

			config, err := loader.LoadConfig(tt.loadOpts, flagSets...)

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContain != "" && !strings.Contains(err.Error(), tt.errContain) {
				t.Errorf("LoadConfig() error = %v, want to contain %v", err, tt.errContain)
				return
			}

			if err == nil && tt.validate != nil {
				if err := tt.validate(config); err != nil {
					t.Errorf("Validation failed: %v", err)
				}
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &Config{
				AppName: "myapp",
				EnvName: "prod",
			},
			wantErr: false,
		},
		{
			name: "empty app name",
			config: &Config{
				AppName: "",
				EnvName: "prod",
			},
			wantErr: true,
			errMsg:  "application name cannot be empty",
		},
		{
			name: "whitespace app name",
			config: &Config{
				AppName: "   ",
				EnvName: "prod",
			},
			wantErr: true,
			errMsg:  "application name cannot be empty",
		},
		{
			name: "empty env name",
			config: &Config{
				AppName: "myapp",
				EnvName: "",
			},
			wantErr: true,
			errMsg:  "environment name cannot be empty",
		},
		{
			name: "whitespace env name",
			config: &Config{
				AppName: "myapp",
				EnvName: "   ",
			},
			wantErr: true,
			errMsg:  "environment name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Validate() error = %v, want to contain %v", err, tt.errMsg)
			}
		})
	}
}

func TestNormalizeFlags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "flag with underscore",
			input:    "test_flag",
			expected: "test-flag",
		},
		{
			name:     "flag with multiple underscores",
			input:    "test_flag_name",
			expected: "test-flag-name",
		},
		{
			name:     "flag without underscore",
			input:    "testflag",
			expected: "testflag",
		},
		{
			name:     "flag with dash",
			input:    "test-flag",
			expected: "test-flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			normalized := normalizeFlags(fs, tt.input)
			if string(normalized) != tt.expected {
				t.Errorf("normalizeFlags() = %v, want %v", normalized, tt.expected)
			}
		})
	}
}

func TestConfigLoader_AddConfigPath(t *testing.T) {
	loader := NewConfigLoader()
	initialLen := len(loader.configPaths)

	loader.AddConfigPath("/new/path")

	if len(loader.configPaths) != initialLen+1 {
		t.Errorf("AddConfigPath() did not add path, got %d paths, want %d",
			len(loader.configPaths), initialLen+1)
	}

	if loader.configPaths[len(loader.configPaths)-1] != "/new/path" {
		t.Errorf("AddConfigPath() added wrong path, got %v, want /new/path",
			loader.configPaths[len(loader.configPaths)-1])
	}

	// Test chaining
	loader.AddConfigPath("/another/path").AddConfigPath("/third/path")

	if len(loader.configPaths) != initialLen+3 {
		t.Errorf("Chained AddConfigPath() failed, got %d paths, want %d",
			len(loader.configPaths), initialLen+3)
	}
}

func TestLoadOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []LoadOption
		expected loadOptions
	}{
		{
			name:    "no options",
			options: nil,
			expected: loadOptions{
				loadConfigFile: false,
				configFileName: "",
			},
		},
		{
			name:    "with config file enabled",
			options: []LoadOption{WithConfigFile(true)},
			expected: loadOptions{
				loadConfigFile: true,
				configFileName: "",
			},
		},
		{
			name:    "with config file disabled",
			options: []LoadOption{WithConfigFile(false)},
			expected: loadOptions{
				loadConfigFile: false,
				configFileName: "",
			},
		},
		{
			name:    "with custom config file",
			options: []LoadOption{WithCustomConfigFile("custom")},
			expected: loadOptions{
				loadConfigFile: true,
				configFileName: "custom",
			},
		},
		{
			name: "multiple options",
			options: []LoadOption{
				WithConfigFile(false),
				WithCustomConfigFile("override"),
			},
			expected: loadOptions{
				loadConfigFile: true, // WithCustomConfigFile overrides
				configFileName: "override",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &loadOptions{}
			for _, option := range tt.options {
				option(opts)
			}

			if opts.loadConfigFile != tt.expected.loadConfigFile {
				t.Errorf("loadConfigFile = %v, want %v", opts.loadConfigFile, tt.expected.loadConfigFile)
			}
			if opts.configFileName != tt.expected.configFileName {
				t.Errorf("configFileName = %v, want %v", opts.configFileName, tt.expected.configFileName)
			}
		})
	}
}
