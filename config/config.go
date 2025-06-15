package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	DefaultAppName = "app"
	DefaultEnvName = "local"
	AppName        = "app-name"
	EnvName        = "env-name"
)

// Config holds all global configuration for the application.
type Config struct {
	AppName string
	EnvName string
}

// ConfigLoader manages the configuration loading process.
type ConfigLoader struct {
	envVarPrefix string
	configPaths  []string
	configType   string
}

// Option defines a function type for configuring ConfigLoader.
type Option func(*ConfigLoader)

// LoadOption defines a function type for configuring the load process.
type LoadOption func(*loadOptions)

// loadOptions holds configuration for the loading process.
type loadOptions struct {
	loadConfigFile bool
	configFileName string
}

// WithEnvPrefix sets the environment variable prefix.
func WithEnvPrefix(prefix string) Option {
	return func(cl *ConfigLoader) {
		cl.envVarPrefix = prefix
	}
}

// WithConfigPaths sets the configuration file search paths.
func WithConfigPaths(paths ...string) Option {
	return func(cl *ConfigLoader) {
		cl.configPaths = paths
	}
}

// WithConfigType sets the configuration file type.
func WithConfigType(configType string) Option {
	return func(cl *ConfigLoader) {
		cl.configType = configType
	}
}

// WithConfigFile enables loading from config file.
func WithConfigFile(enabled bool) LoadOption {
	return func(o *loadOptions) {
		o.loadConfigFile = enabled
	}
}

// WithCustomConfigFile enables loading from a custom config file name.
func WithCustomConfigFile(fileName string) LoadOption {
	return func(o *loadOptions) {
		o.loadConfigFile = true
		o.configFileName = fileName
	}
}

// NewConfigLoader creates a new ConfigLoader with the given options.
func NewConfigLoader(options ...Option) *ConfigLoader {
	cl := &ConfigLoader{
		envVarPrefix: "",
		configPaths:  []string{"./configs", "/configs"},
		configType:   "toml",
	}

	for _, option := range options {
		option(cl)
	}

	return cl
}

// AddConfigPath adds a configuration file search path.
func (cl *ConfigLoader) AddConfigPath(path string) *ConfigLoader {
	cl.configPaths = append(cl.configPaths, path)
	return cl
}

// Validate checks if the configuration values are valid.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.AppName) == "" {
		return errors.New("application name cannot be empty")
	}
	if strings.TrimSpace(c.EnvName) == "" {
		return errors.New("environment name cannot be empty")
	}
	return nil
}

// bindFlags adds all the flags to the provided flag set.
func (cl *ConfigLoader) bindFlags(fs *pflag.FlagSet, config *Config) {
	fs.StringVar(&config.AppName, AppName, config.AppName, "The name of the application.")
	fs.StringVar(&config.EnvName, EnvName, config.EnvName,
		"The environment of the application. Used to load the right config file.")
}

// normalizeFlags changes all flags that contain "_" separators to use "-".
func normalizeFlags(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	}
	return pflag.NormalizedName(name)
}

// setupViper configures viper with environment variable settings.
func (cl *ConfigLoader) setupViper() {
	if cl.envVarPrefix != "" {
		viper.SetEnvPrefix(cl.envVarPrefix)
	}
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}

// loadConfigFile loads the configuration file using viper.
func (cl *ConfigLoader) loadConfigFile(configName string) error {
	viper.SetConfigName(configName)
	viper.SetConfigType(cl.configType)

	for _, path := range cl.configPaths {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return fmt.Errorf("config file '%s.%s' not found in paths %v: %w",
				configName, cl.configType, cl.configPaths, err)
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return nil
}

// LoadConfig loads configuration with the specified options.
func (cl *ConfigLoader) LoadConfig(options []LoadOption, flagSets ...func(fs *pflag.FlagSet)) (*Config, error) {
	opts := &loadOptions{
		loadConfigFile: false,
		configFileName: "",
	}

	for _, option := range options {
		option(opts)
	}

	fs := pflag.NewFlagSet("config", pflag.ExitOnError)

	config := &Config{
		AppName: DefaultAppName,
		EnvName: DefaultEnvName,
	}

	cl.bindFlags(fs, config)

	for _, flagSet := range flagSets {
		flagSet(fs)
	}

	fs.SetNormalizeFunc(normalizeFlags)

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	if err := viper.BindPFlags(fs); err != nil {
		return nil, fmt.Errorf("failed to bind flags to viper: %w", err)
	}

	cl.setupViper()

	if opts.loadConfigFile {
		configName := opts.configFileName
		if configName == "" {
			envName := viper.GetString(EnvName)
			configName = fmt.Sprintf("config.%s", strings.ToLower(envName))
		}

		if err := cl.loadConfigFile(configName); err != nil {
			return nil, err
		}
	}

	config.AppName = viper.GetString(AppName)
	config.EnvName = viper.GetString(EnvName)

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}
