package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexferl/golib/config"
	"github.com/alexferl/golib/http/middleware"
	"github.com/alexferl/golib/http/server"
	"github.com/alexferl/golib/logger"
	"github.com/labstack/echo/v4"
	"github.com/spf13/pflag"
)

type AppConfig struct {
	config.Config
	Server        *server.Config
	Logger        *logger.Config
	RequestID     *middleware.RequestID
	RequestLogger *middleware.RequestLogger
}

func main() {
	appConfig := &AppConfig{
		Config:        config.Config{AppName: "myapp", EnvName: "local"},
		Server:        server.DefaultConfig,
		Logger:        logger.DefaultConfig,
		RequestID:     middleware.DefaultRequestID,
		RequestLogger: middleware.DefaultLogger,
	}

	configLoader := config.NewConfigLoader()

	_, err := configLoader.LoadConfig(
		[]config.LoadOption{},
		func(fs *pflag.FlagSet) {
			fs.AddFlagSet(appConfig.Server.FlagSet())
			fs.AddFlagSet(appConfig.Logger.FlagSet())
			fs.AddFlagSet(appConfig.RequestID.FlagSet())
			fs.AddFlagSet(appConfig.RequestLogger.FlagSet())
		},
	)
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	appLogger, err := logger.New(appConfig.Logger)
	if err != nil {
		log.Fatal("failed to create logger:", err)
	}

	appConfig.RequestLogger.Logger = appLogger

	var middlewares []echo.MiddlewareFunc
	if appConfig.RequestID.Enabled {
		middlewares = append(middlewares, middleware.NewRequestID(appConfig.RequestID))
	}

	if appConfig.RequestLogger.Enabled {
		middlewares = append(middlewares, middleware.NewRequestLogger(appConfig.RequestLogger))
	}

	srv := server.New(*appConfig.Server,
		server.WithLogger(appLogger),
		server.WithMiddleware(middlewares...),
	)

	srv.Echo().GET("/", func(c echo.Context) error {
		srv.Logger().Info().Str("endpoint", "/").Msg("Hello endpoint called")
		return c.JSON(200, map[string]string{"message": "Hello World"})
	})

	errCh := srv.Start()

	srv.Logger().Info().
		Str("name", appConfig.Server.Name).
		Str("version", appConfig.Server.Version).
		Str("addr", appConfig.Server.HTTP.BindAddr).
		Bool("tls", appConfig.Server.TLS.Enabled).
		Msg("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		srv.Logger().Fatal().Err(err).Msg("server error")
	case <-quit:
		srv.Logger().Info().Msg("received shutdown signal")

		ctx, cancel := context.WithTimeout(context.Background(), appConfig.Server.GracefulTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			srv.Logger().Fatal().Err(err).Msg("server forced to shutdown")
		}

		srv.Logger().Info().Msg("server exited")
	}
}
