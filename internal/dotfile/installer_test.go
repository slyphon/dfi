package dotfile

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"

	fsf "github.com/slyphon/dfi/internal/fsfixture"
	pl "github.com/slyphon/dfi/pkg/pathlib"
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

	err := inst.Run(pl.PosixSliceStringer(s.fsFix.Dotfiles), s.fsFix.HomeDir.String())
	r.NoError(err)
	r.NotEmpty(ac.links)
	r.Len(ac.links, 4)

	sort.Sort(byVpath(ac.links))

	validate := func(name string) {
		ld := LinkDataOpsT(ac.links).
			Find(func(ld LinkData) bool {
				return pl.NewPurePath(ld.Vpath).Must().ExMatch("**/" + name)
			})
		s.NotNil(ld)

		reld, err := ld.RelTo(s.fsFix.TempDir.String())
		r.NoError(err)
		r.Equal(reld.Vpath, "home/settings/dotfiles/"+name)
		r.Equal(reld.LinkPath, "home/."+name)
		r.Equal(reld.LinkData, "settings/dotfiles/"+name)
	}

	validate("bashrc")
	validate("vimrc")
	validate("zshrc")
}

func (s *InstallerSuite) TestBinfiles() {
	ac := newApplyCollector()

	inst := &Installer{
		prefix:     "",
		onConflict: Rename,
		apply:      ac.apply,
	}

	r := s.Require()
	err := inst.Run(pl.PosixSliceStringer(s.fsFix.Binfiles), s.fsFix.LocalBinDir.String())

	r.NoError(err)
	r.NotEmpty(ac.links)
	r.Len(ac.links, 3)

	sort.Sort(byVpath(ac.links))

	validate := func(ld LinkData, name string) {
		reld, err := ld.RelTo(s.fsFix.TempDir.String())
		r.NoError(err)
		r.Equal(reld.Vpath, "home/settings/bin/"+name)
		r.Equal(reld.LinkPath, "home/.local/bin/"+name)
		r.Equal(reld.LinkData, "../../settings/bin/"+name)
	}

	validate(ac.links[0], "cat")
	validate(ac.links[1], "dog")
	validate(ac.links[2], "ls")
}

func (s *InstallerSuite) TestInstallerDotFiles() {
	i := NewInstaller(".", ConflictHandlers.Replace)
	err := i.Run(pl.PosixSliceStringer(s.fsFix.Dotfiles), s.fsFix.HomeDir.String())
	s.NoError(err)

	dotPaths := pl.PosixPathOpsT(s.fsFix.Dotfiles).
		Filter(func(pp pl.PosixPath) bool {
			return !pp.Must().ExMatch("**/config")
		})

	sort.Sort(pl.LexOrderPosix(dotPaths))

	home := s.fsFix.HomeDir

	entries, err := home.Glob(".*")
	s.NoError(err)
	sort.Sort(pl.LexOrderPosix(entries))
	s.Len(entries, 5) // also has the ".local" and ".config" dirs

	entries = pl.PosixPathOpsT(entries).
		Filter(func(pp pl.PosixPath) bool {
			return !pp.Must().ExMatch("**/.{config,local}")
		})

	s.Len(entries, 3)

	expectNames := []string{".bashrc", ".vimrc", ".zshrc"}
	for i, pp := range entries {
		x := expectNames[i]
		s.Equal(x, pp.Name())
		b, e := dotPaths[i].SameFile(pp)
		s.NoError(e)
		s.True(b, fmt.Sprintf("expected %+v to be the same file as %+v", dotPaths[i], pp))
	}
}

func (s *InstallerSuite) TestInstallerBinFiles() {
	i := NewInstaller("", ConflictHandlers.Replace)
	err := i.Run(pl.PosixSliceStringer(s.fsFix.Binfiles), s.fsFix.LocalBinDir.String())
	s.NoError(err)

	entries, err := s.fsFix.LocalBinDir.Glob("*")
	s.NoError(err)
	sort.Sort(pl.LexOrderPosix(entries))
	s.Len(entries, 3)

	expectNames := []string{"cat", "dog", "ls"}
	for i, pp := range entries {
		s.Equal(expectNames[i], pp.Name())
	}
}
