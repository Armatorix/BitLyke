package pg

import (
	"log"

	"github.com/Armatorix/BitLyke/pkg/config"
	"github.com/Armatorix/BitLyke/pkg/schema"
	"github.com/avast/retry-go"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
)

type DB struct {
	*pg.DB
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

func (db *DB) InitModels() error {
	return db.Model(&schema.ShortLink{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
}

func (db *DB) TestRequest() error {
	const testNum = 1
	var num int
	_, err := db.Query(pg.Scan(&num), "SELECT ?", testNum)
	if err != nil {
		log.Println("failed test", err)
		return errors.Wrap(err, "connection check failed")
	}
	if num != testNum {
		log.Println("bad respond", num)
		return errors.Wrap(err, "connection check failed: different value")
	}
	return nil
}
