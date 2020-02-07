package pathlib

import (
	"fmt"
	"github.com/gobwas/glob"
	"os"
	"path"
	fp "path/filepath"
	"syscall"
	"time"
)

type (
	FileStat interface {
		Exists() bool
		IsSymlink() bool
		IsFile() bool
		IsDir() bool
		IsOther() bool
		FileInfo() os.FileInfo
	}

	PurePath interface {
		Name() string
		Parent() PurePath
		Join(names... string) PurePath
		Clean() PurePath
		Match(pattern string) (matched bool, err error)

		// ExMatch performs an extended match on this PurePath
		// see github.com/gobwas/glob for syntax.
		// (basically '**' is supported)
		ExMatch(pattern string) (matched bool, err error)
		Split() (dir PurePath, file string)
		ToString() string
		ToPosix() PosixPath
	}

	RichFileInfo interface {
		os.FileInfo
		IsBlockDevice() bool
		IsCharDevice() bool
		IsFifo() bool
		IsFile() bool
		IsSocket() bool
		IsSymlink() bool
	}

	PosixPath interface {
		PurePath
		Stat() (RichFileInfo, error)
		Lstat() (RichFileInfo, error)
		SymlinkTo(path string) error
		Rel(other string) (PosixPath, error)
		Resolve() (PosixPath, error)
		Mkdir(perm os.FileMode) error
		MkdirAll(perm os.FileMode) error
		Remove() error
		RemoveAll() error
		Lexists() (bool, error)
		Exists() (bool, error)
		IsMount() (bool, error)
	}

	pathStr string
)

var _ PurePath = pathStr("")

func NewPurePath(s string) PurePath {
	return pathStr(s)
}

func NewPosixPath(s string) PosixPath {
	return NewPurePath(s).ToPosix()
}

func (p pathStr) Name() string {
	return path.Base(string(p))
}

func (p pathStr) Parent() PurePath {
	return pathStr(path.Dir(string(p)))
}

func (p pathStr) Join(names ...string) PurePath {
	els := append([]string{string(p)}, names...)
	return NewPurePath(path.Join(els...))
}

func (p pathStr) Clean() PurePath {
	return pathStr(path.Clean(string(p)))
}

func (p pathStr) Match(pattern string) (matched bool, err error) {
	return path.Match(pattern, string(p))
}

func (p pathStr) ExMatch(pattern string) (matched bool, err error) {
	g, err := glob.Compile(pattern, '/')
	if err != nil {
		return false, err
	}
	return g.Match(p.ToString()), nil
}

func (p pathStr) Split() (dir PurePath, name string) {
	d, n := path.Split(string(p))
	return pathStr(d), n
}

func (p pathStr) ToPosix() PosixPath {
	return PosixPath(p)
}

func (p pathStr) ToString() string {
	return string(p)
}

type richFileInfo struct {
	info os.FileInfo
}

var _ RichFileInfo = richFileInfo{nil}
var _ os.FileInfo = richFileInfo{nil}

func (r richFileInfo) Name() string {
	return r.info.Name()
}

func (r richFileInfo) Size() int64 {
	return r.info.Size()
}

func (r richFileInfo) Mode() os.FileMode {
	return r.info.Mode()
}

func (r richFileInfo) ModTime() time.Time {
	return r.info.ModTime()
}

func (r richFileInfo) IsDir() bool {
	return r.info.IsDir()
}

func (r richFileInfo) Sys() interface{} {
	return r.info.Sys()
}

func (r richFileInfo) IsBlockDevice() bool {
	return r.info.Mode()&os.ModeDevice != 0
}

func (r richFileInfo) IsCharDevice() bool {
	return r.info.Mode()&os.ModeCharDevice != 0
}

func (r richFileInfo) IsFifo() bool {
	return r.info.Mode()&os.ModeNamedPipe != 0
}

func (r richFileInfo) IsFile() bool {
	return r.info.Mode()&os.ModeType == 0
}

func (r richFileInfo) IsSocket() bool {
	return r.info.Mode()&os.ModeSocket != 0
}

func (r richFileInfo) IsSymlink() bool {
	return r.info.Mode()&os.ModeSymlink != 0
}

func (p pathStr) Stat() (RichFileInfo, error) {
	info, err := os.Stat(string(p))
	return richFileInfo{info}, err
}

func (p pathStr) Lstat() (RichFileInfo, error) {
	info, err := os.Lstat(string(p))
	return richFileInfo{info}, err
}

func (p pathStr) SymlinkTo(path string) error {
	return os.Symlink(path, string(p))
}

func (p pathStr) Rel(other string) (PosixPath, error) {
	pp, err := fp.Rel(string(p), other)
	if err != nil {
		return nil, err
	}
	return PosixPath(pathStr(pp)), nil
}

func (p pathStr) Resolve() (PosixPath, error) {
	pp, err := fp.EvalSymlinks(string(p))
	if err != nil {
		return nil, err
	}
	return PosixPath(pathStr(pp)), nil
}

func (p pathStr) Mkdir(perm os.FileMode) error {
	return os.Mkdir(string(p), perm)
}

func (p pathStr) MkdirAll(perm os.FileMode) error {
	return os.MkdirAll(string(p), perm)
}

func (p pathStr) Remove() error {
	return os.Remove(string(p))
}

func (p pathStr) RemoveAll() error {
	return os.RemoveAll(string(p))
}

func (p pathStr) Lexists() (bool, error) {
	panic("not implemented")
}

func (p pathStr) Exists() (bool, error) {
	_, err := p.Stat()
	switch {
	case err == nil:
		return true, nil
	case os.IsNotExist(err):
		return false, nil
	default:
		return false, err
	}
}

func getStatT(f os.FileInfo) (*syscall.Stat_t, error) {
	statT, ok := f.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("failed to get system Stat_t")
	}
	return statT, nil
}

func (p pathStr) IsMount() (b bool, err error) {
	var self os.FileInfo
	var selfStatT *syscall.Stat_t

	if self, err = p.Stat(); err != nil || !self.IsDir() {
		if os.IsNotExist(err) {
			// swallow the ENOENT here
			return false, nil
		} else {
			return false, err
		}
	}

	if selfStatT, err = getStatT(self); err != nil {
		return false, err
	}

	var myParentInfo os.FileInfo

	if myParentInfo, err = p.Parent().ToPosix().Stat(); err != nil {
		return false, err
	}

	parentStatT, err := getStatT(myParentInfo)
	if err != nil {
		return false, err
	}

	// from linux coreutils:
	// int is_not_mnt = (st_dev == st.st_dev) && (st_ino != st.st_ino);
	//
	isNotMnt := (selfStatT.Dev == parentStatT.Dev) && (selfStatT.Ino != parentStatT.Ino)

	return !isNotMnt, nil
}


