package monitor

import (
	"errors"
	"fmt"
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
	command := "xfconf-query --channel xfce4-desktop --list"
	properties := strings.Split(output(command, x.verbose), "\n")
	if len(properties) == 0 {
		return errors.New("Could not list any properties for Xfce4")
	}
	for _, prop := range properties {
		if strings.HasSuffix(prop, "/last-image") {
			command = fmt.Sprintf("xfconf-query --channel xfce4-desktop --property %s --set %q", prop, imageFilename)
			if err := run(command, x.verbose); err != nil {
				return err
			}
		}
	}
	return nil
}
