package wallutils

// Functions and structs for detecting overlap between monitors / display configurations

import (
	"errors"
	"fmt"
)

type Rect struct {
	x, y, w, h uint
}

func NewRect(x, y, w, h uint) *Rect {
	return &Rect{x, y, w, h}
}

func (r *Rect) String() string {
	return fmt.Sprintf("(%d, %d, %d, %d)", r.x, r.y, r.w, r.h)
}

// gdc2 implements the Euclid algorithm for greatest common divisor
// on two given numbers.
func gcd2(x, y uint) uint {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

// gdc applies the Euclid algorithm for greatest common divisor on
// a slice of numbers.
func gdc(nums []uint) (uint, error) {
	if len(nums) < 2 {
		return 0, errors.New("need more than one number for gdc")
	}

	// This is the result for only 2 integers
	result := gcd2(nums[0], nums[1])

	// For loop in case there're more than 2 ints
	for j := 2; j < len(nums); j++ {
		result = gcd2(result, nums[j])
	}

	return result, nil
}

// overlaps checks if a slice of rectangles overlaps.
// will modify (reduce the size of) the rectangles in the process.
func overlaps(rects []*Rect) bool {
	// Shrink all rectangles down to minimum size by dividing on the
	// common greatest denominator, then draw "pixels" in a grid and check if
	// a "pixel" is drawn twice. Can probably be done currently too.

	var coordinates []uint
	for _, r := range rects {
		coordinates = append(coordinates, r.x, r.y, r.w, r.h)
	}

	d, err := gdc(coordinates)
	if err != nil {
		// Too few rectangles for any overlap
		return false
	}

	//fmt.Println("GDC", d)

	var (
		minx uint = 1<<16 - 1
		maxx uint = 0
		miny uint = 1<<16 - 1
		maxy uint = 0
	)

	// Scale down all rectangles, and find the min/max values
	for _, r := range rects {
		r.x /= d
		r.y /= d
		r.w /= d
		r.h /= d
		if r.x < minx {
			minx = r.x
		}
		if r.x+r.w > maxx {
			maxx = r.x + r.w
		}
		if r.y < miny {
			miny = r.y
		}
		if r.y+r.h > maxy {
			maxy = r.y + r.h
		}
	}

	// Scale the rectangles back up when done
	//defer func() {
	//	for _, r := range rects {
	//		r.x *= d
	//		r.y *= d
	//		r.w *= d
	//		r.h *= d
	//	}
	//}

	// Find the width and height of the boundaries
	width := maxx - minx
	height := maxy - miny

	//fmt.Println("minx, maxx, miny, maxy", minx, maxx, miny, maxy)
	//fmt.Println("width, height", width, height)
	//fmt.Printf("Using rectangle overlap buffer of size %d\n", width*height)

	// For the case of monitor resolutions, this slice of "pixels"
	// should be significantly smaller than the original sizes.
	pixels := make([]int, width*height)

	// now loop over the "pixels" and mark them
	// if one the "pixels" are above 1, there is overlap
	for y := miny; y <= maxy; y++ {
		for x := minx; x <= maxx; x++ {
			for _, r := range rects {
				if r.y <= y && y < (r.y+r.h) {
					if r.x <= x && x < (r.x+r.w) {
						index := y*width + x
						pixels[index]++
						if pixels[index] > 1 {
							return true
						}
					}
				}
			}
		}
	}

	// Found no overlap
	return false
}
