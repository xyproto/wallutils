package monitor

// Deepin windowmanager detector
type Deepin struct {
	verbose bool
}

func (d *Deepin) Name() string {
	return "Deepin"
}

func (d *Deepin) ExecutablesExists() bool {
	return which("deepin-session") != "" && which("dconf") != ""
}

func (d *Deepin) Running() bool {
	return containsE("GDMSESSION", "deepin") || containsE("XDG_SESSION_DESKTOP", "deepin") || containsE("XDG_CURRENT_DESKTOP", "deepin")
}

func (d *Deepin) SetVerbose(verbose bool) {
	d.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (d *Deepin) SetWallpaper(imageFilename string) error {
	return run("dconf write /com/deepin/wrap/gnome/desktop/background/picture-uri \"'"+imageFilename+"'\"", d.verbose)
}
