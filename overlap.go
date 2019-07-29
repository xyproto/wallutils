package wallutils

// Functions and structs for detecting overlap between monitors / display configurations

import (
	"image"
)

// NewRect is a convenience function for creating an image.Rectangle, given
// the upper left corner (x and y), a width and a height.
func NewRect(x, y, w, h uint) image.Rectangle {
	return image.Rect(int(x), int(y), int(x+w), int(y+h))
}

// Overlaps checks if any rectangles in a slice of rectangles overlaps.
func Overlaps(rects []image.Rectangle) bool {
	for _, ar := range rects {
		for _, br := range rects {
			if ar == br {
				continue
			}
			if ar.Overlaps(br) {
				return true
			}
		}
	}
	return false
}
