package dotfile

import (
	fp "path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	ppath "github.com/slyphon/dfi/pkg/pathlib"
)

type (
	ApplyFn func(ld LinkData) error

	Installer struct {
		prefix     string
		onConflict OnConflict
		apply      ApplyFn
	}

	// for testing, collects the LinkData Run calls us with
	applyCollector struct {
		links []LinkData
		apply ApplyFn
	}
)

func newApplyCollector() *applyCollector {
	ac := &applyCollector{}
	ac.apply = func(ls LinkData) error {
		ac.links = append(ac.links, ls)
		return nil
	}

	return ac
}

func runApply(ld LinkData, conflict OnConflict) (err error) {
	var fn func () error

	fn = func () error {
		vpath := ppath.NewPosixPath(ld.Vpath)
		lpath := ppath.NewPosixPath(ld.LinkPath)

		ctx := log.WithFields(log.Fields{
			"Vpath": ld.Vpath,
			"LinkPath": ld.LinkPath,
			"LinkData": ld.LinkData,
		})

		// if the path doesn't exist or it's a symlink
		if !lpath.Exists() || lpath.IsSymlink() {
			// if the path isn't a symlink we can create it and return
			if !lpath.IsSymlink() {
				return lpath.SymlinkTo(ld.LinkData)
			}

			// the path is a symlink, so we have to figure out what to do

			ctx.Debug("found symlink")

			// if the symlink points to the versioned file, we're done
			if same, err := lpath.SameFile(vpath); err != nil || same {
				return err
			}

			// otherwise it's bad, and we delegate to the onConflict.handler
			// to tell us what to do
			switch skip, err := conflict.handle(lpath.String()); {
			case err != nil:
				return err
			case skip: 	// the handler wants us to ignore this path
				return nil
			default:    // the handler (re)moved the lpath, so try again
				return fn()
			}
		} else if lpath.IsFile() || lpath.IsDir() {
			switch skip, err := conflict.handle(lpath.String()); {
			case err != nil:
				return err
			case skip:
				return nil
			default:
				return fn()
			}
		} else {
			ctx.Panic("could not handle conflict")
		}

		return nil
	}

	return fn()
}

func dryRunApply(ld LinkData) error {
	return nil
}

func NewDryRunInstaller(prefix string, onConflict OnConflict) *Installer {
	var apply = dryRunApply

	return &Installer{prefix, onConflict, apply}
}

func NewInstaller(prefix string, onConflict OnConflict) *Installer {
	var applyFn ApplyFn
	applyFn = func (ld LinkData) error {
		return runApply(ld, onConflict)
	}
	return &Installer{prefix, onConflict, applyFn}
}

func (n *Installer) Run(sourcePaths []string, destPath string) (err error) {
	var src []string
	var dst string

	if src, err = mkAbs(sourcePaths); err != nil {
		return err
	}

	if dst, err = fp.Abs(destPath); err != nil {
		return errors.Wrapf(err, "failed to Abs(%#v)", destPath)
	}

	var linkData []LinkData
	if linkData, err = LinkDataForList(src, dst, n.prefix); err != nil {
		return err
	}

	for _, ld := range linkData {
		if err = n.apply(ld); err != nil {
			return err
		}
	}

	return nil
}

func ApplySettings(s *Settings) error {
	return NewInstaller(s.Prefix, s.OnConflict).Run(s.SourcePaths, s.DestPath)
}

