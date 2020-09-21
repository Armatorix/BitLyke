package xgopg

import "strings"

const duplicatedKeyErrorPrefix = "ERROR #23505 duplicate key"

func IsDuplicatedKeyErr(err error) bool {
	return strings.HasPrefix(err.Error(), duplicatedKeyErrorPrefix)
}
