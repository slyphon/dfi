package pathlib

import (
	"os"

	"github.com/spf13/afero"
)

type (
	MemFsPlus struct {
		*afero.MemMapFs
	}
)

var _ Fs = &MemFsPlus{}

func (m *MemFsPlus) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	panic("Not Implemented")
}

func (m *MemFsPlus) Symlink(old, new string) error         { panic("Not Implemented") }
func (m *MemFsPlus) Glob(pattern string) ([]string, error) { panic("Not Implemented") }

func NewMemFsPlus() *MemFsPlus {
	return &MemFsPlus{
		MemMapFs: &afero.MemMapFs{},
	}
}

