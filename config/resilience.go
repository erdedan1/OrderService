package config

import "time"

type ResilienceConfig struct {
	RateLimiter    RateLimiterConfig    `validate:"required"`
	CircuitBreaker CircuitBreakerConfig `validate:"required"`
}

type RateLimiterConfig struct {
	RequestsPerSecond float64 `env:"RATE_LIMITER_RPS" env-default:"50" validate:"gt=0"`
	Burst             int     `env:"RATE_LIMITER_BURST" env-default:"100" validate:"gt=0"`
}

type CircuitBreakerConfig struct {
	ConsecutiveFailures uint32        `env:"CIRCUIT_BREAKER_CONSECUTIVE_FAILURES" env-default:"5" validate:"gt=0"`
	HalfOpenRequests    uint32        `env:"CIRCUIT_BREAKER_HALF_OPEN_REQUESTS" env-default:"3" validate:"gt=0"`
	OpenTimeout         time.Duration `env:"CIRCUIT_BREAKER_OPEN_TIMEOUT" env-default:"10s" validate:"gt=0"`
}
