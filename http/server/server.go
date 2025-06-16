package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/alexferl/golib/logger"
	"github.com/klauspost/compress/gzhttp"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/ziflex/lecho/v3"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// Server represents an HTTP/HTTPS server with configurable middleware and TLS support.
type Server struct {
	config      Config
	logger      *logger.Logger
	echo        *echo.Echo
	httpServer  *http.Server
	httpsServer *http.Server
	errCh       chan error
	ctx         context.Context
	cancel      context.CancelFunc
}

// Option is a function that configures a Server.
type Option func(*Server)

// WithEchoConfig allows custom Echo configuration.
func WithEchoConfig(configFunc func(*echo.Echo)) Option {
	return func(s *Server) {
		configFunc(s.echo)
	}
}

// WithLogger sets a custom logger for the server.
func WithLogger(logger *logger.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithMiddleware adds middleware to the Echo instance.
func WithMiddleware(middlewares ...echo.MiddlewareFunc) Option {
	return func(s *Server) {
		s.echo.Use(middlewares...)
	}
}

// New creates a new Server instance with the given configuration and options.
func New(config Config, options ...Option) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	server := &Server{
		config: config,
		echo:   e,
		errCh:  make(chan error, 10),
		ctx:    ctx,
		cancel: cancel,
	}

	for _, option := range options {
		option(server)
	}

	if server.logger == nil {
		defaultLogger, err := logger.New(logger.DefaultConfig)
		if err != nil {
			panic(err)
		}
		server.logger = defaultLogger
	}

	server.echo.Logger = lecho.From(server.logger.GetLogger())

	server.echo.GET(config.Healthcheck.LivenessEndpoint, config.Healthcheck.LivenessHandler)
	server.echo.GET(config.Healthcheck.ReadinessEndpoint, config.Healthcheck.ReadinessHandler)
	server.echo.GET(config.Healthcheck.StartupEndpoint, config.Healthcheck.StartupHandler)

	if config.Prometheus.Enabled {
		server.echo.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
			Namespace: "",
			Subsystem: config.Name,
		}))
		server.echo.GET(config.Prometheus.Path, echoprometheus.NewHandler())
	}

	return server
}

// Echo returns the underlying Echo instance.
func (s *Server) Echo() *echo.Echo {
	return s.echo
}

// Logger returns the logger instance used by the server.
func (s *Server) Logger() *logger.Logger {
	return s.logger
}

// Start starts the HTTP or HTTPS server and returns a channel for errors.
func (s *Server) Start() <-chan error {
	handler := s.prepareHandler()

	s.httpServer = s.createHTTPServer(s.config.HTTP.BindAddr, handler)

	if !s.config.TLS.Enabled {
		s.startHTTPServer()
	} else {
		s.startHTTPSServer(handler)
	}

	return s.errCh
}

// Shutdown gracefully shuts down the HTTP and HTTPS servers.
func (s *Server) Shutdown(ctx context.Context) error {
	// signal all goroutines to stop
	s.cancel()

	var errs []error

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("HTTP server shutdown error: %w", err))
		}
	}

	if s.httpsServer != nil {
		if err := s.httpsServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("HTTPS server shutdown error: %w", err))
		}
	}

	close(s.errCh)

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// prepareHandler prepares the HTTP handler with optional gzip compression.
func (s *Server) prepareHandler() http.Handler {
	handler := http.Handler(s.echo)

	if s.config.Compress.Enabled {
		gzipHandler, err := gzhttp.NewWrapper(
			gzhttp.MinSize(s.config.Compress.MinLength),
			gzhttp.CompressionLevel(s.config.Compress.Level),
		)
		if err != nil {
			s.errCh <- fmt.Errorf("gzip handler error: %w", err)
			return handler
		}
		handler = gzipHandler(s.echo)
	}

	return handler
}

// createHTTPServer creates a new HTTP server with the given address and handler.
func (s *Server) createHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		IdleTimeout:       s.config.HTTP.IdleTimeout,
		ReadTimeout:       s.config.HTTP.ReadTimeout,
		ReadHeaderTimeout: s.config.HTTP.ReadHeaderTimeout,
		WriteTimeout:      s.config.HTTP.WriteTimeout,
		MaxHeaderBytes:    s.config.HTTP.MaxHeaderBytes,
	}
}

// startHTTPServer starts the HTTP server in a new goroutine.
func (s *Server) startHTTPServer() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errCh <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()
}

// startHTTPSServer starts the HTTPS server with TLS configuration.
func (s *Server) startHTTPSServer(handler http.Handler) {
	s.httpsServer = s.createHTTPServer(s.config.TLS.BindAddr, handler)

	if s.config.TLS.ACME.Enabled {
		s.setupACME()
	} else {
		s.setupManualTLS()
	}

	if s.config.Redirect.HTTPS && !s.config.TLS.ACME.Enabled {
		s.echo.Pre(s.redirectToHTTPS)
		s.startHTTPServer()
	}
}

// setupACME configures the server to use ACME/Let's Encrypt for TLS certificates.
func (s *Server) setupACME() {
	acmeClient := &acme.Client{}
	autocertManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Email:      s.config.TLS.ACME.Email,
		HostPolicy: autocert.HostWhitelist(s.config.TLS.ACME.HostWhitelist...),
		Cache:      autocert.DirCache(s.config.TLS.ACME.CachePath),
	}

	if s.config.TLS.ACME.DirectoryURL != "" {
		acmeClient.DirectoryURL = s.config.TLS.ACME.DirectoryURL
	}

	autocertManager.Client = acmeClient

	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: defaultCurves,
		CipherSuites:     getOptimalDefaultCipherSuites(),
		GetCertificate:   autocertManager.GetCertificate,
	}

	s.httpsServer.TLSConfig = tlsConfig

	// HTTP server that listens on port 80 for challenges
	_, port, err := net.SplitHostPort(s.config.HTTP.BindAddr)
	if err != nil {
		s.errCh <- fmt.Errorf("failed to split host/port: %w", err)
		return
	}

	if port != "80" {
		s.errCh <- fmt.Errorf("bind-addr must be set to port 80 for the challenge server")
		return
	}

	s.httpServer.Handler = autocertManager.HTTPHandler(nil)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errCh <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Start HTTPS server
	_, tlsPort, err := net.SplitHostPort(s.config.TLS.BindAddr)
	if err != nil {
		s.errCh <- fmt.Errorf("failed to split host/port: %w", err)
		return
	}

	if tlsPort != "443" {
		s.errCh <- fmt.Errorf("tls-bind-addr must be set to port 443 for auto TLS")
		return
	}

	go func() {
		if err := s.httpsServer.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errCh <- fmt.Errorf("HTTPS server error: %w", err)
		}
	}()
}

// setupManualTLS configures the server to use manual TLS certificates.
func (s *Server) setupManualTLS() {
	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: defaultCurves,
		CipherSuites:     getOptimalDefaultCipherSuites(),
	}

	s.httpsServer.TLSConfig = tlsConfig

	go func() {
		if err := s.httpsServer.ListenAndServeTLS(
			s.config.TLS.CertFile,
			s.config.TLS.KeyFile,
		); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errCh <- fmt.Errorf("HTTPS server error: %w", err)
		}
	}()
}

// redirectToHTTPS redirects HTTP requests to HTTPS.
func (s *Server) redirectToHTTPS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req, scheme := c.Request(), c.Scheme()
		if scheme != "https" {
			host := req.Host

			if h, _, err := net.SplitHostPort(host); err == nil {
				host = h
			}

			_, tlsPort, err := net.SplitHostPort(s.config.TLS.BindAddr)
			if err != nil {
				return err
			}

			// if TLS port is the default (443), don't include it in the URL
			portSuffix := ""
			if tlsPort != "443" {
				portSuffix = ":" + tlsPort
			}

			url := fmt.Sprintf("https://%s%s%s", host, portSuffix, req.RequestURI)
			return c.Redirect(s.config.Redirect.Code, url)
		}

		return next(c)
	}
}
