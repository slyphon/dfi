package dotfile

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"

	fsf "github.com/slyphon/dfi/internal/fsfixture"
	pl "github.com/slyphon/dfi/pkg/pathlib"
)

type (
	InstallerSuite struct {
		RequireSuite
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

func (s *InstallerSuite) TestInstallerDotFiles() {
	i := NewInstaller(".", ConflictHandlers.Replace)
	err := i.Run(s.fsFix.Dotfiles, s.fsFix.HomeDir)
	s.NoError(err)

	dotPaths := make([]pl.PosixPath, 0, len(s.fsFix.Dotfiles))
	for _, p := range s.fsFix.Dotfiles {
		dotPaths = append(dotPaths, pl.NewPosixPath(p))
	}

	home := pl.NewPosixPath(s.fsFix.HomeDir)

	entries, err := home.Glob(".*")
	s.NoError(err)
	sort.Sort(pl.LexOrder(entries))
	s.Len(entries, 4)                             // also has the ".local" dir
	entries = append(entries[:1], entries[2:]...) // cut out .local

	expectNames := []string{".bashrc", ".vimrc", ".zshrc"}
	for i, pp := range entries {
		x := expectNames[i]
		s.Equal(x, pp.Name())
		b, e := dotPaths[i].SameFile(pp)
		s.NoError(e)
		s.True(b)
	}
}

func (s *InstallerSuite) TestInstallerBinFiles() {
	i := NewInstaller("", ConflictHandlers.Replace)
	err := i.Run(s.fsFix.Binfiles, s.fsFix.LocalBinDir)
	s.NoError(err)

	entries, err := pl.NewPosixPath(s.fsFix.LocalBinDir).Glob("*")
	s.NoError(err)
	sort.Sort(pl.LexOrder(entries))
	s.Len(entries, 3)

	expectNames := []string{"cat", "dog", "ls"}
	for i, pp := range entries {
		s.Equal(expectNames[i], pp.Name())
	}
}
