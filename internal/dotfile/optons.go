package dotfile

type Settings struct {
	Prefix      string
	OnConflict  OnConflict
	DryRun      bool
	SourcePaths []string
	DestPath    string
}
