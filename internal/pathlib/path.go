package pathlib

import (
	"github.com/gobwas/glob"
	"os"
	"path"
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
	}

	Path interface {
		PurePath
		IsBlockDevice() bool
		IsCharDevice() bool
		IsDir() bool
		IsFifo() bool
		IsFile() bool
		IsMount() bool
		IsSocket() bool
		IsSymlink() bool
		Exists() bool
		Lstat() (os.FileInfo, error)
		SymlinkTo(path string) error
		Rel(other string) (Path, error)
		Resolve() (Path, error)
		Mkdir(perm os.FileMode) error
		MkdirAll(perm os.FileMode) error
		Remove() error
		RemoveAll() error
	}

	purePath string
)

var _ PurePath = purePath("")

func NewPurePath(s string) PurePath {
	return purePath(s)
}

func (p purePath) Name() string {
	return path.Base(string(p))
}

func (p purePath) Parent() PurePath {
	return purePath(path.Dir(string(p)))
}

func (p purePath) Join(names ...string) PurePath {
	els := append([]string{string(p)}, names...)
	return NewPurePath(path.Join(els...))
}

func (p purePath) Clean() PurePath {
	return purePath(path.Clean(string(p)))
}

func (p purePath) Match(pattern string) (matched bool, err error) {
	return path.Match(pattern, string(p))
}

func (p purePath) ExMatch(pattern string) (matched bool, err error) {
	g, err := glob.Compile(pattern, '/')
	if err != nil {
		return false, err
	}
	return g.Match(p.ToString()), nil
}

func (p purePath) Split() (dir PurePath, name string) {
	d, n := path.Split(string(p))
	return purePath(d), n
}

func (p purePath) ToString() string {
	return string(p)
}

