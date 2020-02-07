package pathlib

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type PurePathSuite struct {
	suite.Suite
	r *require.Assertions
}

func (s *PurePathSuite) BeforeTest(a, b string) {
	s.r = s.Require()
}

func TestPurePaths(t *testing.T) {
	suite.Run(t, new(PurePathSuite))
}


func (s *PurePathSuite) TestName() {
	pp := NewPurePath("/a/b/c")
	s.r.Implements((*PurePath)(nil), pp)

	s.r.Equal(pp.Name(), "c")
	s.r.Equal(NewPurePath("/a/b/c/").Name(), "c")
}

func (s *PurePathSuite) TestParent() {
	pp := NewPurePath("/a/b/c")
	s.r.Equal("/a/b", pp.Parent().ToString())
	s.r.Equal("/", NewPurePath("/").Parent().ToString())
}

func (s *PurePathSuite) TestJoin() {
	s.r.Implements((*PurePath)(nil), NewPurePath("/a/b/c").Join("x", "y", "z"))

	validate := func (expect, base string, extra... string) {
		x := NewPurePath(expect)
		b := NewPurePath(base)

		result := b.Join(extra...)

		s.r.Equal(x.ToString(), result.ToString())
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
	s.r.Equal("/a/b/c", p.Clean().ToString())
}

func (s *PurePathSuite) TestMatch() {
	b, err := NewPurePath("/foo/xyz.tar").Match("/*/*.tar")
	s.r.NoError(err)
	s.r.True(b)
}

func (s *PurePathSuite) TestExMatch() {
	b, err := NewPurePath("/foo/bar/baz.tar").ExMatch("**/*.tar")
	s.r.NoError(err)
	s.r.True(b)
}

func (s *PurePathSuite) TestSplit() {
	head, tail := NewPurePath("/a/b/c").Split()
	s.r.Equal("/a/b/", head.ToString())
	s.r.Equal("c", tail)
}

type PosixPathSuite struct {
	suite.Suite
	r *require.Assertions
}

func (s *PosixPathSuite) BeforeTest(a, b string) {
	s.r = s.Require()
}

func TestPosixPaths(t *testing.T) {
	suite.Run(t, new(PosixPathSuite))
}

func (s *PosixPathSuite) TestIsMount() {
	pp := NewPosixPath("/")
	b, e := pp.IsMount()
	s.r.NoError(e)
	s.r.True(b)

	wd, err := os.Getwd()
	s.r.NoError(err)
	pp = NewPosixPath(wd)
	b, e = pp.IsMount()
	s.r.NoError(e)
	s.r.False(b)
}


