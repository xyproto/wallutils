package wallutils

import (
	"errors"
	"fmt"
	"strings"
)

// Xfce4 windowmanager detector
type Xfce4 struct {
	mode    string
	verbose bool
}

func (x *Xfce4) Name() string {
	return "Xfce4"
}

func (x *Xfce4) ExecutablesExists() bool {
	return (which("xfconf-query") != "") && (which("xfce4-session") != "")
}

func (x *Xfce4) Running() bool {
	return (containsE("GDMSESSION", "xfce") || containsE("XDG_SESSION_DESKTOP", "xfce") || containsE("DESKTOP_SESSION", "xfce"))
}

func (x *Xfce4) SetMode(mode string) {
	x.mode = mode
}

func (x *Xfce4) SetVerbose(verbose bool) {
	x.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (x *Xfce4) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	// Find a list of all available properties for all monitors
	properties := strings.Split(output("xfconf-query", []string{"--channel", "xfce4-desktop", "--list"}, x.verbose), "\n")
	if len(properties) == 0 {
		return errors.New("could not find any properties for Xfce4")
	}

	// initialize the mode setting (stretched/tiled etc)
	mode := defaultMode
	if x.mode != "" {
		mode = x.mode
	}

	// Wallpaper mode for Xfce4: Auto=0, Centered=1, Tiled=2, Stretched=3, Scaled=4, Zoomed=5
	fillMode := "3"
	if len(mode) == 1 {
		// Single digit
		fillMode = mode
	} else {
		switch mode {
		case "stretch", "stretched":
			fillMode = "3"
		case "auto":
			fillMode = "0"
		case "center", "centered":
			fillMode = "1"
		case "tile", "tiled":
			fillMode = "2"
		case "scale", "scaled", "fit", "fill":
			fillMode = "4"
		case "zoom", "zoomed", "crop", "cropped":
			fillMode = "5"
		default:
			// Invalid and unrecognized desktop wallpaper mode
			return fmt.Errorf("invalid desktop wallpaper mode for Xfce4: %s", x.mode)
		}
	}

	for _, prop := range properties {

		if strings.HasSuffix(prop, "/image-style") {
			if err := run("xfconf-query", []string{"--channel", "xfce4-desktop", "--property", prop, "--set", fillMode}, x.verbose); err != nil {
				return err
			}
		}
		if strings.HasSuffix(prop, "/last-image") {
			if err := run("xfconf-query", []string{"--channel", "xfce4-desktop", "--property", prop, "--set", imageFilename}, x.verbose); err != nil {
				return err
			}
		}
	}
	return nil
}
