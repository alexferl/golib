package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"github.com/alexferl/golib/http/config"
	"github.com/alexferl/golib/http/middleware"
	"github.com/alexferl/golib/http/router"
)

type Server struct {
	*echo.Echo
	*config.Config
}

func New() *Server {
	e := echo.New()
	return &Server{e, config.DefaultConfig}
}

// Start starts the echo HTTP server.
func (s *Server) Start(r *router.Router) {
	middleware.Register(s.Echo)
	router.Register(s.Echo, r)

	// Start server
	go func() {
		addr := fmt.Sprintf("%s:%s", viper.GetString("http-bind-address"), viper.GetString("http-bind-port"))
		if err := s.Echo.Start(addr); err != nil {
			s.Echo.Logger.Info("Received signal, shutting down the server")
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	timeout := time.Duration(viper.GetInt64("http-graceful-timeout")) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.Echo.Shutdown(ctx); err != nil {
		s.Echo.Logger.Fatal(err)
	}
}
