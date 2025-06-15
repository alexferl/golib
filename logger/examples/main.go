package main

import (
	"log"

	"github.com/alexferl/golib/logger"
)

func main() {
	// Example 1: JSON format to stderr with DEBUG level
	config1 := &logger.Config{
		LogLevel:  "DEBUG",
		LogFormat: "json",
		LogOutput: "stderr",
	}

	jsonLogger, err := logger.New(config1)
	if err != nil {
		log.Fatal(err)
	}

	jsonLogger.Debug().Str("format", "json").Msg("This is JSON formatted")

	// Example 2: Text format to stdout with INFO level
	config2 := &logger.Config{
		LogLevel:  "INFO",
		LogFormat: "text",
		LogOutput: "stdout",
	}

	textLogger, err := logger.New(config2)
	if err != nil {
		log.Fatal(err)
	}

	textLogger.Info().Str("format", "text").Msg("This is human-readable text")

	// Example 3: Using default config
	defaultLogger, err := logger.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	defaultLogger.Info().Msg("Using default configuration")
}
