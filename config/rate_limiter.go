package config

type RateLimiterConfig struct {
	RequestsPerSecond float64 `env:"RATE_LIMITER_RPS" env-default:"50" validate:"gt=0"`
	Burst             int     `env:"RATE_LIMITER_BURST" env-default:"100" validate:"gt=0"`
}
