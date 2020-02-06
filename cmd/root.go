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
  "github.com/spf13/cobra"
  log "github.com/sirupsen/logrus"

	homedir "github.com/mitchellh/go-homedir"
	df "github.com/slyphon/dfi/internal/dotfile"
	"github.com/spf13/viper"
)


var (
  cfgFile string
  conflictOpt string

  settings df.Settings = df.Settings{
		Prefix:      "",
		OnConflict:  df.Backup,
		DryRun:      false,
		SourcePaths: nil,
    DestPath:    "",
  }

  // rootCmd represents the base command when called without any subcommands
  rootCmd = &cobra.Command{
    Use:   "dfi sources... dest",
    Short: "Manages dotfile symlinks to version-controlled files",
    Long: `Usage: dfi [flags] sources... dest

Sources should be a list of (non-recursive) globs, files, and directories whose
contents should be symlinked into dest. The contents of directories are noted
with a trailing '/' on the path. To link a directory itself, omit the trailing
slash.

SOURCES:

    foo*    -> files beginning with 'foo'
    a/dir/  -> contents of 'a/dir'
    a/dir   -> create a link to 'dir' itself in dest

For info on the supported glob syntax, see https://github.com/gobwas/glob

MOTIVATION:

The purpose of this utility is to keep configuration files and directories
under source control in a single directory, then symlink them into place.
This program will by default look in the current working directory for a
directory named 'dotfiles', and will create symlinks to all of them in the
parent of the current directory with a '.' prepended.

$HOME/
  .settings/
    dotfiles/
      bashrc
      bash_profile
      zshrc
      ssh/
        config

if you cd into '~/.settings' and run dotinstall it will create the following symlinks:

$HOME/
  .bashrc -> .settings/dotfiles/bashrc
  .bash_profile -> .settings/dotfiles/bash_profile
  .zshrc -> .settings/dotfiles/zshrc
  .ssh -> .settings/dotfiles/ssh

Additionally, it can install links under a 'bin' directory, where the '.' prefix
is not applied. This is useful when you have a number of shell scripts and want
to link them from your version controllled directory into a location in your PATH.

$HOME/
  .settings/
    bin/
      foo
      bar
      baz

can be linked to

$HOME/
  .local/
    bin/
      foo -> ../../.settings/bin/foo
      bar -> ../../.settings/bin/bar
      baz -> ../../.settings/bin/baz

In the case of conflicts (i.e. destination already exists) you can decide
how files and symlinks will be handled.

If a link path already exists, the following strategies are available:

* 'backup': move the file to a unique dated backup location and create the symlink

* 'replace': just delete the file and create the symlink

* 'warn': print a warning that the conflict exists and continue.

* 'fail': stop processing and report an error.

EXAMPLES:

create symlinks for all files in 'dotfiles' in the home directory, with a '.' prefix
for each symlink name. i.e. ~/.bashrc -> .settings/dotfiles/bashrc

  $ dfi -p . ~/.settings/dotfiles/ ~/

create symlinks for all files in 'bin' under ~/.local/bin.
eg. ~/.local/bin/foo -> ../../.settings/bin/foo

  $ dfi ~/.settings/bin/ ~/.local/bin

`,
    Args: cobra.MinimumNArgs(2),
    // Uncomment the following line if your bare application
    // has an action associated with it:
    RunE: func(cmd *cobra.Command, args []string) (err error) {
      if settings.OnConflict, err = df.OnConflictForString(conflictOpt); err != nil {
        return err
      }

      settings.DestPath = args[len(args)-1]

      if err = destIsDir(settings.DestPath); err != nil {
        return err
      }

      settings.SourcePaths = args[0:len(args)-1]

      log.Infof("parsed settings: %+v", settings)

      return nil
    },
  }
)

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
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.

  rootCmd.PersistentFlags().StringVarP(&settings.Prefix, "prefix", "p", "", "A prefix to put before link names, eg. dotfiles have a '.' prefix")
  rootCmd.PersistentFlags().StringVarP(&conflictOpt, "on-conflct", "C", "backup", "Action to take when the symlink location exists: backup, replace, warn, fail")
  rootCmd.PersistentFlags().BoolVarP(&settings.DryRun, "dry-run", "n", false, "Show what would be done, take no action")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".dfi-go" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".dfi")
  }

  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}

