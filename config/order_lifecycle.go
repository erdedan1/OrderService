package config

import "time"

type OrderLifecycleConfig struct {
	MaxDuration time.Duration `env:"ORDER_LIFECYCLE_MAX_DURATION" env-default:"1m" validate:"gt=0"`
}
