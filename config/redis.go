package config

import "time"

type RedisConfig struct {
	Host         string        `env:"REDIS_HOST" validate:"required,hostname|ip"`
	Port         string        `env:"REDIS_PORT" validate:"required,numeric"`
	Password     string        `env:"REDIS_PASSWORD" validate:"-"`
	DB           int           `env:"REDIS_DB" validate:"gte=0"`
	MinIdleConns int           `env:"REDIS_MIN_IDLE_CONNS" validate:"gte=0"`
	PoolSize     int           `env:"REDIS_POOL_SIZE" validate:"gte=0"`
	PoolTimeout  time.Duration `env:"REDIS_POOL_TIMEOUT" validate:"gte=0"`
}
