package pathlib

import (
	"fmt"
	"path"

	"github.com/gobwas/glob"
)

type (
	PurePath interface {
		fmt.Stringer
		FlavorCasts

		Name() string
		Parent() PurePath
		Join(names ...string) PurePath
		Clean() PurePath
		Match(pattern string) (matched bool, err error)

		// ExMatch performs an extended match on this PurePath
		// see github.com/gobwas/glob for syntax.
		// (basically '**' is supported)
		ExMatch(pattern string) (matched bool, err error)
		Split() (dir PurePath, file string)
	}

	pureStr string
)

var _ PurePath = pureStr("")
var _ fmt.Stringer = pureStr("")

func NewPurePath(s string) PurePath { return pureStr(s) }
func (p pureStr) Name() string      { return path.Base(string(p)) }
func (p pureStr) Parent() PurePath  { return pureStr(path.Dir(string(p))) }
func (p pureStr) Clean() PurePath   { return pureStr(path.Clean(string(p))) }
func (p pureStr) Posix() PosixPath  { return pure2posix(p) }
func (p pureStr) Pure() PurePath    { return p }
func (p pureStr) Must() MustActions { return posix2must(pure2posix(p)) }
func (p pureStr) String() string    { return string(p) }

func (p pureStr) Join(names ...string) PurePath {
	els := append([]string{string(p)}, names...)
	return NewPurePath(path.Join(els...))
}

func (p pureStr) Match(pattern string) (matched bool, err error) {
	return path.Match(pattern, string(p))
}

func (p pureStr) ExMatch(pattern string) (matched bool, err error) {
	var g glob.Glob
	g, err = glob.Compile(pattern, '/')
	if err != nil {
		return false, err
	}
	return g.Match(p.String()), nil
}

func (p pureStr) Split() (dir PurePath, name string) {
	d, n := path.Split(string(p))
	return pureStr(d), n
}
