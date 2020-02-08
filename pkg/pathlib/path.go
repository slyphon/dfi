package pathlib

import (
	"fmt"
	"os"
	"path"
	fp "path/filepath"
	"syscall"
	"time"

	"github.com/gobwas/glob"
	log "github.com/sirupsen/logrus"
)

type (
	PurePath interface {
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
		String() string
		ToPosix() PosixPath
	}

	FileInfoAdditions interface {
		IsBlockDevice() bool
		IsCharDevice() bool
		IsFifo() bool
		IsFile() bool
		IsSocket() bool
		IsSymlink() bool
	}

	RichFileInfo interface {
		os.FileInfo
		FileInfoAdditions
		String() string
		getInfo() os.FileInfo
	}

	// Provide path manipulation like PurePath, but also
	// posix actions on the real filesystem. Note that the
	// FileInfoAdditions actions on a PosixPath will ignore all
	// errors and return false. This is the strategy that
	// the Python stdlib 'pathlib' follows in the interest of
	// ease of use. If you really need to know about every
	// error type, then use Stat() and inspect the error case before
	// interrogating the FileInfoAdditions method.
	PosixPath interface {
		PurePath
		FileInfoAdditions
		Stat() (RichFileInfo, error)
		Lstat() (RichFileInfo, error)
		SymlinkTo(path string) error
		Rel(other string) (PosixPath, error)
		Resolve() (PosixPath, error)
		Mkdir(perm os.FileMode) error
		MustMkdir(perm os.FileMode)
		MkdirAll(perm os.FileMode) error
		MustMkdirAll(perm os.FileMode)
		Remove() error
		RemoveAll() error
		Lexists() bool
		Exists() bool
		IsMount() bool
		SameFile(other PosixPath) (bool, error)
		IsDir() bool
		Glob(pattern string) ([]PosixPath, error)
	}

	pathStr string
)

var _ PurePath = pathStr("")
var _ fmt.Stringer = pathStr("")

func NewPurePath(s string) PurePath {
	return pathStr(s)
}

func NewPosixPath(s string) PosixPath {
	return NewPurePath(s).ToPosix()
}

func (p pathStr) Name() string {
	return path.Base(string(p))
}

func (p pathStr) Parent() PurePath {
	return pathStr(path.Dir(string(p)))
}

func (p pathStr) Join(names ...string) PurePath {
	els := append([]string{string(p)}, names...)
	return NewPurePath(path.Join(els...))
}

func (p pathStr) Clean() PurePath {
	return pathStr(path.Clean(string(p)))
}

func (p pathStr) Match(pattern string) (matched bool, err error) {
	return path.Match(pattern, string(p))
}

func (p pathStr) ExMatch(pattern string) (matched bool, err error) {
	g, err := glob.Compile(pattern, '/')
	if err != nil {
		return false, err
	}
	return g.Match(p.String()), nil
}

func (p pathStr) Split() (dir PurePath, name string) {
	d, n := path.Split(string(p))
	return pathStr(d), n
}

func (p pathStr) ToPosix() PosixPath {
	return PosixPath(p)
}

func (p pathStr) String() string {
	return string(p)
}

type richFileInfo struct {
	info os.FileInfo
}

var _ RichFileInfo = richFileInfo{nil}
var _ os.FileInfo = richFileInfo{nil}

func (r richFileInfo) getInfo() os.FileInfo {
	return r.info
}

func (r richFileInfo) Name() string {
	return r.info.Name()
}

func (r richFileInfo) Size() int64 {
	return r.info.Size()
}

func (r richFileInfo) Mode() os.FileMode {
	return r.info.Mode()
}

func (r richFileInfo) ModTime() time.Time {
	return r.info.ModTime()
}

func (r richFileInfo) IsDir() bool {
	return r.info.IsDir()
}

func (r richFileInfo) Sys() interface{} {
	return r.info.Sys()
}

func (r richFileInfo) IsBlockDevice() bool {
	return r.info.Mode()&os.ModeDevice != 0
}

func (r richFileInfo) IsCharDevice() bool {
	return r.info.Mode()&os.ModeCharDevice != 0
}

func (r richFileInfo) IsFifo() bool {
	return r.info.Mode()&os.ModeNamedPipe != 0
}

func (r richFileInfo) IsFile() bool {
	return r.info.Mode()&os.ModeType == 0
}

func (r richFileInfo) IsSocket() bool {
	return r.info.Mode()&os.ModeSocket != 0
}

func (r richFileInfo) IsSymlink() bool {
	return r.info.Mode()&os.ModeSymlink != 0
}

func (r richFileInfo) String() string {
	return fmt.Sprintf("%#v", struct {
		Name    string
		Mode    string
		Size    int64
		ModTime time.Time
		IsDir   bool
	}{
		Name:    r.Name(),
		Mode:    r.Mode().String(),
		Size:    r.Size(),
		ModTime: r.ModTime(),
		IsDir:   r.IsDir(),
	})
}

func (p pathStr) Stat() (RichFileInfo, error) {
	info, err := os.Stat(string(p))
	return richFileInfo{info}, err
}

func (p pathStr) Lstat() (RichFileInfo, error) {
	info, err := os.Lstat(string(p))
	return richFileInfo{info}, err
}

func (p pathStr) SymlinkTo(path string) error {
	return os.Symlink(path, string(p))
}

func (p pathStr) Rel(other string) (PosixPath, error) {
	pp, err := fp.Rel(string(p), other)
	if err != nil {
		return nil, err
	}
	return PosixPath(pathStr(pp)), nil
}

func (p pathStr) Resolve() (PosixPath, error) {
	pp, err := fp.EvalSymlinks(string(p))
	if err != nil {
		return nil, err
	}
	return PosixPath(pathStr(pp)), nil
}

func (p pathStr) Mkdir(perm os.FileMode) error {
	return os.Mkdir(string(p), perm)
}

func (p pathStr) MustMkdir(perm os.FileMode) {
	err := p.Mkdir(perm)
	if err != nil {
		log.Panicf("error in Mkdir for %v, perms: %v, err: %+v", p.String(), perm.String(), err.Error())
	}
}

func (p pathStr) MkdirAll(perm os.FileMode) error {
	return os.MkdirAll(string(p), perm)
}

func (p pathStr) MustMkdirAll(perm os.FileMode) {
	err := p.MkdirAll(perm)
	if err != nil {
		log.Panicf("error in MkdirAll for %v, perms: %v, err: %+v", p.String(), perm.String(), err.Error())
	}
}

func (p pathStr) Remove() error {
	return os.Remove(string(p))
}

func (p pathStr) RemoveAll() error {
	return os.RemoveAll(string(p))
}

func (p pathStr) Lexists() bool {
	_, err := p.Lstat()
	return err == nil
}

func (p pathStr) Exists() bool {
	_, err := p.Stat()
	return err == nil
}

func getStatT(f os.FileInfo) (*syscall.Stat_t, error) {
	statT, ok := f.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("failed to get system Stat_t")
	}
	return statT, nil
}

func (p pathStr) IsMount() bool {
	var self os.FileInfo
	var selfStatT *syscall.Stat_t
	var err error

	if self, err = p.Stat(); err != nil || !self.IsDir() {
		return false
	}

	if selfStatT, err = getStatT(self); err != nil {
		return false
	}

	var myParentInfo os.FileInfo

	if myParentInfo, err = p.Parent().ToPosix().Stat(); err != nil {
		return false
	}

	parentStatT, err := getStatT(myParentInfo)
	if err != nil {
		return false
	}

	// from linux coreutils:
	// int is_not_mnt = (st_dev == st.st_dev) && (st_ino != st.st_ino);
	//
	isNotMnt := (selfStatT.Dev == parentStatT.Dev) && (selfStatT.Ino != parentStatT.Ino)

	return !isNotMnt
}

func (p pathStr) SameFile(other PosixPath) (b bool, err error) {
	this, err := p.Stat()
	if err != nil {
		return false, err
	}
	that, err := other.Stat()
	if err != nil {
		return false, err
	}
	return os.SameFile(this.getInfo(), that.getInfo()), nil
}

func (p pathStr) Glob(pattern string) ([]PosixPath, error) {
	matches, err := fp.Glob(fp.Join(p.String(), pattern))
	if err != nil {
		return nil, err
	}
	paths := make([]PosixPath, 0, len(matches))
	for _, m := range matches {
		paths = append(paths, NewPosixPath(m))
	}
	return paths, err
}

func (p pathStr) IsBlockDevice() bool {
	info, err := p.Lstat()
	return err == nil && info.IsBlockDevice()
}

func (p pathStr) IsCharDevice() bool {
	info, err := p.Lstat()
	return err == nil && info.IsCharDevice()
}

func (p pathStr) IsFifo() bool {
	info, err := p.Lstat()
	return err == nil && info.IsFifo()
}

func (p pathStr) IsFile() bool {
	info, err := p.Lstat()
	return err == nil && info.IsFile()
}

func (p pathStr) IsSocket() bool {
	info, err := p.Lstat()
	return err == nil && info.IsSocket()
}

func (p pathStr) IsSymlink() bool {
	info, err := p.Lstat()
	return err == nil && info.IsSymlink()
}

func (p pathStr) IsDir() bool {
	info, err := p.Lstat()
	return err == nil && info.IsDir()
}

type LexOrder []PosixPath
func (o LexOrder) Len() int           { return len(o) }
func (o LexOrder) Less(i, j int) bool { return o[i].String() < o[j].String() }
func (o LexOrder) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
