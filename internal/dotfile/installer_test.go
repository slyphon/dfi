package dotfile

import (
	"testing"

	"github.com/stretchr/testify/suite"
	fsf "github.com/slyphon/dfi/internal/fsfixture"
)

type (
	InstallerSuite struct {
		suite.Suite
		fsFix fsf.FsFixture
	}
)

var (
	_ suite.AfterTest = &InstallerSuite{}
	_ suite.BeforeTest = &InstallerSuite{}
)

func TestInstaller(t *testing.T) {
	suite.Run(t, new(InstallerSuite))
}

func (s *InstallerSuite) BeforeTest(a, b string) {
	s.fsFix = fsf.NewFsFixture()
}

func (s *InstallerSuite) AfterTest(a, b string) {
	s.fsFix.Cleanup()
}

func (s *InstallerSuite) TestDotfiles() {
	stg := Settings{
		
	}
}
