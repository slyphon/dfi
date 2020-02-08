package fsfixture

import (
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	pl "github.com/slyphon/dfi/pkg/pathlib"
)

const (
	timeFormat      string = "20060102150405"
	TimeFormatFloat        = timeFormat + ".000000000"
)

func timestamp() string {
	return time.Now().Format(TimeFormatFloat)
}


const DirPerms os.FileMode = 0755

func mustTempDir(prefix string) string {
	td, err := ioutil.TempDir("", prefix)
	if err != nil {
		log.Panicf("Failed to create TempDir, err: %+v", err)
	}
	return td
}

type (
	FsFixture struct {
		TempDir     pl.PosixPath
		HomeDir     pl.PosixPath
		SettingsDir pl.PosixPath
		DotfileDir  pl.PosixPath
		BinDir      pl.PosixPath
		LocalBinDir pl.PosixPath
		Dotfiles    []pl.PosixPath
		Binfiles    []pl.PosixPath
	}
)


// layout:
//
// 	temp/
//  . home/
//  . .	settings/
//	.	.	.	dotfiles/
//  . .	. .	bashrc
//	.	.	.	.	vimrc
// 	.	.	.	.	zshrc
//  . . . . config/
//  . . . . . yarn/
//  . . . . . . global/
//  . . . . . . . package.json
//	.	.	.	bin/
//	.	.	.	.	cat
// 	.	.	.	. dog
//	.	.	.	.	ls
//	. . .local
//  .	.	.	bin/

var (
	BinfileNames = []string{"cat", "dog", "ls"}
	DotfileNames = []string{"bashrc", "config", "vimrc", "zshrc"}
)

func NewFsFixture() FsFixture {
	TempDir := pl.NewPosixPath(mustTempDir(timestamp()))
	HomeDir := TempDir.Join("home")
	SettingsDir := HomeDir.Join("settings")
	DotfileDir := SettingsDir.Join("dotfiles")
	BinDir := SettingsDir.Join("bin")
	LocalBinDir := HomeDir.Join(".local/bin")
	Dotfiles := make([]pl.PosixPath, 0, len(DotfileNames))
	Binfiles := make([]pl.PosixPath, 0, len(BinfileNames))

	dirs := []pl.PosixPath{TempDir, HomeDir, SettingsDir, DotfileDir, BinDir, LocalBinDir}

	for _, d := range dirs {
		d.Must().MkdirAll(DirPerms)
	}

	for _, f := range DotfileNames {
		p := DotfileDir.Join(f).Posix()
		p.Must().Touch(0o644, false /* existOk */)
		Dotfiles = append(Dotfiles, p)
	}

	for _, f := range BinfileNames {
		p := BinDir.Join(f).Posix()
		p.Must().Touch(0o644, false /* existOk */)
		Binfiles = append(Binfiles, p)
	}

	return FsFixture{
		TempDir, HomeDir, SettingsDir, DotfileDir, BinDir, LocalBinDir, Dotfiles, Binfiles,
	}
}

func (f FsFixture) Cleanup() {
	err := f.TempDir.RemoveAll()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("ignoring error in FsFixture.Cleanup")
	}
}
