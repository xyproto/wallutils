// monitor is a package for dealing with monitors, resolution, dpi, wallpapers, wallpaper collections, timed wallpapers and converting to the Simple Timed Wallpaper format.
package monitor

import (
	"errors"
	"fmt"
)

// Monitor contains an ID, the width in pixels and the height in pixels
type Monitor struct {
	ID     uint // monitor number, from 0 and up
	Width  uint // width, in pixels
	Height uint // height, in pixels
	DPIw   uint // DPI, if available (width)
	DPIh   uint // DPI, if available (height)
}

var errNoWaylandNoX = errors.New("could not detect either Wayland or X")

const VersionString = "4.1.1"

// String returns a string with monitor ID and resolution
func (m Monitor) String() string {
	return fmt.Sprintf("[%d] %dx%d", m.ID, m.Width, m.Height)
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

// Detect returns information about all monitors, regardless of if it's under
// Wayland or X11. Will use additional plugins, if available.
func Detect() ([]Monitor, error) {
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
	var monitors = []Monitor{}
	for i, ID := range IDs {
		monitors = append(monitors, Monitor{ID, widths[i], heights[i], wDPIs[i], hDPIs[i]})
	}
	return monitors, nil
}
