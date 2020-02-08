package pathlib

type (
	FileInfoAdditions interface {
		IsBlockDevice() bool
		IsCharDevice() bool
		IsFifo() bool
		IsFile() bool
		IsSocket() bool
		IsSymlink() bool
	}

	Posixly interface {
		Posix() PosixPath
	}

	Purely interface {
		Pure() PurePath
	}

	Dangerously interface {
		Must() MustActions
	}

	FlavorCasts interface {
		Posixly
		Purely
		Dangerously
	}
)

func pure2posix(p pureStr) posixStr  { return posixStr(string(p)) }
func posix2pure(p posixStr) pureStr  { return pureStr(string(p)) }
func posix2must(p posixStr) mustPath { return mustPath(string(p)) }
func must2posix(p mustPath) posixStr { return posixStr(string(p)) }

func PosixSliceStringer(ps []PosixPath) (str []string) {
	str = make([]string, len(ps))
	for i, pp := range ps {
		str[i] = pp.String()
	}
	return
}

