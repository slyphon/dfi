package dotfile

import (
	"github.com/pkg/errors"
	"os"
	fp "path/filepath"
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


func exists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, errors.Wrapf(err, "unexpected error doing Lstat on %#v", path)
	}
}

func dryRunApply(ld LinkData) error {
	//fn := func () error {
	//	fp.Dir(ld.LinkPath)
	//}
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
