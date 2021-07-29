package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

const DefaultPort = 8080

type ServerConfig struct {
	Port     int     `env:"PORT"`
	LogLevel log.Lvl `env:"SERVER_LOG_LEVEL"`
}

type PostgresConfig struct {
	URI string `env:"DATABASE_URL"`
}

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
}

func New() *Config {
	return &Config{
		Server: ServerConfig{
			Port: DefaultPort,
		},
		Postgres: PostgresConfig{
			URI: "postgres://postgres:example@localhost:5432/bitlyke?sslmode=disable",
		},
	}
}

func FromEnv() (*Config, error) {
	cfg := New()
	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(err, "env parsing: %w")
	}
	return cfg, nil
}
