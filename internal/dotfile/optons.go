package dotfile

import (
	"github.com/pkg/errors"
	fp "path/filepath"
)

type Settings struct {
	Prefix      string
	OnConflict  OnConflict
	DryRun      bool
	SourcePaths []string
	DestPath    string
}

func mkAbs(paths []string) ([]string, error) {
	abs := make([]string, len(paths))

	for i, p := range paths {
		ap, err := fp.Abs(p)
		if e := errors.Wrap(err, "Abs on path failed"); e != nil {
			return nil, e
		}
		abs[i] = ap
	}

	return abs, nil
}

func mkAbsSettings(s *Settings) (*Settings, error) {
	var err error
	if s.SourcePaths, err = mkAbs(s.SourcePaths); err != nil {
		return nil, err
	}
	if s.DestPath, err = fp.Abs(s.DestPath); err != nil {
		return nil, err
	}
	return s, nil
}

// AbsPaths returns a reference to a copy of the receiver with
// SourcePaths and DestPaths converted to absolute paths
//
func (s Settings) AbsPaths() (*Settings, error) {
	return mkAbsSettings(&s)
}
