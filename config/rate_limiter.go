package config

type RateLimiterConfig struct {
	GlobalRequestsPerSecond float64 `env:"RATE_LIMITER_GLOBAL_RPS" env-default:"50" validate:"gt=0"`
	GlobalBurst             float64 `env:"RATE_LIMITER_GLOBAL_BURST" env-default:"100" validate:"gt=0"`

	ClientRequestsPerSecond float64 `env:"RATE_LIMITER_CLIENT_RPS" env-default:"50" validate:"gt=0"`
	ClientBurst             float64 `env:"RATE_LIMITER_CLIENT_BURST" env-default:"100" validate:"gt=0"`
}
