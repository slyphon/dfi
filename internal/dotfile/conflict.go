package dotfile

import (
	"fmt"
	"os"
	fp "path/filepath"
	str "strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type (
	OnConflict int
	ConflictHandler interface {
		Handle(linkPath string) (skip bool, err error)
	}
)

const (
	Rename OnConflict = iota
	Replace
	Warn
	Fail
)

var ConflictHandlers = struct {
	Rename OnConflict
	Replace OnConflict
	Warn OnConflict
	Fail OnConflict
} {Rename, Replace, Warn, Fail}

const (
	TimeFormat string = "20060102150405"
)

var _ ConflictHandler = OnConflict(0)

func timestamp() string {
	return time.Now().Format(TimeFormat)
}

func isSymlink(fm os.FileMode) bool {
	return fm&os.ModeSymlink != 0
}

func nameForMode(fi os.FileInfo) string {
	switch m := fi.Mode(); {
	case m.IsDir():
		return "directory"
	case m.IsRegular():
		return "file"
	case m&os.ModeSymlink != 0:
		return "symlink"
	case m&os.ModeNamedPipe != 0:
		return "fifo"
	case m&os.ModeDevice != 0:
		return "dev"
	case m&os.ModeCharDevice != 0:
		return "chardev"
	case m&os.ModeSocket != 0:
		return "socket"
	case m&os.ModeIrregular != 0:
		return "irregular"
	default:
		return "unknown"
	}
}

func canRename(path string) (err error) {
	var info os.FileInfo

	if info, err = os.Lstat(path); err != nil {
		if os.IsNotExist(err) {
			// ok, well, now, it doesn't exist so I guess
			// we just continue?
			return nil
		} else {
			return errors.Wrapf(err, "failed to stat path %#v", path)
		}
	}

	mode := info.Mode()

	if !(isSymlink(mode) || mode.IsRegular() || mode.IsDir()) {
		return errors.Errorf("dest path %#v is a %s, cannot back up", path, nameForMode(info))
	}

	return nil
}

func doRename(path string) (err error) {
	if err = canRename(path); err != nil {
		return err
	}

	for i := 0; i < 100; i++ {
		bak := fp.Join(fp.Dir(path), fmt.Sprintf("%s.dfi_%s_%d", fp.Base(path), timestamp(), i))
		if err = os.Rename(path, bak); err != nil && !os.IsExist(err) {
			return errors.Wrapf(err, "falied to rename dest path %#v to %#v", path, bak)
		} else if err == nil {
			return nil
		}
	}

	return errors.Errorf("failed to back up path %#v", path)
}

// tis is actually 'unlink' as we remove the path that's in our way
// we will not remove a directory.
func doReplace(path string) error {
	return errors.Wrapf(os.Remove(path), "failed to remove %#v", path)
}

func (oc OnConflict) Handle(linkPath string) (skip bool, err error) {
	switch oc {
	case Rename:
		return false, doRename(linkPath)
	case Replace:
		return false, doReplace(linkPath)
	case Warn:
		log.Warnf("Destination %+v exists, skipping", linkPath)
		return true, nil
	case Fail:
		return false, errors.Errorf("Destination %#v exists, exiting", linkPath)
	default:
		panic(fmt.Sprintf("should never reach here: oc value: %#v", oc))
	}
}

func OnConflictForString(s string) (OnConflict, error) {
	switch str.ToLower(s) {
	case "rename":
		return Rename, nil
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
