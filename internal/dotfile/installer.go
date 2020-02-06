package dotfile

import (
	"github.com/pkg/errors"
	fp "path/filepath"
)

type (
	Installer struct {
		prefix     string
		onConflict OnConflict
		apply      func(ld LinkData) error
	}
)

func dryRunApply(ld LinkData) error {
	return nil
}

func NewDryRunInstaller(prefix string, onConflict OnConflict) *Installer {
	var apply = dryRunApply

	return &Installer{prefix, onConflict, apply}
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
