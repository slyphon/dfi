package dotfile

import (
	"sort"
	"testing"

	fsf "github.com/slyphon/dfi/internal/fsfixture"
	"github.com/stretchr/testify/suite"
)

type (
	InstallerSuite struct {
		suite.Suite
		fsFix fsf.FsFixture
	}
)

var (
	_ suite.AfterTest  = &InstallerSuite{}
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
	ac := newApplyCollector()

	inst := &Installer{
		prefix:     ".",
		onConflict: Rename,
		apply:      ac.apply,
	}

	r := s.Require()

	err := inst.Run(s.fsFix.Dotfiles, s.fsFix.HomeDir)
	r.NoError(err)
	r.NotEmpty(ac.links)
	r.Len(ac.links, 3)

	sort.Sort(byVpath(ac.links))

	validate := func(ld LinkData, name string) {
		reld, err := ld.RelTo(s.fsFix.TempDir)
		r.NoError(err)
		r.Equal(reld.Vpath, "home/settings/dotfiles/"+name)
		r.Equal(reld.LinkPath, "home/."+name)
		r.Equal(reld.LinkData, "settings/dotfiles/"+name)
	}

	validate(ac.links[0], "bashrc")
	validate(ac.links[1], "vimrc")
	validate(ac.links[2], "zshrc")
}

func (s *InstallerSuite) TestBinfiles() {
	ac := newApplyCollector()

	inst := &Installer{
		prefix:     "",
		onConflict: Rename,
		apply:      ac.apply,
	}

	r := s.Require()
	err := inst.Run(s.fsFix.Binfiles, s.fsFix.LocalBinDir)

	r.NoError(err)
	r.NotEmpty(ac.links)
	r.Len(ac.links, 3)

	sort.Sort(byVpath(ac.links))

	validate := func(ld LinkData, name string) {
		reld, err := ld.RelTo(s.fsFix.TempDir)
		r.NoError(err)
		r.Equal(reld.Vpath, "home/settings/bin/"+name)
		r.Equal(reld.LinkPath, "home/.local/bin/"+name)
		r.Equal(reld.LinkData, "../../settings/bin/"+name)
	}

	validate(ac.links[0], "cat")
	validate(ac.links[1], "dog")
	validate(ac.links[2], "ls")
}
