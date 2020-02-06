package fsfixture

import (
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	fp "path/filepath"

	df "github.com/slyphon/dfi/internal/dotfile"
)

const (
	TimeFormatFloat string = df.TimeFormat + ".000000000"
)

func timestamp() string {
	return time.Now().Format(TimeFormatFloat)
}

func mustTouch(path ...string) string {
	joined := fp.Join(path...)
	var fp *os.File
	var err error

	if fp, err = os.Create(joined); err != nil {
		log.Panicf("failed to touch file %@v", joined)
	}

	if err = fp.Close(); err != nil {
		log.Panicf("failed to close file handle associated with %#v", joined)
	}

	return joined
}

const DirPerms os.FileMode = 0755

func mustMkDirAll(path ...string) string {
	joined := fp.Join(path...)
	err := os.MkdirAll(joined, DirPerms)
	if err != nil {
		log.Panicf("failed to mkdir %#v, err: %+v", path, err)
	}
	return joined
}

func mustTempDir(prefix string) string {
	td, err := ioutil.TempDir("", prefix)
	if err != nil {
		log.Panicf("Failed to create TempDir, err: %+v", err)
	}
	return td
}

type (
	FsFixture struct {
		TempDir     string
		HomeDir     string
		SettingsDir string
		DotfileDir  string
		BinDir      string
		LocalBinDir string
		Dotfiles    []string
		Binfiles    []string
	}
)

var (
	binfileNames = []string{"ls", "cat", "dog"}
	dotfileNames = []string{"bashrc", "zshrc", "vimrc"}
)

func NewFsFixture() FsFixture {
	TempDir := mustTempDir(timestamp())
	HomeDir := mustMkDirAll(TempDir, "home")
	SettingsDir := mustMkDirAll(HomeDir, "settings")
	DotfileDir := mustMkDirAll(SettingsDir, "dotfiles")
	BinDir := mustMkDirAll(DotfileDir, "bin")
	LocalBinDir := mustMkDirAll(HomeDir, ".local/bin")
	Dotfiles := make([]string, 0, len(dotfileNames))
	Binfiles := make([]string, 0, len(binfileNames))

	for _, f := range dotfileNames {
		Dotfiles = append(Dotfiles, mustTouch(DotfileDir, f))
	}

	for _, f := range binfileNames {
		Binfiles = append(Binfiles, mustTouch(BinDir, f))
	}

	return FsFixture{
		TempDir, HomeDir, SettingsDir, DotfileDir, BinDir, LocalBinDir, Dotfiles, Binfiles,
	}
}

func (f FsFixture) Cleanup() {
	err := os.RemoveAll(f.TempDir)
	log.WithFields(log.Fields{
		"err": err,
	}).Error("ignoring error in FsFixture.Cleanup")
}
