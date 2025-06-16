package middleware

import (
	"testing"
	"time"
)

func TestRateLimiter_FlagSet(t *testing.T) {
	config := &RateLimiter{
		Enabled: true,
		Store:   LimiterStoreMemory,
		Memory: RateLimiterMemoryStore{
			Rate:      10.5,
			Burst:     100,
			ExpiresIn: 5 * time.Minute,
		},
	}

	fs := config.FlagSet()

	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	enabledFlag := fs.Lookup(RateLimiterEnabled)
	if enabledFlag == nil {
		t.Errorf("Flag %s not found", RateLimiterEnabled)
	} else {
		if enabledFlag.DefValue != "true" {
			t.Errorf("Flag %s default value = %v, want true", RateLimiterEnabled, enabledFlag.DefValue)
		}
	}

	storeFlag := fs.Lookup(RateLimiterStoreType)
	if storeFlag == nil {
		t.Errorf("Flag %s not found", RateLimiterStoreType)
	} else {
		if storeFlag.DefValue != "memory" {
			t.Errorf("Flag %s default value = %v, want memory", RateLimiterStoreType, storeFlag.DefValue)
		}
	}

	rateFlag := fs.Lookup(RateLimiterMemoryRate)
	if rateFlag == nil {
		t.Errorf("Flag %s not found", RateLimiterMemoryRate)
	} else {
		if rateFlag.DefValue != "10.5" {
			t.Errorf("Flag %s default value = %v, want 10.5", RateLimiterMemoryRate, rateFlag.DefValue)
		}
	}

	burstFlag := fs.Lookup(RateLimiterMemoryBurst)
	if burstFlag == nil {
		t.Errorf("Flag %s not found", RateLimiterMemoryBurst)
	} else {
		if burstFlag.DefValue != "100" {
			t.Errorf("Flag %s default value = %v, want 100", RateLimiterMemoryBurst, burstFlag.DefValue)
		}
	}

	expiresFlag := fs.Lookup(RateLimiterMemoryExpires)
	if expiresFlag == nil {
		t.Errorf("Flag %s not found", RateLimiterMemoryExpires)
	} else {
		if expiresFlag.DefValue != "5m0s" {
			t.Errorf("Flag %s default value = %v, want 5m0s", RateLimiterMemoryExpires, expiresFlag.DefValue)
		}
	}
}

func TestRateLimiter_FlagSet_Parse(t *testing.T) {
	config := &RateLimiter{
		Enabled: false,
		Store:   LimiterStoreMemory,
		Memory: RateLimiterMemoryStore{
			Rate:      0,
			Burst:     0,
			ExpiresIn: 3 * time.Minute,
		},
	}

	fs := config.FlagSet()

	args := []string{
		"--rate-limiter-enabled",
		"--rate-limiter-store", "memory",
		"--rate-limiter-memory-rate", "25.5",
		"--rate-limiter-memory-burst", "50",
		"--rate-limiter-memory-expires", "10m",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	if !config.Enabled {
		t.Errorf("Enabled = %v, want true", config.Enabled)
	}
	if config.Store != LimiterStoreMemory {
		t.Errorf("Store = %v, want %v", config.Store, LimiterStoreMemory)
	}
	if config.Memory.Rate != 25.5 {
		t.Errorf("Memory.Rate = %v, want 25.5", config.Memory.Rate)
	}
	if config.Memory.Burst != 50 {
		t.Errorf("Memory.Burst = %v, want 50", config.Memory.Burst)
	}
	if config.Memory.ExpiresIn != 10*time.Minute {
		t.Errorf("Memory.ExpiresIn = %v, want 10m", config.Memory.ExpiresIn)
	}
}

func TestDefaultRateLimiter(t *testing.T) {
	if DefaultRateLimiter == nil {
		t.Fatal("DefaultRateLimiter is nil")
	}

	if DefaultRateLimiter.Enabled != false {
		t.Errorf("DefaultRateLimiter.Enabled = %v, want false", DefaultRateLimiter.Enabled)
	}
	if DefaultRateLimiter.Store != LimiterStoreMemory {
		t.Errorf("DefaultRateLimiter.Store = %v, want %v", DefaultRateLimiter.Store, LimiterStoreMemory)
	}
	if DefaultRateLimiter.Memory.Rate != 0 {
		t.Errorf("DefaultRateLimiter.Memory.Rate = %v, want 0", DefaultRateLimiter.Memory.Rate)
	}
	if DefaultRateLimiter.Memory.Burst != 0 {
		t.Errorf("DefaultRateLimiter.Memory.Burst = %v, want 0", DefaultRateLimiter.Memory.Burst)
	}
	if DefaultRateLimiter.Memory.ExpiresIn != 3*time.Minute {
		t.Errorf("DefaultRateLimiter.Memory.ExpiresIn = %v, want 3m", DefaultRateLimiter.Memory.ExpiresIn)
	}
}

func TestRateLimiter_FlagSet_DefaultValues(t *testing.T) {
	config := &RateLimiter{
		Enabled: true,
		Store:   LimiterStoreMemory,
		Memory: RateLimiterMemoryStore{
			Rate:      5.0,
			Burst:     20,
			ExpiresIn: 1 * time.Hour,
		},
	}

	fs := config.FlagSet()

	enabledFlag := fs.Lookup(RateLimiterEnabled)
	if enabledFlag == nil {
		t.Fatal("Enabled flag not found")
	}
	if enabledFlag.DefValue != "true" {
		t.Errorf("Enabled flag default = %v, want true", enabledFlag.DefValue)
	}

	storeFlag := fs.Lookup(RateLimiterStoreType)
	if storeFlag == nil {
		t.Fatal("Store flag not found")
	}
	if storeFlag.DefValue != "memory" {
		t.Errorf("Store flag default = %v, want memory", storeFlag.DefValue)
	}

	rateFlag := fs.Lookup(RateLimiterMemoryRate)
	if rateFlag == nil {
		t.Fatal("Rate flag not found")
	}
	if rateFlag.DefValue != "5" {
		t.Errorf("Rate flag default = %v, want 5", rateFlag.DefValue)
	}
}

func TestRateLimiter_FlagSet_DisabledByDefault(t *testing.T) {
	config := DefaultRateLimiter

	fs := config.FlagSet()

	var args []string

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse empty flags: %v", err)
	}

	if config.Enabled {
		t.Errorf("Enabled = %v, want false (default)", config.Enabled)
	}
}

func TestRateLimiterStore_String(t *testing.T) {
	tests := []struct {
		store RateLimiterStore
		want  string
	}{
		{LimiterStoreMemory, "memory"},
		{RateLimiterStore("invalid"), "unknown store: invalid"},
	}

	for _, tt := range tests {
		got := tt.store.String()
		if got != tt.want {
			t.Errorf("RateLimiterStore(%s).String() = %v, want %v", tt.store, got, tt.want)
		}
	}
}

func TestRateLimiterStore_Set(t *testing.T) {
	tests := []struct {
		value   string
		want    RateLimiterStore
		wantErr bool
	}{
		{"memory", LimiterStoreMemory, false},
		{"MEMORY", LimiterStoreMemory, false},
		{"invalid", RateLimiterStore(""), true},
	}

	for _, tt := range tests {
		var store RateLimiterStore
		err := store.Set(tt.value)
		if tt.wantErr {
			if err == nil {
				t.Errorf("RateLimiterStore.Set(%q) expected error but got nil", tt.value)
			}
		} else {
			if err != nil {
				t.Errorf("RateLimiterStore.Set(%q) unexpected error: %v", tt.value, err)
			}
			if store != tt.want {
				t.Errorf("RateLimiterStore.Set(%q) = %v, want %v", tt.value, store, tt.want)
			}
		}
	}
}

func TestRateLimiterStore_Type(t *testing.T) {
	var store RateLimiterStore
	if got := store.Type(); got != "string" {
		t.Errorf("RateLimiterStore.Type() = %v, want string", got)
	}
}

func TestNewRateLimiter(t *testing.T) {
	config := &RateLimiter{
		Enabled: true,
		Store:   LimiterStoreMemory,
		Memory: RateLimiterMemoryStore{
			Rate:      10,
			Burst:     20,
			ExpiresIn: 5 * time.Minute,
		},
	}

	middleware := NewRateLimiter(config)
	if middleware == nil {
		t.Fatal("NewRateLimiter() returned nil")
	}
}

func TestNewRateLimiter_DefaultConfig(t *testing.T) {
	middleware := NewRateLimiter(DefaultRateLimiter)
	if middleware == nil {
		t.Fatal("NewRateLimiter() with DefaultRateLimiter returned nil")
	}
}

func TestNewRateLimiter_InvalidStore(t *testing.T) {
	config := &RateLimiter{
		Enabled: true,
		Store:   RateLimiterStore("invalid"),
	}

	middleware := NewRateLimiter(config)
	if middleware != nil {
		t.Error("NewRateLimiter() with invalid store should return nil")
	}
}
