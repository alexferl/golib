package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alexferl/golib/logger"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNew(t *testing.T) {
	config := Config{
		Name:    "test",
		Version: "1.0.0",
		HTTP: HTTPConfig{
			BindAddr: ":8080",
		},
	}

	server := New(config)

	if server == nil {
		t.Fatal("New() returned nil")
	}

	if server.config.Name != "test" {
		t.Errorf("Expected config name 'test', got '%s'", server.config.Name)
	}

	if server.echo == nil {
		t.Error("Echo instance is nil")
	}

	if server.logger == nil {
		t.Error("Logger instance is nil")
	}
}

func TestWithLogger(t *testing.T) {
	config := Config{}
	customLogger, err := logger.New(logger.DefaultConfig)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	server := New(config, WithLogger(customLogger))

	if server.logger != customLogger {
		t.Error("Custom logger was not set")
	}
}

func TestWithMiddleware(t *testing.T) {
	config := Config{}
	middlewareCalled := false

	testMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			middlewareCalled = true
			return next(c)
		}
	}

	server := New(config, WithMiddleware(testMiddleware))

	e := server.Echo()
	e.GET("/test", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if !middlewareCalled {
		t.Error("Middleware was not called")
	}

	if rec.Code != 200 {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestWithEchoConfig(t *testing.T) {
	config := Config{}
	configCalled := false

	server := New(config, WithEchoConfig(func(e *echo.Echo) {
		configCalled = true
		e.Debug = true
	}))

	if !configCalled {
		t.Error("Echo config function was not called")
	}

	if !server.echo.Debug {
		t.Error("Echo debug was not set")
	}
}

func TestEcho(t *testing.T) {
	config := Config{}
	server := New(config)

	e := server.Echo()
	if e == nil {
		t.Error("Echo() returned nil")
	}

	if e != server.echo {
		t.Error("Echo() returned different instance")
	}
}

func TestLogger(t *testing.T) {
	config := Config{}
	server := New(config)

	l := server.Logger()
	if l == nil {
		t.Error("Logger() returned nil")
	}

	if l != server.logger {
		t.Error("Logger() returned different instance")
	}
}

func TestShutdown(t *testing.T) {
	config := Config{
		HTTP: HTTPConfig{
			BindAddr: ":0", // Use random port
		},
	}

	server := New(config)

	errCh := server.Start()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	// Check if error channel is closed
	select {
	case _, ok := <-errCh:
		if ok {
			t.Error("Error channel should be closed after shutdown")
		}
	default:
		// Channel might be buffered, this is OK
	}
}

func TestPrepareHandler(t *testing.T) {
	config := Config{
		Compress: CompressConfig{
			Enabled: false,
		},
	}

	server := New(config)
	handler := server.prepareHandler()

	if handler == nil {
		t.Error("prepareHandler() returned nil")
	}
}

func TestPrepareHandlerWithCompression(t *testing.T) {
	config := Config{
		Compress: CompressConfig{
			Enabled:   true,
			Level:     6,
			MinLength: 1024,
		},
	}

	server := New(config)
	handler := server.prepareHandler()

	if handler == nil {
		t.Error("prepareHandler() with compression returned nil")
	}
}

func TestCreateHTTPServer(t *testing.T) {
	config := Config{
		HTTP: HTTPConfig{
			BindAddr:          ":8080",
			IdleTimeout:       60 * time.Second,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			MaxHeaderBytes:    1024,
		},
	}

	server := New(config)
	httpServer := server.createHTTPServer(":8080", http.DefaultServeMux)

	if httpServer != nil {
		if httpServer.Addr != ":8080" {
			t.Errorf("Expected address ':8080', got '%s'", httpServer.Addr)
		}

		if httpServer.IdleTimeout != 60*time.Second {
			t.Errorf("Expected IdleTimeout 60s, got %v", httpServer.IdleTimeout)
		}
	} else {
		t.Error("createHTTPServer() returned nil")
	}
}

func TestHealthcheckEndpoints(t *testing.T) {
	config := *DefaultConfig
	server := New(config)

	req := httptest.NewRequest("GET", config.Healthcheck.LivenessEndpoint, nil)
	rec := httptest.NewRecorder()
	c := server.echo.NewContext(req, rec)

	err := config.Healthcheck.LivenessHandler(c)
	if err != nil {
		t.Errorf("LivenessHandler returned error: %v", err)
	}

	if rec.Code != 200 {
		t.Errorf("Expected status 200 for liveness, got %d", rec.Code)
	}

	req = httptest.NewRequest("GET", config.Healthcheck.ReadinessEndpoint, nil)
	rec = httptest.NewRecorder()
	c = server.echo.NewContext(req, rec)

	err = config.Healthcheck.ReadinessHandler(c)
	if err != nil {
		t.Errorf("ReadinessHandler returned error: %v", err)
	}

	if rec.Code != 200 {
		t.Errorf("Expected status 200 for readiness, got %d", rec.Code)
	}

	req = httptest.NewRequest("GET", config.Healthcheck.StartupEndpoint, nil)
	rec = httptest.NewRecorder()
	c = server.echo.NewContext(req, rec)

	err = config.Healthcheck.StartupHandler(c)
	if err != nil {
		t.Errorf("StartupHandler returned error: %v", err)
	}

	if rec.Code != 200 {
		t.Errorf("Expected status 200 for startup, got %d", rec.Code)
	}
}

func TestDefaultHealthcheckHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/livez", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := defaultHealthcheckHandler(c)
	if err != nil {
		t.Errorf("defaultHealthcheckHandler returned error: %v", err)
	}

	if rec.Code != 200 {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	expectedBody := "ok"
	if strings.TrimSpace(rec.Body.String()) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, rec.Body.String())
	}
}

func TestHealthcheckEndpointsRegistered(t *testing.T) {
	config := *DefaultConfig
	server := New(config)

	e := server.Echo()

	req := httptest.NewRequest("GET", "/livez", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Expected status 200 for liveness endpoint, got %d", rec.Code)
	}

	req = httptest.NewRequest("GET", "/readyz", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Expected status 200 for readiness endpoint, got %d", rec.Code)
	}

	req = httptest.NewRequest("GET", "/startupz", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Expected status 200 for startup endpoint, got %d", rec.Code)
	}
}

func resetPrometheusRegistry() {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	prometheus.DefaultGatherer = prometheus.DefaultRegisterer.(prometheus.Gatherer)
}

func TestPrometheusEndpoint(t *testing.T) {
	resetPrometheusRegistry()
	defer resetPrometheusRegistry()

	config := Config{
		Name: "testapp",
		Prometheus: PrometheusConfig{
			Enabled: true,
			Path:    "/metrics",
		},
	}

	server := New(config)
	e := server.Echo()

	req := httptest.NewRequest("GET", "/metrics", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Expected status 200 for Prometheus endpoint, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "# HELP") || !strings.Contains(body, "# TYPE") {
		t.Error("Response does not contain Prometheus metrics format")
	}
}

func TestPrometheusDisabled(t *testing.T) {
	resetPrometheusRegistry()
	defer resetPrometheusRegistry()

	config := Config{
		Name: "testapp",
		Prometheus: PrometheusConfig{
			Enabled: false,
			Path:    "/metrics",
		},
	}

	server := New(config)
	e := server.Echo()

	req := httptest.NewRequest("GET", "/metrics", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != 404 {
		t.Errorf("Expected status 404 for disabled Prometheus endpoint, got %d", rec.Code)
	}
}

func TestPrometheusCustomPath(t *testing.T) {
	resetPrometheusRegistry()
	defer resetPrometheusRegistry()

	config := Config{
		Name: "testapp",
		Prometheus: PrometheusConfig{
			Enabled: true,
			Path:    "/custom/metrics",
		},
	}

	server := New(config)
	e := server.Echo()

	req := httptest.NewRequest("GET", "/custom/metrics", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Expected status 200 for custom Prometheus endpoint, got %d", rec.Code)
	}

	req = httptest.NewRequest("GET", "/metrics", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != 404 {
		t.Errorf("Expected status 404 for default metrics path when custom path is used, got %d", rec.Code)
	}
}
