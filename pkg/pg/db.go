package pg

import (
	"fmt"
	"strings"

	"github.com/Armatorix/BitLyke/pkg/config"
	"github.com/avast/retry-go"
	"github.com/go-pg/pg/v9"
)

const (
	unix = "unix"
	tcp  = "tcp"
)

type DB struct {
	*pg.DB
}

func getNetwork(addr string) string {
	if strings.Contains(addr, "/") {
		return unix
	}
	return tcp
}

func New(cfg config.PostgresConfig) (*DB, error) {
	db := &DB{
		pg.Connect(&pg.Options{
			Network:  getNetwork(cfg.Address),
			Addr:     cfg.Address,
			User:     cfg.User,
			Password: cfg.Password,
			Database: cfg.Database,
		})}

	err := retry.Do(db.TestRequest)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) TestRequest() error {
	const testNum = 1
	var num int
	_, err := db.Query(pg.Scan(&num), "SELECT ?", testNum)
	if err != nil {
		return fmt.Errorf("connection check failed: %w", err)
	}
	if num != testNum {
		return fmt.Errorf("connection check failed: should have %d, was %d", testNum, num)
	}
	return nil
}
