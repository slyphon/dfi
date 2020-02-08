package pathlib

import (
	"fmt"
	"os"
	fp "path/filepath"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type (
	// Provide path manipulation like PurePath, but also
	// posix actions on the real filesystem. Note that the
	// FileInfoAdditions actions on a PosixPath will ignore all
	// errors and return false. This is the strategy that
	// the Python stdlib 'pathlib' follows in the interest of
	// ease of use. If you really need to know about every
	// error type, then use Stat() and inspect the error case before
	// interrogating the FileInfoAdditions method.
	PosixPath interface {
		FileInfoAdditions
		FlavorCasts
		fmt.Stringer

		Name() string
		Parent() PosixPath
		Join(names ...string) PosixPath
		Clean() PosixPath
		Match(pattern string) (matched bool, err error)
		ExMatch(pattern string) (matched bool, err error)
		Split() (dir PosixPath, file string)
		Stat() (RichFileInfo, error)
		Lstat() (RichFileInfo, error)
		SymlinkTo(path string) error
		Rel(other string) (PosixPath, error)
		Resolve() (PosixPath, error)
		Mkdir(perm os.FileMode) error
		MkdirAll(perm os.FileMode) error
		Remove() error
		RemoveAll() error
		Lexists() bool
		Exists() bool
		IsMount() bool
		SameFile(other PosixPath) (bool, error)
		IsDir() bool
		Glob(pattern string) ([]PosixPath, error)
		Touch(perm os.FileMode, existOk bool) error
	}

	posixStr string
)

var _ PosixPath = posixStr("")

func NewPosixPath(s string) PosixPath             { return posixStr(s) }
func (p posixStr) Name() string                   { return posix2pure(p).Name() }
func (p posixStr) Parent() PosixPath              { return posix2pure(p).Parent().Posix() }
func (p posixStr) Clean() PosixPath               { return posix2pure(p).Clean().Posix() }
func (p posixStr) Join(names ...string) PosixPath { return posix2pure(p).Join(names...).Posix() }
func (p posixStr) String() string                 { return string(p) }
func (p posixStr) Pure() PurePath                 { return posix2pure(p) }
func (p posixStr) Posix() PosixPath               { return p }

func (p posixStr) Match(pattern string) (matched bool, err error) {
	return posix2pure(p).Match(pattern)
}

func (p posixStr) ExMatch(pattern string) (matched bool, err error) {
	return posix2pure(p).ExMatch(pattern)
}

func (p posixStr) Split() (dir PosixPath, file string) {
	d, f := posix2pure(p).Split()
	return d.Posix(), f
}

func (p posixStr) Stat() (RichFileInfo, error) {
	info, err := fs.Stat(string(p))
	return richInfo{info}, err
}

func (p posixStr) Lstat() (RichFileInfo, error) {
	info, didLstat, err := fs.LstatIfPossible(string(p))
	if !didLstat {
		return nil, errors.Errorf("Lstat not available on the filesystem %#v", fs.Name())
	}
	return richInfo{info}, err
}

func (p posixStr) SymlinkTo(path string) error {
	return fs.Symlink(path, string(p))
}

func (p posixStr) Rel(other string) (PosixPath, error) {
	pp, err := fp.Rel(string(p), other)
	if err != nil {
		return nil, err
	}
	return PosixPath(posixStr(pp)), nil
}

func (p posixStr) Resolve() (PosixPath, error) {
	pp, err := fp.EvalSymlinks(string(p))
	if err != nil {
		return nil, err
	}
	return PosixPath(posixStr(pp)), nil
}

func (p posixStr) Mkdir(perm os.FileMode) error {
	return fs.Mkdir(string(p), perm)
}

func (p posixStr) MkdirAll(perm os.FileMode) error {
	return fs.MkdirAll(string(p), perm)
}

func (p posixStr) Remove() error {
	return fs.Remove(string(p))
}

func (p posixStr) RemoveAll() error {
	return fs.RemoveAll(string(p))
}

func (p posixStr) Lexists() bool {
	_, err := p.Lstat()
	return err == nil
}

func (p posixStr) Exists() bool {
	_, err := p.Stat()
	return !(os.IsNotExist(err) || IsNotDir(err)) || err == nil
}

func getStatT(f os.FileInfo) (*syscall.Stat_t, error) {
	statT, ok := f.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("failed to get system Stat_t")
	}
	return statT, nil
}

func (p posixStr) IsMount() bool {
	var self os.FileInfo
	var selfStatT *syscall.Stat_t
	var err error

	if self, err = p.Stat(); err != nil || !self.IsDir() {
		return false
	}

	if selfStatT, err = getStatT(self); err != nil {
		return false
	}

	var myParentInfo os.FileInfo

	if myParentInfo, err = p.Parent().Posix().Stat(); err != nil {
		return false
	}

	parentStatT, err := getStatT(myParentInfo)
	if err != nil {
		return false
	}

	// from linux coreutils:
	// int is_not_mnt = (st_dev == st.st_dev) && (st_ino != st.st_ino);
	//
	isNotMnt := (selfStatT.Dev == parentStatT.Dev) && (selfStatT.Ino != parentStatT.Ino)

	return !isNotMnt
}

func (p posixStr) SameFile(other PosixPath) (b bool, err error) {
	this, err := p.Stat()
	if err != nil {
		return false, err
	}
	that, err := other.Stat()
	if err != nil {
		return false, err
	}
	return os.SameFile(this.getInfo(), that.getInfo()), nil
}

func (p posixStr) Glob(pattern string) ([]PosixPath, error) {
	matches, err := fs.Glob(fp.Join(p.String(), pattern))
	if err != nil {
		return nil, err
	}
	paths := make([]PosixPath, 0, len(matches))
	for _, m := range matches {
		paths = append(paths, NewPosixPath(m))
	}
	return paths, err
}

func (p posixStr) IsBlockDevice() bool {
	info, err := p.Stat()
	return err == nil && info.IsBlockDevice()
}

func (p posixStr) IsCharDevice() bool {
	info, err := p.Stat()
	return err == nil && info.IsCharDevice()
}

func (p posixStr) IsFifo() bool {
	info, err := p.Stat()
	return err == nil && info.IsFifo()
}

func (p posixStr) IsFile() bool {
	info, err := p.Stat()
	return err == nil && info.IsFile()
}

func (p posixStr) IsSocket() bool {
	info, err := p.Stat()
	return err == nil && info.IsSocket()
}

func (p posixStr) IsSymlink() bool {
	info, err := p.Lstat()
	return err == nil && info.IsSymlink()
}

func (p posixStr) IsDir() bool {
	info, err := p.Stat()
	return err == nil && info.IsDir()
}

func (p posixStr) Touch(perm os.FileMode, existOk bool) (err error) {
	now := time.Now()

	if existOk {
		switch err = fs.Chtimes(p.String(), now, now); {
		case os.IsNotExist(err): // guess we have to create it
		case err == nil:
			return nil // file existed, we changed the time, we're done here
		default:
			return err // something went wrong, abort! abort!
		}
	}

	flags := os.O_CREATE | os.O_WRONLY
	if !existOk {
		flags = flags | os.O_EXCL
	}

	var file afero.File
	file, err = fs.OpenFile(p.String(), flags, perm)
	defer func() {
		if file != nil {
			if cerr := file.Close(); cerr != nil {
				err = cerr
			}
		}
	}()

	return
}

func (p posixStr) Must() MustActions {
	return NewMustAction(p.String())
}

type (
	PosixPathOps interface {
		Filter(fn func(pp PosixPath) bool) []PosixPath
		Find(fn func(pp PosixPath) bool) *PosixPath
	}
	PosixPathOpsT []PosixPath
)

func (p PosixPathOpsT) Filter(fn func(pp PosixPath) bool) []PosixPath {
	rv := make([]PosixPath, 0, len(p))
	for _, pp := range p {
		if fn(pp) {
			rv = append(rv, pp)
		}
	}
	return rv
}

func (p PosixPathOpsT) Find(fn func(pp PosixPath) bool) (found *PosixPath) {
	for _, pp := range p {
		if fn(pp) {
			return &pp
		}
	}
	return nil
}
