package config

import "time"

type OrderLifecycleConfig struct {
	StepInterval time.Duration `env:"ORDER_LIFECYCLE_STEP_INTERVAL" env-default:"5s" validate:"gt=0"`
}
