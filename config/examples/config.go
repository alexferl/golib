package main

import (
	"fmt"

	libconfig "github.com/alexferl/golib/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program.
type Config struct {
	*libconfig.Config
	MyKey string
}

// NewConfig creates a Config instance.
func NewConfig() *Config {
	c := &Config{
		Config: libconfig.New("app"),
		MyKey:  "value",
	}
	return c
}

const (
	MyKey = "my-key"
)

// addFlags adds all the flags from the command line.
func (c *Config) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.MyKey, MyKey, c.MyKey, "My key.")
}

// BindFlags normalizes and parses the command line flags.
func (c *Config) BindFlags() {
	c.addFlags(pflag.CommandLine)
	err := c.Config.BindFlags() // Bind the default flags from x/config
	if err != nil {
		panic(fmt.Sprintf("failed binding flags: %v\n", err))
	}
}

func main() {
	c := NewConfig()
	c.BindFlags()
	fmt.Println(viper.GetString(libconfig.AppName)) // from libconfig, overloaded in configs/config.dev.toml
	fmt.Println(viper.GetString(MyKey))
}
