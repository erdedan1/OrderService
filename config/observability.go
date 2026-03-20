package config

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
	Sampling float64 `env:"TELEMETRY_SAMPLING" validate:"required"`
}
