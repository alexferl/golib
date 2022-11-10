package main

import (
	"fmt"

	libLog "github.com/alexferl/golib/log"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

func main() {
	c := libLog.DefaultConfig // use default settings
	// c := &libLog.Config{LogWriter: "json"} // use custom settings
	c.BindFlags(pflag.CommandLine)
	pflag.Parse()

	err := libLog.New(c)
	if err != nil {
		panic(fmt.Sprintf("Error initializing logger: '%v'", err))
	}

	log.Info().Msg("Hello, world!")
	log.Info().Msgf("Hello, %s!", "world")
	log.Warn().Msg("Hello, warn!")
	log.Error().Msg("Hello, error!")
}
