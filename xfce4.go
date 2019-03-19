package wallutils

import (
	"errors"
	"strings"
)

// Xfce4 windowmanager detector
type Xfce4 struct {
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

func (x *Xfce4) SetVerbose(verbose bool) {
	x.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (x *Xfce4) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return errors.New(imageFilename + " does not exist")
	}
	properties := strings.Split(output("xfconf-query", []string{"--channel", "xfce4-desktop", "--list"}, x.verbose), "\n")
	if len(properties) == 0 {
		return errors.New("Could not list any properties for Xfce4")
	}
	for _, prop := range properties {
		if strings.HasSuffix(prop, "/last-image") {
			if err := run("xfconf-query", []string{"--channel", "xfce4-desktop", "--property", prop, "--set", imageFilename}, x.verbose); err != nil {
				return err
			}
		}
	}
	return nil
}
