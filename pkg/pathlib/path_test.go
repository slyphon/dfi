package pathlib

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

// same as suite.Suite only the assertions are require and stop the test
// on the first failure
type RequireSuite struct {
	suite.Suite
	*require.Assertions
}

func (rs *RequireSuite) SetT(t *testing.T) {
	rs.Suite.SetT(t)
	rs.Assertions = rs.Suite.Require()
}

type PurePathSuite struct {
	RequireSuite
}

func TestPurePaths(t *testing.T) {
	suite.Run(t, new(PurePathSuite))
}

func (s *PurePathSuite) TestName() {
	pp := NewPurePath("/a/b/c")
	s.Implements((*PurePath)(nil), pp)
	s.Equal(pp.Name(), "c")
	s.Equal(NewPurePath("/a/b/c/").Name(), "c")
}

func (s *PurePathSuite) TestParent() {
	pp := NewPurePath("/a/b/c")
	s.Equal("/a/b", pp.Parent().String())
	s.Equal("/", NewPurePath("/").Parent().String())
}

func (s *PurePathSuite) TestJoin() {
	s.Implements((*PurePath)(nil), NewPurePath("/a/b/c").Join("x", "y", "z"))

	validate := func(expect, base string, extra ...string) {
		x := NewPurePath(expect)
		b := NewPurePath(base)

		result := b.Join(extra...)

		s.Equal(x.String(), result.String())
	}

	validate("/a/b/c/x/y/z", "/a/b/c", "x", "y", "z")
	validate("/a/q", "/a/b/c", "..", "..", "q")
	validate("/y", "/a/b/c", "..", "..", "..", "..", "y")
	validate("/a/b/x/y/z", "/a/b//x", "y", "z")
	validate("/a/b/c/1/2/3", "/a/b/c", "/1/2/3")
	validate("a/b/c/1/2/3", "a/b/c", "/1/2/3")
}

func (s *PurePathSuite) TestClean() {
	p := NewPurePath("/a///b/./../b//c")
	s.Equal("/a/b/c", p.Clean().String())
}

func (s *PurePathSuite) TestMatch() {
	b, err := NewPurePath("/foo/xyz.tar").Match("/*/*.tar")
	s.NoError(err)
	s.True(b)
}

func (s *PurePathSuite) TestExMatch() {
	b, err := NewPurePath("/foo/bar/baz.tar").ExMatch("**/*.tar")
	s.NoError(err)
	s.True(b)
}

func (s *PurePathSuite) TestSplit() {
	head, tail := NewPurePath("/a/b/c").Split()
	s.Equal("/a/b/", head.String())
	s.Equal("c", tail)
}

