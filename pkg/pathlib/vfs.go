package pathlib

import (
	"os"
	fp "path/filepath"

	"github.com/spf13/afero"
)

// integration with afero's filesystem abstractions
// to allow for testing

type (
	Symlinker interface {
		Symlink(old, new string) error
	}

	Globber interface {
		Glob(pattern string) ([]string, error)
	}


	Fs interface {
		afero.Fs
		afero.Lstater
		Symlinker
		Globber
	}

	pathLibOsFs struct {
		*afero.OsFs
	}
)

func (p pathLibOsFs) Symlink(old, new string) error {
	return os.Symlink(old, new)
}

func (p pathLibOsFs) Glob(pattern string) ([]string, error) {
	return fp.Glob(pattern)
}

func newPathLibOsFs() pathLibOsFs {
	return pathLibOsFs{
		OsFs: &afero.OsFs{},
	}
}

var fs Fs

func init() {
	fs = newPathLibOsFs()
}

func SetFs(newFs Fs) (oldFs Fs) {
	oldFs, fs = fs, newFs
	return oldFs
}

func GetFs() Fs {
	return fs
}
