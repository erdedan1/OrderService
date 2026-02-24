package config

type PostgresDB struct {
	Password string `env:"DB_PASSWORD"`
	User     string `env:"DB_USER"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Name     string `env:"DB_NAME"`
}
