package config

import "time"

type GRPCClientConfig struct {
	ConnectTimeout    time.Duration `env:"GRPC_CLIENT_CONNECT_TIMEOUT" validate:"gte=0"`
	MaxBackoffDelay   time.Duration `env:"GRPC_CLIENT_MAX_BACKOFF_DELAY" validate:"gte=0"`
	BaseBackoffDelay  time.Duration `env:"GRPC_CLIENT_BASE_BACKOFF_DELAY" validate:"gte=0"`
	BackoffMultiplier float64       `env:"GRPC_CLIENT_BACKOFF_MULTIPLIER" validate:"gte=1"`
	BackoffJitter     float64       `env:"GRPC_CLIENT_BACKOFF_JITTER" validate:"gte=0"`
}
