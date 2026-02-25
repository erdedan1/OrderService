package config

import (
	tel "github.com/erdedan1/shared/telemetry"
)

type Observability struct {
	Loglvl    string    `env:"LOG_LVL" validate:"required"`
	Telemetry Telemetry `validate:"required"`
}

type Telemetry struct {
	ServiceName    string `env:"TELEMETRY_SERVICE_NAME" validate:"required"`
	ServiceVersion string `env:"TELEMETRY_SERVICE_VERSION" validate:"required"`
	Environment    string `env:"TELEMETRY_ENVIROMENT" validate:"required"`

	Host     string  `env:"TELEMETRY_HOST" validate:"required"`
	Port     string  `env:"TELEMETRY_PORT" validate:"required"`
	Enabled  bool    `env:"TELEMETRY_ENABLED" validate:"required"`
	Sampling float64 `env:"TELEMETRY_SAMPLING" validate:"required"`
}

func NewTelemetryConfig(t Telemetry) tel.Config {
	return tel.Config{
		ServiceName:    t.ServiceName,
		ServiceVersion: t.ServiceVersion,
		Environment:    t.Environment,
		Host:           t.Host,
		Port:           t.Port,
		Enabled:        t.Enabled,
		Sampling:       t.Sampling,
	}
}
