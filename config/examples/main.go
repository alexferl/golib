package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexferl/golib/config"
	"github.com/spf13/pflag"
)

func main() {
	fmt.Println("Try these commands:")
	fmt.Println("  go run main.go")
	fmt.Println("  go run main.go --app-name myservice --env-name prod --port 3000 --debug")
	fmt.Println("  MYAPP_APP_NAME=webapi MYAPP_ENV_NAME=staging go run main.go")
	fmt.Println("  go run main.go --env-name prod  # loads config.prod.toml")
	fmt.Println("  go run main.go --env-name dev   # loads config.dev.toml")
	fmt.Println("  go run main.go --help")
	fmt.Println()

	loader := config.NewConfigLoader(
		config.WithEnvPrefix("MYAPP"),
		config.WithConfigPaths("./configs"),
		config.WithConfigType("toml"),
	)

	var port int
	var debug bool
	var version bool

	cfg, err := loader.LoadConfig(
		[]config.LoadOption{config.WithConfigFile(true)},
		func(fs *pflag.FlagSet) {
			fs.IntVar(&port, "port", 8080, "Server port")
			fs.BoolVar(&debug, "debug", false, "Enable debug mode")
			fs.BoolVar(&version, "version", false, "Show version")
		},
	)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if version {
		fmt.Println("Version: 1.0.0")
		os.Exit(0)
	}

	fmt.Printf("\nConfiguration loaded:\n")
	fmt.Printf("  App Name: %s\n", cfg.AppName)
	fmt.Printf("  Environment: %s\n", cfg.EnvName)
	fmt.Printf("  Port: %d\n", port)
	fmt.Printf("  Debug: %t\n", debug)
	fmt.Println()

	fmt.Println("Configuration sources (in precedence order):")
	fmt.Println("  1. Command line flags (highest)")
	fmt.Println("  2. Environment variables (MYAPP_*)")
	fmt.Printf("  3. Config file (config.%s.toml)\n", cfg.EnvName)
	fmt.Println("  4. Default values (lowest)")
	fmt.Println()

	fmt.Printf("Starting %s server in %s environment on port %d\n",
		cfg.AppName, cfg.EnvName, port)

	if debug {
		fmt.Println("🐛 Debug mode is enabled")
	}

	fmt.Println("✅ Application started successfully!")
}
