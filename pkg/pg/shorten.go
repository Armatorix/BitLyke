package pg

import (
	"errors"
	"fmt"

	"github.com/Armatorix/BitLyke/pkg/model"
	"github.com/go-pg/pg/v9"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrAlreadyInUse = errors.New("already in use")
)

func (db *DB) InsertLinkShorten(l *model.ShortenLink) (*model.ShortenLink, error) {
	_, err := db.GetDestinationLink(l.ShortenPath)
	if err == nil {
		return nil, ErrAlreadyInUse
	}
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}
	err = db.Insert(l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (db *DB) GetDestinationLink(shorten string) (*model.ShortenLink, error) {
	l := &model.ShortenLink{
		ShortenPath: shorten,
	}
	err := db.Model(l).Where("shorten_path = ?", l.ShortenPath).Select()
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
