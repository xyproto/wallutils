//go:build !cgo
// +build !cgo

package wallutils

// WMs contains all available backends for changing the wallpaper
// Only backends that do not require cgo should be included here.
var WMs = []WM{
	&Hyprpaper{},
	//&Sway{},
	&Deepin{},
	&Xfce4{},
	&Mate{},
	&Cinnamon{},
	&Plasma{},
	&Gnome3{},
	&Gnome2{},
	&Pekwm{},
	&Weston{},
	// xbg.New(), // X11
	&Feh{}, // use feh for X11
}
