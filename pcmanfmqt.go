package wallutils

import (
	"fmt"
	"os"
)

// PCManFMQt windowmanager detector
type PCManFMQt struct {
	mode    string
	verbose bool
}

// Name returns the name of this window manager or desktop environment
func (pcmq *PCManFMQt) Name() string {
	return "PCManFM-Qt"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (pcmq *PCManFMQt) ExecutablesExists() bool {
	return which("pcmanfm-qt") != ""
}

// Running examines environment variables to try to figure out if this backend is currently running
func (pcmq *PCManFMQt) Running() bool {
	// Detect LxQt
	return (os.Getenv("XDG_MENU_PREFIX") == "lxqt-") || (os.Getenv("DESKTOP_SESSION") == "lxqt")
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (pcmq *PCManFMQt) SetMode(mode string) {
	pcmq.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (pcmq *PCManFMQt) SetVerbose(verbose bool) {
	pcmq.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (pcmq *PCManFMQt) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	// initialize the mode setting (stretched/tiled etc)
	mode := defaultMode
	if pcmq.mode != "" {
		mode = pcmq.mode
	}

	switch mode {
	case "none":
		break
	case "color":
		mode = "color"
	case "centered", "center":
		mode = "center"
	case "zoom":
		mode = "zoom"
	case "stretch", "stretched":
		mode = "stretch"
	case "fill", "scale", "scaled", "spanned":
		mode = "fit"
	case "wallpaper", "tile":
		mode = "tile"
	default:
		// Invalid and unrecognized desktop wallpaper mode
		return fmt.Errorf("invalid desktop wallpaper mode for PCManFM-Qt: %s", mode)
	}

	// Set the wallpaper image with the selected mode
	return run("pcmanfm-qt", []string{"--wallpaper-mode", mode, "--set-wallpaper", imageFilename}, pcmq.verbose)
}
