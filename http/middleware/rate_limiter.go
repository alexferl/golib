package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
	"golang.org/x/time/rate"
)

// RateLimiterMemoryStore holds configuration for memory-based rate limiting.
type RateLimiterMemoryStore struct {
	// Rate specifies the rate limit per second.
	// Optional. Default value 0.
	Rate float64

	// Burst specifies the burst limit.
	// Optional. Default value 0.
	Burst int

	// ExpiresIn specifies the expiration time for rate limit entries.
	// Optional. Default value 3 minutes.
	ExpiresIn time.Duration
}

type RateLimiterStore string

const (
	limiterStoreMemory = "memory"
)

const (
	LimiterStoreMemory RateLimiterStore = limiterStoreMemory
)

var LimiterStores = []string{limiterStoreMemory}

func (s *RateLimiterStore) String() string {
	switch *s {
	case LimiterStoreMemory:
		return limiterStoreMemory
	default:
		return fmt.Sprintf("unknown store: %s", *s)
	}
}

func (s *RateLimiterStore) Set(value string) error {
	switch strings.ToLower(value) {
	case limiterStoreMemory:
		*s = LimiterStoreMemory
		return nil
	default:
		return fmt.Errorf("invalid rate limiter store: %s (must be one of: %s)", value, strings.Join(LimiterStores, ", "))
	}
}

func (s *RateLimiterStore) Type() string {
	return "string"
}

// RateLimiter holds configuration for rate limiting middleware.
type RateLimiter struct {
	// Enabled indicates whether rate limiting middleware is enabled.
	// Optional. Default value false.
	Enabled bool

	// Store specifies the rate limiter store type.
	// Optional. Default value "memory".
	Store RateLimiterStore

	// Memory holds memory store configuration.
	// Optional. Default value with 3 minute expiration.
	Memory RateLimiterMemoryStore
}

// DefaultRateLimiter provides default RateLimiter configuration.
var DefaultRateLimiter = &RateLimiter{
	Enabled: false,
	Store:   LimiterStoreMemory,
	Memory: RateLimiterMemoryStore{
		Rate:      0,
		Burst:     0,
		ExpiresIn: 3 * time.Minute,
	},
}

const (
	RateLimiterEnabled       = "rate-limiter-enabled"
	RateLimiterStoreType     = "rate-limiter-store"
	RateLimiterMemoryRate    = "rate-limiter-memory-rate"
	RateLimiterMemoryBurst   = "rate-limiter-memory-burst"
	RateLimiterMemoryExpires = "rate-limiter-memory-expires"
)

// FlagSet returns a pflag.FlagSet for CLI configuration.
func (r *RateLimiter) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("Rate Limiter", pflag.ExitOnError)
	fs.BoolVar(&r.Enabled, RateLimiterEnabled, r.Enabled, "Enable rate limiting middleware")
	fs.Var(&r.Store, RateLimiterStoreType, fmt.Sprintf("Rate limiter store type\nValues: %s", strings.Join(LimiterStores, ", ")))
	fs.Float64Var(&r.Memory.Rate, RateLimiterMemoryRate, r.Memory.Rate, "Rate limit per second for memory store")
	fs.IntVar(&r.Memory.Burst, RateLimiterMemoryBurst, r.Memory.Burst, "Burst limit for memory store")
	fs.DurationVar(&r.Memory.ExpiresIn, RateLimiterMemoryExpires, r.Memory.ExpiresIn, "Expiration time for rate limit entries")
	return fs
}

// NewRateLimiter creates a new rate limiter middleware with the given configuration.
func NewRateLimiter(config *RateLimiter) echo.MiddlewareFunc {
	switch config.Store {
	case LimiterStoreMemory:
		s := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(config.Memory.Rate),
			Burst:     config.Memory.Burst,
			ExpiresIn: config.Memory.ExpiresIn,
		})

		return middleware.RateLimiter(s)
	}
	return nil
}
