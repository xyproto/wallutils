package monitor

// Plasma windowmanager detector
type Plasma struct {
	verbose bool
}

func (p *Plasma) Name() string {
	return "Plasma"
}

func (p *Plasma) ExecutablesExists() bool {
	return (which("dbus-send") != "") && ((which("kwin_x11") != "") || (which("kwin_wayland") != "") || (which("kwin") != ""))
}

func (p *Plasma) Running() bool {
	return containsE("GDMSESSION", "plasma") || containsE("XDG_SESSION_DESKTOP", "plasma") || containsE("XDG_CURRENT_DESKTOP", "plasma")
}

func (p *Plasma) SetVerbose(verbose bool) {
	p.verbose = verbose
}

// SetWallpaper sets the desktop wallpaper, given an image filename.
// The image must exist and be readable.
func (p *Plasma) SetWallpaper(imageFilename string) error {
	return run(`dbus-send --session --dest=org.kde.plasmashell --type=method_call /PlasmaShell org.kde.PlasmaShell.evaluateScript 'string:
    var Desktops = desktops();
    for (i=0;i<Desktops.length;i++) {
            d = Desktops[i];
            d.wallpaperPlugin = "org.kde.image";
            d.currentConfigGroup = Array("Wallpaper",
                                         "org.kde.image",
                                         "General");
            d.writeConfig("Image", "file://`+imageFilename+`");
    }'`, p.verbose)
}
