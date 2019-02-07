package monitor

import (
	"errors"
	"fmt"
	"strings"
)

// Xfce4 windowmanager detector
type Xfce4 struct {
}

func (s *Xfce4) Name() string {
	return "Xfce4"
}

func (s *Xfce4) ExecutablesExists() bool {
	return (which("xfconf-query") != "") && (which("xfce4-session") != "")
}

func (s *Xfce4) Running() bool {
	return (containsE("GDMSESSION", "xfce") || containsE("XDG_SESSION_DESKTOP", "xfce") || containsE("DESKTOP_SESSION", "xfce"))
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (s *Xfce4) SetWallpaper(imageFilename string) error {
	properties := strings.Split(output("xfconf-query --channel xfce4-desktop --list"), "\n")
	if len(properties) == 0 {
		return errors.New("Could not list any properties for Xfce4")
	}
	for _, prop := range properties {
		if strings.HasSuffix(prop, "/last-image") {
			if err := run(fmt.Sprintf("xfconf-query --channel xfce4-desktop --property %s --set %q", prop, imageFilename)); err != nil {
				return err
			}
		}
	}
	return nil
}
