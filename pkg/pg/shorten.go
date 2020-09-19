package db

import (
	"errors"
	"fmt"

	"github.com/Armatorix/BitLyke/pkg/model"
	"github.com/go-pg/pg/v9"
)

var ErrNotFound = errors.New("not found")

func (db *DB) InsertLinkShorten(l *model.ShortenLink) (*model.ShortenLink, error) {
	err := db.Insert(l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (db *DB) GetDestinationLink(shorten string) (*model.ShortenLink, error) {
	l := &model.ShortenLink{
		ShortenPath: shorten,
	}
	err := db.Model(&l).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("unexpected select failure: %w", err)
	}
	return l, nil
}

func (db *DB) GetLinkShortens() ([]model.ShortenLink, error) {
	ls := []model.ShortenLink{}
	err := db.Model(&ls).Select()
	if err != nil {
		return nil, err
	}
	return ls, nil
}
