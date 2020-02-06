package dotfile

import (
	log "github.com/sirupsen/logrus"
)

func mkLinkData(s *Settings) ([]LinkData, error) {
	var err error
	ld := make([]LinkData, len(s.SourcePaths))

	for i, sp := range s.SourcePaths {
		if ld[i], err = LinkDataFor(sp, s.DestPath, s.Prefix); err != nil {
			return nil, err
		}
	}

	return ld, nil
}

func DryRun(s Settings) error {
	var err error
	var stg *Settings

	if stg, err = s.AbsPaths(); err != nil {
		return err
	}

	var linkData []LinkData

	if linkData, err = mkLinkData(stg); err != nil {
		return err
	}

	for _, ld := range linkData {
		log.WithFields(log.Fields{
			"Vpath":    ld.Vpath,
			"LinkPath": ld.LinkPath,
			"LinkData": ld.LinkData,
		}).Debug("LinkData")
	}

	return nil
}
