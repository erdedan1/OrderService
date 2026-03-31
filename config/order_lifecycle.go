package config

import "time"

type OrderLifecircuitConfig struct {
	StepInterval time.Duration `env:"ORDER_LIFECIRCUIT_STEP_INTERVAL" env-default:"5s" validate:"gt=0"`
	TimeOut      time.Duration `env:"ORDER_LIFECIRCUIT_TIMEOUT" env-default:"12h" validate:"gt=0"`
}
