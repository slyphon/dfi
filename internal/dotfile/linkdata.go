package dotfile

import (
		fp "path/filepath"
		"github.com/pkg/errors"
		log "github.com/sirupsen/logrus"
)

const SLASH uint8 = 0x2f

func FindCommonRoot(a string, b string) string {
	var lastSlash int = 0

	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			break
		} else if a[i] == SLASH {
			lastSlash = i
		}
	}

	return a[0:lastSlash]
}

type LinkData struct {
	// Bpath is the "versioned" path where the config file resides
	Vpath string

	// LinkPath is the location of the symlink we install
	LinkPath string

	// LinkData is the contents of the symlink
	LinkData string
}

var emptyLinkData = LinkData{Vpath: "", LinkPath: "", LinkData: ""}

func LinkDataForList(vpaths []string, targetDir string, prefix string) (data []LinkData, err error) {
	data = make([]LinkData, len(vpaths))

	for i, vp := range vpaths {
		if data[i], err = LinkDataFor(vp, targetDir, prefix); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func LinkDataFor(vpath string, targetDir string, prefix string) (LinkData, error) {
	linkPath := fp.Join(targetDir, prefix + fp.Base(vpath))

	common := FindCommonRoot(vpath, linkPath)
	// we found a common root, now relativize the link data to point at the
	// versioned file
	if common != "" {
		// logging context for errors
		ctx := log.WithFields(log.Fields{
			"vpath": vpath,
			"linkPath": linkPath,
			"targetDir": targetDir,
			"common": common,
		})

		rel, err := fp.Rel(targetDir, common)
		if err != nil {
			ctx.Error("failed to relativize targetDir with common")
			return emptyLinkData, errors.WithStack(err)
		}

		vpRel, err := fp.Rel(common, vpath)
		if err != nil {
			ctx.Error("failed to relativize common with vpath")
			return emptyLinkData, errors.WithStack(err)
		}

		return LinkData{
			Vpath: vpath,
			LinkPath: linkPath,
			LinkData: fp.Join(rel, vpRel),
		}, nil
	}

	// no common path, just use an abspath
	return LinkData{
		Vpath: vpath,
		LinkPath: linkPath,
		LinkData: vpath,
	}, nil
}
