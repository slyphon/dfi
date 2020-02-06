package dotfile

import (
	"github.com/pkg/errors"
	str "strings"
)

type OnConflict int

const (
	Backup OnConflict = iota
	Replace
	Warn
	Fail
)

func OnConflictForString(s string) (OnConflict, error) {
	switch str.ToLower(s) {
	case "backup":
		return Backup, nil
	case "replace":
		return Replace, nil
	case "warn":
		return Warn, nil
	case "fail":
		return Fail, nil
	default:
		return -1, errors.Errorf("invalid OnConflict string: %v", s)
	}
}
