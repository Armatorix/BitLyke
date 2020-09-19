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

func (db *DB) InsertShort(l *model.ShortLink) (*model.ShortLink, error) {
	_, err := db.GetDestinationLink(l.ShortPath)
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

func (db *DB) GetDestinationLink(short string) (*model.ShortLink, error) {
	l := &model.ShortLink{
		ShortPath: short,
	}
	err := db.Model(l).Where("short_path = ?", l.ShortPath).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("unexpected select failure: %w", err)
	}
	return l, nil
}

func (db *DB) GetLinkShorts() ([]model.ShortLink, error) {
	ls := []model.ShortLink{}
	err := db.Model(&ls).Select()
	if err != nil {
		return nil, err
	}
	return ls, nil
}

func (db *DB) DeleteShort(short string) (*model.ShortLink, error) {
	l := &model.ShortLink{}
	_, err := db.Model(l).Where("short_path = ?", l.ShortPath).Delete()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("unexpected select failure: %w", err)
	}
	return l, nil
}
