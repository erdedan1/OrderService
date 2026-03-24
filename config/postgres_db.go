package config

type PostgresDB struct {
	Password string `env:"DB_PASSWORD" validate:"required"`
	User     string `env:"DB_USER" "validate:"required"`
	Host     string `env:"DB_HOST" "validate:"required"`
	Port     string `env:"DB_PORT" "validate:"required"`
	Name     string `env:"DB_NAME" "validate:"required"`
}
