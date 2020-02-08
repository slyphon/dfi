package dotfile

func dryRunApply(ld LinkData) error {
	return nil
}

func NewDryRunInstaller(prefix string, onConflict OnConflict) *Installer {
	var apply = dryRunApply

	return &Installer{prefix, onConflict, apply}
}

