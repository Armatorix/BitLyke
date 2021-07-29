package pg

import (
	"strings"

	"github.com/Armatorix/BitLyke/pkg/config"
	"github.com/avast/retry-go"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
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
	opts, err := pg.ParseURL(cfg.URI)
	if err != nil {
		return nil, err
	}
	db := &DB{pg.Connect(opts)}
	err = retry.Do(db.TestRequest)
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
		return errors.Wrap(err, "connection check failed")
	}
	if num != testNum {
		return errors.Wrap(err, "connection check failed: different value")
	}
	return nil
}
