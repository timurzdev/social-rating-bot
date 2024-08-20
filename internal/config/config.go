package config

import "fmt"

type Config struct {
	Token string `env:"TOKEN"`
	//DSN              string `env:"DSN"`
	MigrationsPath   string `env:"MIGRATIONS_PATH"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresPort     int    `env:"POSTGRES_PORT"`
	PostgresDBName   string `env:"POSTGRES_DB_NAME"`
}

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@postgres:%d/%s?sslmode=disable",
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresPort,
		c.PostgresDBName,
	)
}
