//go:build cgo
// +build cgo

package wallutils

// WMs contains all available backends for changing the wallpaper
// Some backends may require cgo (sway + x11)
var WMs = []WM{
	&Hyprpaper{},
	&Sway{},
	&Deepin{},
	&Xfce4{},
	&Mate{},
	&Cinnamon{},
	&Plasma{},
	&Gnome3{},
	&Gnome2{},
	&Pekwm{},
	&SwayBG{},
	&Weston{},
	// xbg.New(), // X11
	&Feh{}, // use feh for X11, for now
}

// Info returns a long info string that looks different for Wayland and for X.
// The string contains all available information about the connected monitors.
func Info() (string, error) {
	if WaylandCanConnect() {
		return WaylandInfo()
	} else if XCanConnect() {
		return XInfo()
	}
	return "", errNoWaylandNoX
}

// Monitors returns information about all monitors, regardless of if it's under
// Wayland or X11. Will use additional plugins, if available.
func Monitors() ([]Monitor, error) {
	IDs, widths, heights, wDPIs, hDPIs := []uint{}, []uint{}, []uint{}, []uint{}, []uint{}
	if WaylandCanConnect() {
		if err := WaylandMonitors(&IDs, &widths, &heights, &wDPIs, &hDPIs); err != nil {
			return []Monitor{}, err
		}
	} else if XCanConnect() {
		if err := XMonitors(&IDs, &widths, &heights, &wDPIs, &hDPIs); err != nil {
			return []Monitor{}, err
		}
	}
	if len(IDs) == 0 {
		return []Monitor{}, errNoWaylandNoX
	}
	// Build and return a []Monitor slice
	monitors := []Monitor{}
	for i, ID := range IDs {
		monitors = append(monitors, Monitor{ID, widths[i], heights[i], wDPIs[i], hDPIs[i]})
	}
	return monitors, nil
}

// AverageResolution returns the average resolution for all connected monitors.
func AverageResolution() (*Res, error) {
	monitors, err := Monitors()
	if err != nil {
		return nil, err
	}
	var ws, hs uint
	for _, mon := range monitors {
		ws += mon.Width
		hs += mon.Height
	}
	ws /= uint(len(monitors))
	hs /= uint(len(monitors))
	return NewRes(ws, hs), nil
}
