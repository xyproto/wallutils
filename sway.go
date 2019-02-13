package monitor

// Sway windowmanager detector
type Sway struct {
	verbose bool
}

func (s *Sway) Name() string {
	return "Sway"
}

func (s *Sway) ExecutablesExists() bool {
	return which("sway") != "" && which("swaymsg") != ""
}

func (s *Sway) Running() bool {
	return hasE("SWAYSOCK") && (containsE("GDMSESSION", "sway") || containsE("XDG_SESSION_DESKTOP", "sway") || containsE("XDG_CURRENT_DESKTOP", "sway"))
}

func (s *Sway) SetVerbose(verbose bool) {
	s.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (s *Sway) SetWallpaper(imageFilename string) error {
	return run("swaymsg 'output \"*\" background "+imageFilename+" fill'", s.verbose)
}
