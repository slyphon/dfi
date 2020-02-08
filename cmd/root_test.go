package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"

	df "github.com/slyphon/dfi/internal/dotfile"
	testhelp "github.com/slyphon/dfi/pkg/testhelper"
)

type (
	RootCmdSuite struct {
		testhelp.DFISuite
		tmpdir string
	}

	RunMock struct {
		settings *df.Settings
	}
)

func (rm *RunMock) Run(s *df.Settings) error {
	rm.settings = s
	return nil
}

func TestRootCmdSuite(t *testing.T) {
	s := new(RootCmdSuite)
	s.AddBeforeHook(func (a, b string) {
		var err error
		s.tmpdir, err = ioutil.TempDir("", "rootcmdsuite")
		s.NoError(err)
	})
	s.AddAfterHook(func (a, b string){
		if (s.tmpdir != "") {
			err := os.RemoveAll(s.tmpdir)
			log.Errorf("ignoring error tearing down tmpdir: %+v", err)
		}
	})
	suite.Run(t, s)
}

func (s *RootCmdSuite) TestCommandSanity() {
	rm := &RunMock{}
	var rootCmd *cobra.Command

	rootCmd = NewRootCommand(rm.Run)
	rootCmd.SetArgs([]string{"-p", ".", "/a/b/c/settings", "/a/b/c/home"})
	err := rootCmd.Execute()
	s.NoError(err)
	s.Equal(".", rm.settings.Prefix)
	s.Equal([]string{"/a/b/c/settings"}, rm.settings.SourcePaths)
	s.Equal("/a/b/c/home", rm.settings.DestPath)
}
