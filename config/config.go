package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	GRPCServer     GRPCServerConfig     `validate:"required"`
	GRPCApi        GRPCApiConfig        `validate:"required"`
	GRPCClient     GRPCClientConfig     `validate:"required"`
	Infrastructure InfrastructureConfig `validate:"required"`
	PostgresDB     PostgresDB
}

type InfrastructureConfig struct {
	RedisConfig RedisConfig `validate:"required"`
	Log_LVL     string      `env:"LOG_LVL" validate:"required"`
}

type GRPCApiConfig struct {
	SpotInstrumentServiceHost string `env:"GRPC_API_SPOT_INSTRUMENT_SERVICE_HOST" validate:"required"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	v := validator.New() //todo мб поменять валидатор
	if err := v.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
