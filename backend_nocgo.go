// +build !cgo

package wallutils

// WMs contains all available backends for changing the wallpaper
// Only backends that do not require cgo should be included here.
var WMs = []WM{
	//&Sway{},
	&Deepin{},
	&Xfce4{},
	&Mate{},
	&Cinnamon{},
	&Plasma{},
	&Gnome3{},
	&Gnome2{},
	&Weston{},
	&Feh{}, // using feh
	//&X11{},
}
