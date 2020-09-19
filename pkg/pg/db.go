package db

import (
	"strings"

	"github.com/Armatorix/BitLyke/pkg/config"
	"github.com/go-pg/pg/v9"
)

type DB struct {
	*pg.DB
}

func getNetwork(addr string) string {
	if strings.Contains(addr, "/") {
		return "unix"
	}
	return "tcp"
}

func New(cfg config.PostgresConfig) *DB {
	return &DB{
		pg.Connect(&pg.Options{
			Network:  getNetwork(cfg.Address),
			Addr:     cfg.Address,
			User:     cfg.User,
			Password: cfg.Password,
			Database: cfg.Database,
		})}
}
