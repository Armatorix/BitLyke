package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Address string `env:"SERVER_ADDRESS"`
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
			Address: ":8081",
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
		return nil, fmt.Errorf("env parsing: %w", err)
	}
	return cfg, nil
}
