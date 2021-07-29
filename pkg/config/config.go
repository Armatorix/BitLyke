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
	URI string `env:"DATABASE"`
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
			URI: "postgres://postgres:example@localhost:5432/bitlyke",
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
