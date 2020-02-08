package pathlib

import (
	"io/ioutil"
	"os"
	fp "path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero/mem"
	"github.com/stretchr/testify/suite"

	"github.com/slyphon/dfi/pkg/testhelper"
)

type PosixPathSuite struct {
	testhelper.DFISuite
	tmpdir *string
	origFs Fs
}

var _ suite.AfterTest = &PosixPathSuite{}

func (s *PosixPathSuite) TempDir() string {
	if s.tmpdir == nil {
		path, err := ioutil.TempDir("", "posixpathsuite")
		s.NoError(err)
		s.tmpdir = &path
	}

	return *s.tmpdir
}

// some tests need the real filesystem, this will put it back in place
func (s *PosixPathSuite) usingOsFs() {
	if s.origFs != nil {
		_ = SetFs(s.origFs)
		s.origFs = nil
	}
}

func TestPosixPaths(t *testing.T) {
	s := new(PosixPathSuite)

	s.AddBeforeHook(func (a, b string) {
		s.origFs = SetFs(NewMemFsPlus())
	})

	s.AddAfterHook(func (a, b string) {
		s.usingOsFs()
		if s.tmpdir != nil {
			err := os.RemoveAll(*s.tmpdir)
			s.tmpdir = nil
			if err != nil {
				log.Errorf("got error in AfterTest: %#v", err)
			}
		}
	})

	suite.Run(t, s)
}

func (s *PosixPathSuite) TestIsMount() {
	s.usingOsFs()
	pp := NewPosixPath("/")
	s.True(pp.IsMount())

	wd, err := os.Getwd()
	s.NoError(err)
	s.False(NewPosixPath(wd).IsMount())
}

func (s *PosixPathSuite) touch(path string) {
	file, err := os.Create(path)
	s.NoError(err)
	s.NoError(file.Close())
}

func (s *PosixPathSuite) TestIsSameFile() {
	s.usingOsFs()
	d := s.TempDir()

	filePath := fp.Join(d, "file")
	linkPath := fp.Join(d, "link")

	s.touch(filePath)
	s.NoError(os.Symlink("./file", linkPath))

	b, err := NewPosixPath(filePath).SameFile(NewPosixPath(linkPath))
	s.True(b)
	s.NoError(err)
}

func (s *PosixPathSuite) TestExMatch() {
	pp := NewPosixPath("/this/is/a/dir/the.tar")

	shouldMatch := func(pattern string) {
		b, e := pp.ExMatch(pattern)
		s.NoError(e)
		s.True(b)
	}

	shouldMatch("**/*.tar")
	shouldMatch("**/*.{tar,tar.gz}")
}

func (s *PosixPathSuite) createFileWithMode(path string, mode os.FileMode) PosixPath {
	fh, err := fs.Create(path)
	s.NoError(err)
	under, ok := fh.(*mem.File)
	s.True(ok)
	fdata := under.Data()
	mem.SetMode(fdata, mode)
	return NewPosixPath(path)
}

const PermBits = 0o644

func (s *PosixPathSuite) TestIsDir() {
	path := "/path/to/dir"
	s.NoError(fs.Mkdir(path, 0o755))
	pp := NewPosixPath(path)
	s.True(pp.IsDir())
}

func (s *PosixPathSuite) TestIsCharDevice() {
	pp := s.createFileWithMode("/path/to/chardev", os.ModeCharDevice|PermBits)
	s.True(pp.IsCharDevice())
}

func (s *PosixPathSuite) TestIsFifo() {
	pp := s.createFileWithMode("/path", os.ModeNamedPipe|PermBits)
	s.True(pp.IsFifo())
}

func (s *PosixPathSuite) TestIsSocket() {
	pp := s.createFileWithMode("/path", os.ModeSocket|PermBits)
	s.True(pp.IsSocket())
}

func (s *PosixPathSuite) TestIsBlockDevice() {
	pp := s.createFileWithMode("/path", os.ModeDevice|PermBits)
	s.True(pp.IsBlockDevice())
}


