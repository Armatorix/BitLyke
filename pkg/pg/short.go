package pg

import (
	"github.com/Armatorix/BitLyke/pkg/schema"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrDuplicatedEntry = errors.New("duplicated entry")
)

func isDuplicatedKeyErr(err error) bool {
	pgerr, ok := err.(pg.Error)
	return !ok || pgerr.IntegrityViolation()
}

func (db *DB) InsertShort(l *schema.ShortLink) (*schema.ShortLink, error) {
	_, err := db.Model(l).Insert()
	if err != nil {
		if isDuplicatedKeyErr(err) {
			return nil, ErrDuplicatedEntry
		}
		return nil, err
	}
	return l, nil
}

func (db *DB) GetDestinationLink(short string) (*schema.ShortLink, error) {
	l := &schema.ShortLink{
		ShortPath: short,
	}
	err := db.Model(l).Where("short_path = ?", l.ShortPath).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "unexpected select failure")
	}
	return l, nil
}

func (db *DB) GetLinkShorts() ([]schema.ShortLink, error) {
	ls := []schema.ShortLink{}
	err := db.Model(&ls).Select()
	if err != nil {
		return nil, err
	}
	return ls, nil
}

func (db *DB) DeleteShort(short string) error {
	res, err := db.Model((*schema.ShortLink)(nil)).Where("\"short_path\" = ?", short).Delete()
	if err != nil {
		return errors.Wrap(err, "unexpected select failure")
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
