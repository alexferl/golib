package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all global configuration for our program
type Config struct {
	AppName      string
	EnvName      string
	EnvVarPrefix string
}

var DefaultConfig = &Config{
	AppName:      "app",
	EnvName:      "dev",
	EnvVarPrefix: "app",
}

// New creates a Config instance
func New() *Config {
	return DefaultConfig
}

// bindFlags adds all the flags from the command line
func (c *Config) bindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.AppName, "app-name", c.AppName, "The name of the application.")
	fs.StringVar(&c.EnvName, "env-name", c.EnvName, "The environment of the application. "+
		"Used to load the right configs file.")
	fs.StringVar(&c.EnvVarPrefix, "env-var-prefix", c.EnvVarPrefix,
		"Used to prefix environment variables.")
}

// wordSepNormalizeFunc changes all flags that contain "_" separators
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}

// BindFlags normalizes and parses the command line flags
func (c *Config) BindFlags(flagSets ...func(fs *pflag.FlagSet)) error {
	for _, flagSet := range flagSets {
		flagSet(pflag.CommandLine)
	}

	c.bindFlags(pflag.CommandLine)
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return err
	}

	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	pflag.Parse()

	n := viper.GetString("app-name")
	if len(n) < 1 {
		return errors.New("application name cannot be empty")
	}

	viper.SetEnvPrefix(n)
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	configName := fmt.Sprintf("config.%s", strings.ToLower(viper.GetString("env-name")))
	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/configs")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.New(fmt.Sprintf("config file not found: '%v'", err))
		} else {
			return errors.New(fmt.Sprintf("couldn't load config file: '%v'", err))
		}
	}

	return nil
}
