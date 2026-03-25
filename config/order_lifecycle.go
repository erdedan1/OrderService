package config

import "time"

type OrderLifecircuitConfig struct {
	StepInterval time.Duration `env:"ORDER_LIFECIRCUIT_STEP_INTERVAL" env-default:"5s" validate:"gt=0"`
}
