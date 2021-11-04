package wallutils

import (
	"fmt"

	"github.com/xyproto/env"
)

// Plasma windowmanager detector
type Plasma struct {
	mode    string
	verbose bool
}

// Name returns the name of this window manager or desktop environment
func (p *Plasma) Name() string {
	return "Plasma"
}

// ExecutablesExists checks if executables associated with this backend exists in the PATH
func (p *Plasma) ExecutablesExists() bool {
	return (which("dbus-send") != "") && ((which("kwin_x11") != "") || (which("kwin_wayland") != "") || (which("kwin") != ""))
}

// Running examines environment variables to try to figure out if this backend is currently running
func (p *Plasma) Running() bool {
	return env.Contains("GDMSESSION", "plasma") || env.Contains("XDG_SESSION_DESKTOP", "KDE") || env.Contains("XDG_CURRENT_DESKTOP", "KDE") || env.Contains("XDG_SESSION_DESKTOP", "plasma") || env.Contains("XDG_CURRENT_DESKTOP", "plasma")
}

// SetMode will set the current way to display the wallpaper (stretched, tiled etc)
func (p *Plasma) SetMode(mode string) {
	p.mode = mode
}

// SetVerbose can be used for setting the verbose field to true or false.
// This will cause this backend to output information about what is is doing on stdout.
func (p *Plasma) SetVerbose(verbose bool) {
	p.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (p *Plasma) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	// initialize the mode setting (stretched/tiled etc)
	mode := defaultMode
	if p.mode != "" {
		mode = p.mode
	}

	var fillMode string
	if len(mode) == 1 {
		// Single digit
		fillMode = mode
	} else {
		// Drawing inspiration from https://github.com/KDE/plasma-workspace/blob/master/wallpapers/image/imagepackage/contents/ui/config.qml
		switch mode {
		case "stretch", "stretched":
			// stretch the picture to match the screen
			fillMode = "0"
		case "fill", "fit", "scale", "scaled":
			// fit and scale, but keep aspect ratio
			fillMode = "1"
		case "zoom", "zoomed", "crop", "cropped":
			// zoom
			fillMode = "2"
		case "tile", "tiled":
			// tiled
			fillMode = "3"
		case "hfill", "vtile":
			// fill horizontally, tile vertically
			fillMode = "4"
		case "vfill", "htile":
			// fill vertically, tile horizontally
			fillMode = "5"
		case "center", "centered":
			// center
			fillMode = "6"
		default:
			// Invalid and unrecognized desktop wallpaper mode
			return fmt.Errorf("invalid desktop wallpaper mode for Plasma: %s", p.mode)
		}
	}

	dbusScript := `string:
    var Desktops = desktops();
    for (i=0;i<Desktops.length;i++) {
            d = Desktops[i];
            d.wallpaperPlugin = "org.kde.image";
            d.currentConfigGroup = Array("Wallpaper",
                                         "org.kde.image",
                                         "General");
            d.writeConfig("FillMode", ` + fillMode + `);
            d.writeConfig("Image", "file://` + imageFilename + `");
    }`
	return run("dbus-send", []string{
		"--session",
		"--dest=org.kde.plasmashell",
		"--type=method_call",
		"/PlasmaShell",
		"org.kde.PlasmaShell.evaluateScript",
		dbusScript}, p.verbose)
}
