package wallutils

// Functions and structs for detecting overlap between monitors / display configurations

import (
	"errors"
	"fmt"
	"image"
	"math"
)

func NewRect(x, y, w, h uint) image.Rectangle {
	return image.Rect(int(x), int(y), int(x+w), int(y+h))
}

// gdc applies the Euclid algorithm for greatest common divisor on
// a slice of numbers.
func gdc(nums []int) (int, error) {
	if len(nums) < 2 {
		return 0, errors.New("need more than one number for gdc")
	}

	// This is the result for only 2 integers
	x := nums[0]
	y := nums[1]
	for y != 0 {
		x, y = y, x%y
	}
	result := x

	// For loop in case there're more than 2 ints
	for j := 2; j < len(nums); j++ {
		x = result
		y = nums[j]
		for y != 0 {
			x, y = y, x%y
		}
		result = x
	}

	return result, nil
}

// overlaps checks if any rectangles in a slice of rectangles overlaps.
func overlaps(rects []image.Rectangle) bool {
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

// overlaps checks if any rectangles in a slice of rectangles overlaps.
func overlapsOld(rects []image.Rectangle) bool {
	// Shrink all rectangles down to minimum size by dividing on the
	// common greatest denominator, then draw "pixels" in a grid and check if
	// a "pixel" is drawn twice. Can probably be done currently too.

	var coordinates []int
	for _, r := range rects {
		coordinates = append(coordinates, r.Min.X, r.Min.Y, r.Max.X-r.Min.X, r.Max.Y-r.Min.Y)
	}

	d, err := gdc(coordinates)
	if err != nil {
		// Too few rectangles for any overlap
		return false
	}

	var (
		minx int = math.MaxInt16
		maxx int
		miny int = math.MaxInt16
		maxy int
	)

	// Scale down all rectangles, and find the min/max values
	for _, r := range rects {
		r.Min.X /= d
		r.Min.Y /= d
		r.Max.X /= d
		r.Max.Y /= d
		if r.Min.X < minx {
			minx = r.Min.X
		}
		if r.Max.X > maxx {
			maxx = r.Max.X
		}
		if r.Min.Y < miny {
			miny = r.Min.Y
		}
		if r.Max.Y > maxy {
			maxy = r.Max.Y
		}
	}

	// Find the width and height of the boundaries
	width := maxx - minx
	height := maxy - miny

	fmt.Println("minx, maxx, miny, maxy", minx, maxx, miny, maxy)
	fmt.Println("width, height", width, height)
	fmt.Printf("Using rectangle overlap buffer of size %d\n", width*height)

	// For the case of monitor resolutions, this slice of "pixels"
	// should be significantly smaller than the original sizes.
	pixels := make([]int, width*height)

	// now loop over the scaled down "pixels" and mark them
	// if one the "pixels" are above 1, there is overlap
	for y := miny; y <= maxy; y++ {
		for x := minx; x <= maxx; x++ {
			for _, r := range rects {
				if r.Min.Y <= y && y < r.Max.Y {
					if r.Min.X <= x && x < r.Max.X {
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
