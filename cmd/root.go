/*
Copyright © 2020 Jonathan Simms <slyphon@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	df "github.com/slyphon/dfi/internal/dotfile"
)

func NewRootCommand() (rootCmd *cobra.Command) {
	conflictOpt := ""
	settings := &df.Settings{}

	rootCmd = &cobra.Command{
		Use:   "dfi sources... dest",
		Short: "Manages dotfile symlinks to version-controlled files",
		Long: `Usage: dfi [flags] sources... dest

Sources should be a list of paths that should have symlinks created in
dest. Note that the case of duplicate filenames (which would create
two sources to the same symlink) is considered an error. Sources can
also be '-' which means to read source paths from stdin, one per line,
or if the -z flag is given, separated by null bytes.

In the case of conflicts (i.e. destination already exists) you can decide
how files and symlinks will be handled.

If a link path already exists, the following strategies are available:

* 'rename': move the file to a unique dated location and create the symlink

* 'replace': just delete the file and create the symlink

* 'warn': print a warning that the conflict exists and continue.

* 'fail': stop processing and report an error.

`,
		Args: cobra.MinimumNArgs(2),

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if settings.OnConflict, err = df.OnConflictForString(conflictOpt); err != nil {
				return err
			}

			settings.DestPath = args[len(args)-1]

			if err = destIsDir(settings.DestPath); err != nil {
				return err
			}

			settings.SourcePaths = args[0 : len(args)-1]

			log.Tracef("parsed settings: %+v", settings)

			return df.ApplySettings(settings)
		},
	}

	rootCmd.PersistentFlags().StringVarP(
		&settings.Prefix,
		"prefix", "p", "",
		"A prefix to put before link names, eg. dotfiles have a '.' prefix",
	)

	rootCmd.PersistentFlags().StringVarP(
		&conflictOpt,
		"on-conflct", "C",
		"rename",
		"Action to take when the symlink location exists: rename, replace, warn, fail",
	)

	return rootCmd
}

func destIsDir(dest string) error {
	info, err := os.Stat(dest)
	if err != nil {
		return errors.Wrapf(err, "failed to stat dest: %v", dest)
	}
	if !info.Mode().IsDir() {
		return errors.Errorf("dest was not a directory: %v", dest)
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(rootCmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
