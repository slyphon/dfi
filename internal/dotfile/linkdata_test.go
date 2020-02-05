package dotfile

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type LinkDataSuite struct {
	suite.Suite
}

func TestHookUpSuite(t *testing.T) {
	suite.Run(t, new(LinkDataSuite))
}

func (s *LinkDataSuite) TestFindCommonRoot() {
	s.Equal("/a/b/c/d", FindCommonRoot("/a/b/c/d/e/f", "/a/b/c/d/q/e/r"))
	s.Equal("/a/b", FindCommonRoot("/a/b/c/d", "/a/b/c"))
	s.Equal("", FindCommonRoot("/", "/a/b/c"))
	s.Equal("", FindCommonRoot("/qwer", "/a/b/c"))
}

func (s *LinkDataSuite) TestLinkDataFor() {
	ld, err := LinkDataFor("/home/x/.settings/bashrc", "/home/x", ".")

	s.NoError(err)

	s.Equal(
		LinkData{
				Vpath:    "/home/x/.settings/bashrc",
				LinkPath: "/home/x/.bashrc",
				LinkData: ".settings/bashrc",
		},
		ld,
	)

	ld, err = LinkDataFor("/home/x/.settings/bin/foo", "/home/x/.local/bin", "")
	s.NoError(err)

	s.Equal(
		LinkData{
				Vpath:    "/home/x/.settings/bin/foo",
				LinkPath: "/home/x/.local/bin/foo",
				LinkData: "../../.settings/bin/foo",
		},
		ld,
	)

	ld, err = LinkDataFor("/home/x/.settings/bin/foo", "/path/to/blah", "")
	s.NoError(err)

	s.Equal(
		LinkData{
				Vpath:    "/home/x/.settings/bin/foo",
				LinkPath: "/path/to/blah/foo",
				LinkData: "/home/x/.settings/bin/foo",
		},
		ld,
	)
}
