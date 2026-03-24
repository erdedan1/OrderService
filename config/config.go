package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
)

var Global Config

type Config struct {
	GRPCServer     GRPCServerConfig     `validate:"required"`
	GRPCApi        GRPCApiConfig        `validate:"required"`
	GRPCClient     GRPCClientConfig     `validate:"required"`
	Infrastructure InfrastructureConfig `validate:"required"`
	PostgresDB     PostgresDB
}

type InfrastructureConfig struct {
	RedisConfig   RedisConfig   `validate:"required"`
	Observability Observability `validate:"required"`
}

type GRPCApiConfig struct {
	SpotInstrumentServiceHost string `env:"GRPC_API_SPOT_INSTRUMENT_SERVICE_HOST" validate:"required"`
}

func New() error {
	if err := env.Parse(Global); err != nil {
		return err
	}

	v := validator.New()
	if err := v.Struct(Global); err != nil {
		return err
	}

	return nil
}
