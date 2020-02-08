package pathlib

import (
	"fmt"
	"os"
	"time"
)

type (
	RichFileInfo interface {
		os.FileInfo
		FileInfoAdditions
		String() string
		getInfo() os.FileInfo
	}
	richInfo struct {
		info os.FileInfo
	}
)

var _ RichFileInfo = richInfo{nil}
var _ os.FileInfo = richInfo{nil}

func (r richInfo) getInfo() os.FileInfo { return r.info }
func (r richInfo) Name() string         { return r.info.Name() }
func (r richInfo) Size() int64          { return r.info.Size() }
func (r richInfo) Mode() os.FileMode    { return r.info.Mode() }
func (r richInfo) ModTime() time.Time   { return r.info.ModTime() }
func (r richInfo) IsDir() bool          { return r.info.IsDir() }
func (r richInfo) Sys() interface{}     { return r.info.Sys() }

func (r richInfo) IsBlockDevice() bool { return r.info.Mode()&os.ModeDevice != 0 }
func (r richInfo) IsCharDevice() bool  { return r.info.Mode()&os.ModeCharDevice != 0 }
func (r richInfo) IsFifo() bool        { return r.info.Mode()&os.ModeNamedPipe != 0 }
func (r richInfo) IsFile() bool        { return r.info.Mode()&os.ModeType == 0 }
func (r richInfo) IsSocket() bool      { return r.info.Mode()&os.ModeSocket != 0 }
func (r richInfo) IsSymlink() bool     { return r.info.Mode()&os.ModeSymlink != 0 }

func (r richInfo) String() string {
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
