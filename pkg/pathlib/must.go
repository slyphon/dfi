package pathlib

import (
	"fmt"
	"os"
)

type (
	// these provide panic-inducing versions of the PosixPath actions
	// for when you're pretty sure it's just gonna work and want the convenience
	MustActions interface {
		FlavorCasts
		fmt.Stringer

		Mkdir(perm os.FileMode)
		MkdirAll(perm os.FileMode)
		Remove()
		RemoveAll()
		Touch(perm os.FileMode, existOk bool)

		// Match returns true if the pattern matches the last path component
		// of this path. Note that ths wlll only blow up if the pattern itself
		// is invalid
		Match(pattern string) bool

		// ExMatch performs an extended match on the entire path string using
		// the globbing library github.com/gobwas/glob. This function will only
		// blow up if the pattern is not valid
		ExMatch(pattern string) bool
	}

	mustPath string
)

var _ MustActions = mustPath("")

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func NewMustAction(s string) MustActions { return MustActions(mustPath(s)) }

func (m mustPath) Mkdir(perm os.FileMode)               { must(m.Posix().Mkdir(perm)) }
func (m mustPath) MkdirAll(perm os.FileMode)            { must(m.Posix().MkdirAll(perm)) }
func (m mustPath) Remove()                              { must(m.Posix().Remove()) }
func (m mustPath) RemoveAll()                           { must(m.Posix().RemoveAll()) }
func (m mustPath) Touch(perm os.FileMode, existOk bool) { must(m.Posix().Touch(perm, existOk)) }
func (m mustPath) Pure() PurePath                       { return must2posix(m).Pure() }
func (m mustPath) Posix() PosixPath                     { return NewPosixPath(string(m)) }
func (m mustPath) Must() MustActions                    { return m }
func (m mustPath) String() string                       { return string(m) }

func (m mustPath) Match(pattern string) bool {
	b, e := must2posix(m).Match(pattern)
	must(e)
	return b
}

func (m mustPath) ExMatch(pattern string) bool {
	b, e := must2posix(m).ExMatch(pattern)
	must(e)
	return b
}
