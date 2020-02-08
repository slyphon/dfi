package dotfile

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type (
	ReaderSuite struct {
		RequireSuite
	}
)

func TestReader(t *testing.T) {
	suite.Run(t, new(ReaderSuite))
}

func testSplit(s *ReaderSuite, splitFn bufio.SplitFunc, expect []string, input string) {
	lines, err := ReadSources(strings.NewReader(input), splitFn)
	s.NoError(err)
	s.Equal(expect, lines)
}

var expect = []string{"/a/b/c", "/a/b/d", "/c/b/a"}

func (s *ReaderSuite) TestSplitOnNewlines() {
	testSplit(
		s,
		SplitOnNewlines,
		expect,
		"/a/b/c\n/a/b/d\n/c/b/a\n",
	)
}

func (s *ReaderSuite) TestSplitOnNewlinesMissingEOL() {
	testSplit(
		s,
		SplitOnNewlines,
		expect,
		"/a/b/c\n/a/b/d\n/c/b/a",
	)
}

func (s *ReaderSuite) TestSpltOnNull() {
	testSplit(
		s,
		SplitOnNullByte,
		expect,
		"/a/b/c\x00/a/b/d\x00/c/b/a\x00",
	)
}

func (s *ReaderSuite) TestSpltOnNullMissingLast() {
	testSplit(
		s,
		SplitOnNullByte,
		expect,
		"/a/b/c\x00/a/b/d\x00/c/b/a",
	)
}
