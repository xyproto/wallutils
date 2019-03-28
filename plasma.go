package wallutils

import (
	"fmt"
)

// Plasma windowmanager detector
type Plasma struct {
	mode    string
	verbose bool
}

func (p *Plasma) Name() string {
	return "Plasma"
}

func (p *Plasma) ExecutablesExists() bool {
	return (which("dbus-send") != "") && ((which("kwin_x11") != "") || (which("kwin_wayland") != "") || (which("kwin") != ""))
}

func (p *Plasma) Running() bool {
	return containsE("GDMSESSION", "plasma") || containsE("XDG_SESSION_DESKTOP", "KDE") || containsE("XDG_CURRENT_DESKTOP", "KDE") || containsE("XDG_SESSION_DESKTOP", "plasma") || containsE("XDG_CURRENT_DESKTOP", "plasma")
}

func (p *Plasma) SetMode(mode string) {
	p.mode = mode
}

func (p *Plasma) SetVerbose(verbose bool) {
	p.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (p *Plasma) SetWallpaper(imageFilename string) error {
	if !exists(imageFilename) {
		return fmt.Errorf("no such file: %s", imageFilename)
	}

	mode := defaultMode

	// If p.mode is specified, do not use the default value
	if p.mode != "" {
		mode = p.mode
	}

	fillMode := "0"
	if len(mode) == 1 {
		// Single digit
		fillMode = mode
	} else {
		// Drawing inspiration from https://github.com/KDE/plasma-workspace/blob/master/wallpapers/image/imagepackage/contents/ui/config.qml
		switch mode {
		case "crop":
			// Image.PreserveAspectCrop
			fillMode = "0"
		case "stretch":
			// Image.Stretch
			fillMode = "1"
		case "scale", "scaled":
			// Image.PreserveAaspectFit
			fillMode = "2"
		case "center", "centered":
			// Image.Pad
			fillMode = "3"
		case "tile", "tiled":
			// Image.Tile
			fillMode = "4"
		case "fill":
			fillMode = "5"
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
            d.writeConfig("Image", "file://` + imageFilename + `");
            d.writeConfig("FillMode", ` + fillMode + `);
    }`
	return run("dbus-send", []string{
		"--session",
		"--dest=org.kde.plasmashell",
		"--type=method_call",
		"/PlasmaShell",
		"org.kde.PlasmaShell.evaluateScript",
		dbusScript}, p.verbose)
}
