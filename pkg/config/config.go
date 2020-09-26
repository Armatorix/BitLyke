package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

type ServerConfig struct {
	Address  string  `env:"SERVER_ADDRESS"`
	LogLevel log.Lvl `env:"SERVER_LOG_LEVEL"`
}

type PostgresConfig struct {
	Address  string `env:"PSQL_ADDRESS"`
	Database string `env:"PSQL_DATABASE"`
	User     string `env:"PSQL_USER"`
	Password string `env:"PSQL_PASSWORD"`
}

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
}

func New() *Config {
	return &Config{
		Server: ServerConfig{
			Address: ":8080",
		},
		Postgres: PostgresConfig{
			Address:  "localhost:5432",
			Database: "bitlyke",
			User:     "postgres",
			Password: "example",
		},
	}
}

func FromEnv() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(err, "env parsing: %w")
	}
	return cfg, nil
}
