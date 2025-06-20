package server

import (
	"net/http"
	"testing"
	"time"
)

func TestConfig_FlagSet(t *testing.T) {
	config := &Config{
		Name:            "testapp",
		Version:         "2.0.0",
		GracefulTimeout: 45 * time.Second,
		HTTP: HTTPConfig{
			BindAddr:          ":9000",
			IdleTimeout:       90 * time.Second,
			ReadTimeout:       15 * time.Second,
			ReadHeaderTimeout: 8 * time.Second,
			WriteTimeout:      12 * time.Second,
			MaxHeaderBytes:    2 << 20,
		},
		TLS: TLSConfig{
			Enabled:  true,
			BindAddr: ":9443",
			CertFile: "/path/to/cert.pem",
			KeyFile:  "/path/to/key.pem",
			ACME: ACMEConfig{
				Enabled:       true,
				Email:         "test@example.com",
				HostWhitelist: []string{"example.com", "www.example.com"},
				CachePath:     "/tmp/certs",
				DirectoryURL:  "https://acme-staging-v02.api.letsencrypt.org/directory",
			},
		},
		Compress: CompressConfig{
			Enabled:   true,
			Level:     9,
			MinLength: 2048,
		},
		Redirect: RedirectConfig{
			HTTPS: true,
			Code:  302,
		},
		Healthcheck: HealthcheckConfig{
			LivenessEndpoint:  "/health/live",
			ReadinessEndpoint: "/health/ready",
			StartupEndpoint:   "/health/startup",
		},
		Prometheus: PrometheusConfig{
			Enabled: true,
			Path:    "/custom/metrics",
		},
	}

	fs := config.FlagSet()

	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	// Test basic flags
	nameFlag := fs.Lookup(ServerName)
	if nameFlag == nil {
		t.Errorf("Flag %s not found", ServerName)
	} else {
		if nameFlag.DefValue != "testapp" {
			t.Errorf("Flag %s default value = %v, want testapp", ServerName, nameFlag.DefValue)
		}
	}

	versionFlag := fs.Lookup(ServerVersion)
	if versionFlag == nil {
		t.Errorf("Flag %s not found", ServerVersion)
	} else {
		if versionFlag.DefValue != "2.0.0" {
			t.Errorf("Flag %s default value = %v, want 2.0.0", ServerVersion, versionFlag.DefValue)
		}
	}

	gracefulTimeoutFlag := fs.Lookup(ServerGracefulTimeout)
	if gracefulTimeoutFlag == nil {
		t.Errorf("Flag %s not found", ServerGracefulTimeout)
	} else {
		if gracefulTimeoutFlag.DefValue != "45s" {
			t.Errorf("Flag %s default value = %v, want 45s", ServerGracefulTimeout, gracefulTimeoutFlag.DefValue)
		}
	}

	// Test HTTP flags
	httpBindAddrFlag := fs.Lookup(ServerHTTPBindAddr)
	if httpBindAddrFlag == nil {
		t.Errorf("Flag %s not found", ServerHTTPBindAddr)
	} else {
		if httpBindAddrFlag.DefValue != ":9000" {
			t.Errorf("Flag %s default value = %v, want :9000", ServerHTTPBindAddr, httpBindAddrFlag.DefValue)
		}
	}

	// Test TLS flags
	tlsEnabledFlag := fs.Lookup(ServerTLSEnabled)
	if tlsEnabledFlag == nil {
		t.Errorf("Flag %s not found", ServerTLSEnabled)
	} else {
		if tlsEnabledFlag.DefValue != "true" {
			t.Errorf("Flag %s default value = %v, want true", ServerTLSEnabled, tlsEnabledFlag.DefValue)
		}
	}

	// Test ACME flags
	acmeEmailFlag := fs.Lookup(ServerTLSACMEEmail)
	if acmeEmailFlag == nil {
		t.Errorf("Flag %s not found", ServerTLSACMEEmail)
	} else {
		if acmeEmailFlag.DefValue != "test@example.com" {
			t.Errorf("Flag %s default value = %v, want test@example.com", ServerTLSACMEEmail, acmeEmailFlag.DefValue)
		}
	}

	// Test Healthcheck flags
	livenessFlag := fs.Lookup(ServerHealthcheckLivenessEndpoint)
	if livenessFlag == nil {
		t.Errorf("Flag %s not found", ServerHealthcheckLivenessEndpoint)
	} else {
		if livenessFlag.DefValue != "/health/live" {
			t.Errorf("Flag %s default value = %v, want /health/live", ServerHealthcheckLivenessEndpoint, livenessFlag.DefValue)
		}
	}

	readinessFlag := fs.Lookup(ServerHealthcheckReadinessEndpoint)
	if readinessFlag == nil {
		t.Errorf("Flag %s not found", ServerHealthcheckReadinessEndpoint)
	} else {
		if readinessFlag.DefValue != "/health/ready" {
			t.Errorf("Flag %s default value = %v, want /health/ready", ServerHealthcheckReadinessEndpoint, readinessFlag.DefValue)
		}
	}

	startupFlag := fs.Lookup(ServerHealthcheckStartupEndpoint)
	if startupFlag == nil {
		t.Errorf("Flag %s not found", ServerHealthcheckStartupEndpoint)
	} else {
		if startupFlag.DefValue != "/health/startup" {
			t.Errorf("Flag %s default value = %v, want /health/startup", ServerHealthcheckStartupEndpoint, startupFlag.DefValue)
		}
	}

	// Test Prometheus flags
	prometheusEnabledFlag := fs.Lookup(ServerPrometheusEnabled)
	if prometheusEnabledFlag == nil {
		t.Errorf("Flag %s not found", ServerPrometheusEnabled)
	} else {
		if prometheusEnabledFlag.DefValue != "true" {
			t.Errorf("Flag %s default value = %v, want true", ServerPrometheusEnabled, prometheusEnabledFlag.DefValue)
		}
	}

	prometheusPathFlag := fs.Lookup(ServerPrometheusPath)
	if prometheusPathFlag == nil {
		t.Errorf("Flag %s not found", ServerPrometheusPath)
	} else {
		if prometheusPathFlag.DefValue != "/custom/metrics" {
			t.Errorf("Flag %s default value = %v, want /custom/metrics", ServerPrometheusPath, prometheusPathFlag.DefValue)
		}
	}
}

func TestConfig_FlagSet_Parse(t *testing.T) {
	config := &Config{
		Name:            "app",
		Version:         "1.0.0",
		GracefulTimeout: 30 * time.Second,
		HTTP: HTTPConfig{
			BindAddr:          "localhost:8080",
			IdleTimeout:       60 * time.Second,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
		TLS: TLSConfig{
			Enabled:  false,
			BindAddr: "localhost:8443",
		},
		Compress: CompressConfig{
			Enabled:   false,
			Level:     6,
			MinLength: 1024,
		},
		Redirect: RedirectConfig{
			HTTPS: false,
			Code:  301,
		},
		Healthcheck: HealthcheckConfig{
			LivenessEndpoint:  "/livez",
			ReadinessEndpoint: "/readyz",
			StartupEndpoint:   "/startupz",
		},
		Prometheus: PrometheusConfig{
			Enabled: false,
			Path:    "/metrics",
		},
	}

	fs := config.FlagSet()

	args := []string{
		"--server-name", "myapp",
		"--server-version", "3.0.0",
		"--server-graceful-timeout", "60s",
		"--server-http-bind-addr", ":8081",
		"--server-http-idle-timeout", "120s",
		"--server-http-read-timeout", "20s",
		"--server-http-read-header-timeout", "10s",
		"--server-http-write-timeout", "15s",
		"--server-http-max-header-bytes", "2097152",
		"--server-tls-enabled",
		"--server-tls-bind-addr", ":8444",
		"--server-tls-cert-file", "/etc/ssl/cert.pem",
		"--server-tls-key-file", "/etc/ssl/key.pem",
		"--server-tls-acme-enabled",
		"--server-tls-acme-email", "admin@example.com",
		"--server-tls-acme-host-whitelist", "example.com,api.example.com",
		"--server-tls-acme-cache-path", "/var/cache/certs",
		"--server-tls-acme-directory-url", "https://acme-v02.api.letsencrypt.org/directory",
		"--server-compress-enabled",
		"--server-compress-level", "8",
		"--server-compress-min-length", "512",
		"--server-redirect-https",
		"--server-redirect-code", "308",
		"--server-healthcheck-liveness-endpoint", "/custom/live",
		"--server-healthcheck-readiness-endpoint", "/custom/ready",
		"--server-healthcheck-startup-endpoint", "/custom/startup",
		"--server-prometheus-enabled",
		"--server-prometheus-path", "/custom/metrics",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Verify basic config
	if config.Name != "myapp" {
		t.Errorf("Name = %v, want myapp", config.Name)
	}
	if config.Version != "3.0.0" {
		t.Errorf("Version = %v, want 3.0.0", config.Version)
	}
	if config.GracefulTimeout != 60*time.Second {
		t.Errorf("GracefulTimeout = %v, want 60s", config.GracefulTimeout)
	}

	// Verify HTTP config
	if config.HTTP.BindAddr != ":8081" {
		t.Errorf("HTTP.BindAddr = %v, want :8081", config.HTTP.BindAddr)
	}
	if config.HTTP.IdleTimeout != 120*time.Second {
		t.Errorf("HTTP.IdleTimeout = %v, want 120s", config.HTTP.IdleTimeout)
	}
	if config.HTTP.MaxHeaderBytes != 2097152 {
		t.Errorf("HTTP.MaxHeaderBytes = %v, want 2097152", config.HTTP.MaxHeaderBytes)
	}

	// Verify TLS config
	if !config.TLS.Enabled {
		t.Errorf("TLS.Enabled = %v, want true", config.TLS.Enabled)
	}
	if config.TLS.BindAddr != ":8444" {
		t.Errorf("TLS.BindAddr = %v, want :8444", config.TLS.BindAddr)
	}
	if config.TLS.CertFile != "/etc/ssl/cert.pem" {
		t.Errorf("TLS.CertFile = %v, want /etc/ssl/cert.pem", config.TLS.CertFile)
	}

	// Verify ACME config
	if !config.TLS.ACME.Enabled {
		t.Errorf("TLS.ACME.Enabled = %v, want true", config.TLS.ACME.Enabled)
	}
	if config.TLS.ACME.Email != "admin@example.com" {
		t.Errorf("TLS.ACME.Email = %v, want admin@example.com", config.TLS.ACME.Email)
	}

	// Verify compression config
	if !config.Compress.Enabled {
		t.Errorf("Compress.Enabled = %v, want true", config.Compress.Enabled)
	}
	if config.Compress.Level != 8 {
		t.Errorf("Compress.Level = %v, want 8", config.Compress.Level)
	}

	// Verify redirect config
	if !config.Redirect.HTTPS {
		t.Errorf("Redirect.HTTPS = %v, want true", config.Redirect.HTTPS)
	}
	if config.Redirect.Code != 308 {
		t.Errorf("Redirect.Code = %v, want 308", config.Redirect.Code)
	}

	// Verify healthcheck config
	if config.Healthcheck.LivenessEndpoint != "/custom/live" {
		t.Errorf("Healthcheck.LivenessEndpoint = %v, want /custom/live", config.Healthcheck.LivenessEndpoint)
	}
	if config.Healthcheck.ReadinessEndpoint != "/custom/ready" {
		t.Errorf("Healthcheck.ReadinessEndpoint = %v, want /custom/ready", config.Healthcheck.ReadinessEndpoint)
	}
	if config.Healthcheck.StartupEndpoint != "/custom/startup" {
		t.Errorf("Healthcheck.StartupEndpoint = %v, want /custom/startup", config.Healthcheck.StartupEndpoint)
	}

	// Verify prometheus config
	if !config.Prometheus.Enabled {
		t.Errorf("Prometheus.Enabled = %v, want true", config.Prometheus.Enabled)
	}
	if config.Prometheus.Path != "/custom/metrics" {
		t.Errorf("Prometheus.Path = %v, want /custom/metrics", config.Prometheus.Path)
	}
}

func TestDefaultConfig(t *testing.T) {
	if DefaultConfig == nil {
		t.Fatal("DefaultConfig is nil")
	}

	if DefaultConfig.Name != "app" {
		t.Errorf("DefaultConfig.Name = %v, want app", DefaultConfig.Name)
	}
	if DefaultConfig.Version != "1.0.0" {
		t.Errorf("DefaultConfig.Version = %v, want 1.0.0", DefaultConfig.Version)
	}
	if DefaultConfig.GracefulTimeout != 30*time.Second {
		t.Errorf("DefaultConfig.GracefulTimeout = %v, want 30s", DefaultConfig.GracefulTimeout)
	}

	// Test HTTP defaults
	if DefaultConfig.HTTP.BindAddr != "localhost:8080" {
		t.Errorf("DefaultConfig.HTTP.BindAddr = %v, want localhost:8080", DefaultConfig.HTTP.BindAddr)
	}
	if DefaultConfig.HTTP.MaxHeaderBytes != 1<<20 {
		t.Errorf("DefaultConfig.HTTP.MaxHeaderBytes = %v, want %v", DefaultConfig.HTTP.MaxHeaderBytes, 1<<20)
	}

	// Test TLS defaults
	if DefaultConfig.TLS.Enabled != false {
		t.Errorf("DefaultConfig.TLS.Enabled = %v, want false", DefaultConfig.TLS.Enabled)
	}
	if DefaultConfig.TLS.BindAddr != "localhost:8443" {
		t.Errorf("DefaultConfig.TLS.BindAddr = %v, want localhost:8443", DefaultConfig.TLS.BindAddr)
	}

	// Test compression defaults
	if DefaultConfig.Compress.Enabled != false {
		t.Errorf("DefaultConfig.Compress.Enabled = %v, want false", DefaultConfig.Compress.Enabled)
	}
	if DefaultConfig.Compress.Level != 6 {
		t.Errorf("DefaultConfig.Compress.Level = %v, want 6", DefaultConfig.Compress.Level)
	}

	// Test redirect defaults
	if DefaultConfig.Redirect.HTTPS != false {
		t.Errorf("DefaultConfig.Redirect.HTTPS = %v, want false", DefaultConfig.Redirect.HTTPS)
	}
	if DefaultConfig.Redirect.Code != http.StatusMovedPermanently {
		t.Errorf("DefaultConfig.Redirect.Code = %v, want %v", DefaultConfig.Redirect.Code, http.StatusMovedPermanently)
	}

	// Test healthcheck defaults
	if DefaultConfig.Healthcheck.LivenessEndpoint != "/livez" {
		t.Errorf("DefaultConfig.Healthcheck.LivenessEndpoint = %v, want /livez", DefaultConfig.Healthcheck.LivenessEndpoint)
	}
	if DefaultConfig.Healthcheck.ReadinessEndpoint != "/readyz" {
		t.Errorf("DefaultConfig.Healthcheck.ReadinessEndpoint = %v, want /readyz", DefaultConfig.Healthcheck.ReadinessEndpoint)
	}
	if DefaultConfig.Healthcheck.StartupEndpoint != "/startupz" {
		t.Errorf("DefaultConfig.Healthcheck.StartupEndpoint = %v, want /startupz", DefaultConfig.Healthcheck.StartupEndpoint)
	}
	if DefaultConfig.Healthcheck.LivenessHandler == nil {
		t.Error("DefaultConfig.Healthcheck.LivenessHandler is nil")
	}
	if DefaultConfig.Healthcheck.ReadinessHandler == nil {
		t.Error("DefaultConfig.Healthcheck.ReadinessHandler is nil")
	}
	if DefaultConfig.Healthcheck.StartupHandler == nil {
		t.Error("DefaultConfig.Healthcheck.StartupHandler is nil")
	}

	// Test prometheus defaults
	if DefaultConfig.Prometheus.Enabled != false {
		t.Errorf("DefaultConfig.Prometheus.Enabled = %v, want false", DefaultConfig.Prometheus.Enabled)
	}
	if DefaultConfig.Prometheus.Path != "/metrics" {
		t.Errorf("DefaultConfig.Prometheus.Path = %v, want /metrics", DefaultConfig.Prometheus.Path)
	}
}

func TestConfig_FlagSet_DefaultValues(t *testing.T) {
	config := &Config{
		Name:            "customapp",
		Version:         "2.1.0",
		GracefulTimeout: 45 * time.Second,
		HTTP: HTTPConfig{
			BindAddr:       ":9090",
			IdleTimeout:    90 * time.Second,
			MaxHeaderBytes: 2 << 20,
		},
		TLS: TLSConfig{
			Enabled:  true,
			BindAddr: ":9443",
		},
		Healthcheck: HealthcheckConfig{
			LivenessEndpoint:  "/custom/live",
			ReadinessEndpoint: "/custom/ready",
			StartupEndpoint:   "/custom/startup",
		},
		Prometheus: PrometheusConfig{
			Enabled: true,
			Path:    "/custom/prometheus",
		},
	}

	fs := config.FlagSet()

	nameFlag := fs.Lookup(ServerName)
	if nameFlag == nil {
		t.Fatal("Name flag not found")
	}
	if nameFlag.DefValue != "customapp" {
		t.Errorf("Name flag default = %v, want customapp", nameFlag.DefValue)
	}

	versionFlag := fs.Lookup(ServerVersion)
	if versionFlag == nil {
		t.Fatal("Version flag not found")
	}
	if versionFlag.DefValue != "2.1.0" {
		t.Errorf("Version flag default = %v, want 2.1.0", versionFlag.DefValue)
	}

	httpBindAddrFlag := fs.Lookup(ServerHTTPBindAddr)
	if httpBindAddrFlag == nil {
		t.Fatal("HTTP BindAddr flag not found")
	}
	if httpBindAddrFlag.DefValue != ":9090" {
		t.Errorf("HTTP BindAddr flag default = %v, want :9090", httpBindAddrFlag.DefValue)
	}

	tlsEnabledFlag := fs.Lookup(ServerTLSEnabled)
	if tlsEnabledFlag == nil {
		t.Fatal("TLS Enabled flag not found")
	}
	if tlsEnabledFlag.DefValue != "true" {
		t.Errorf("TLS Enabled flag default = %v, want true", tlsEnabledFlag.DefValue)
	}

	livenessFlag := fs.Lookup(ServerHealthcheckLivenessEndpoint)
	if livenessFlag == nil {
		t.Fatal("Healthcheck Liveness flag not found")
	}
	if livenessFlag.DefValue != "/custom/live" {
		t.Errorf("Healthcheck Liveness flag default = %v, want /custom/live", livenessFlag.DefValue)
	}

	prometheusEnabledFlag := fs.Lookup(ServerPrometheusEnabled)
	if prometheusEnabledFlag == nil {
		t.Fatal("Prometheus Enabled flag not found")
	}
	if prometheusEnabledFlag.DefValue != "true" {
		t.Errorf("Prometheus Enabled flag default = %v, want true", prometheusEnabledFlag.DefValue)
	}

	prometheusPathFlag := fs.Lookup(ServerPrometheusPath)
	if prometheusPathFlag == nil {
		t.Fatal("Prometheus Path flag not found")
	}
	if prometheusPathFlag.DefValue != "/custom/prometheus" {
		t.Errorf("Prometheus Path flag default = %v, want /custom/prometheus", prometheusPathFlag.DefValue)
	}
}

func TestConfig_FlagSet_EmptyParse(t *testing.T) {
	config := DefaultConfig

	fs := config.FlagSet()

	var args []string

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse empty flags: %v", err)
	}

	// Should retain default values
	if config.Name != "app" {
		t.Errorf("Name = %v, want app (default)", config.Name)
	}
	if config.TLS.Enabled != false {
		t.Errorf("TLS.Enabled = %v, want false (default)", config.TLS.Enabled)
	}
	if config.Healthcheck.LivenessEndpoint != "/livez" {
		t.Errorf("Healthcheck.LivenessEndpoint = %v, want /livez (default)", config.Healthcheck.LivenessEndpoint)
	}
	if config.Prometheus.Enabled != false {
		t.Errorf("Prometheus.Enabled = %v, want false (default)", config.Prometheus.Enabled)
	}
	if config.Prometheus.Path != "/metrics" {
		t.Errorf("Prometheus.Path = %v, want /metrics (default)", config.Prometheus.Path)
	}
}
