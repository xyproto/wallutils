package wallutils

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

var errNoWaylandNoX = errors.New("could not detect neither Wayland nor X")

// String returns a string with monitor ID and resolution
func (m Monitor) String() string {
	return fmt.Sprintf("[%d] %dx%d", m.ID, m.Width, m.Height)
}
