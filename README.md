The purpose of this utility is to keep configuration files and directories under source control in a single directory, then symlink them into place. This is useful for lightweight deployments of personal dev environment settings, and for keeping them in sync.

`dfi` is written in classic UNIX style, in that it takes a list of files, a destination, and creates symlinks, that's all.

```
$HOME/
  .settings/
    dotfiles/
      bashrc
      bash_profile
      zshrc
      ssh/
        config
```

if you run `dfi --prefix=. .settings/dotfiles/* ./` in your home directory, it will create the following symlinks:

```
$HOME/
  .bashrc -> .settings/dotfiles/bashrc
  .bash_profile -> .settings/dotfiles/bash_profile
  .zshrc -> .settings/dotfiles/zshrc
  .ssh -> .settings/dotfiles/ssh
```

Additionally, it can install links under a `bin` directory, where the `.` prefix
is not applied. This is useful when you have a number of shell scripts and want
to link them from your version controllled directory into a location in your PATH.

Given:

```
$HOME/
  .settings/
    bin/
      foo
      bar
      baz
```

In your HOME directory, if you run `dfi .settings/bin/* .local/bin/`, will result in the following links being created.

```
$HOME/
  .local/
    bin/
      foo -> ../../.settings/bin/foo
      bar -> ../../.settings/bin/bar
      baz -> ../../.settings/bin/baz
```

In the case of conflicts (i.e. destination already exists) you can decide
how files and symlinks will be handled.

If a link path already exists, the following strategies are available:

* `rename`: move the file to a unique dated location and create the symlink

* `replace`: just delete the file and create the symlink

* `warn`: print a warning that the conflict exists and continue.

* `fail`: stop processing and report an error.


Note that `dfi` will never `replace` a directory (i.e. `rm -rf` it), rather that is treated as an error and execution will halt with a non-zero return code.

