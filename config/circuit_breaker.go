package config

import "time"

type CircuitBreakerConfig struct {
	ConsecutiveFailures uint32        `env:"CIRCUIT_BREAKER_CONSECUTIVE_FAILURES" env-default:"5" validate:"gt=0"`
	HalfOpenRequests    uint32        `env:"CIRCUIT_BREAKER_HALF_OPEN_REQUESTS" env-default:"3" validate:"gt=0"`
	OpenTimeout         time.Duration `env:"CIRCUIT_BREAKER_OPEN_TIMEOUT" env-default:"10s" validate:"gt=0"`
}
