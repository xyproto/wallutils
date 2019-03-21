package wallutils

// Read ~/.config/monitors.xml, if available
// This will make it possible to detect overlapping coordinates when
// changing the desktop background.

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type MonitorConfiguration struct {
	XMLName        xml.Name         `xml:"monitors"`
	Version        string           `xml:"version,attr"`
	Configurations []MConfiguration `xml:"configuration"`
}

type MConfiguration struct {
	XMLName xml.Name  `xml:"configuration"`
	Clone   string    `xml:"clone,omitempty"`
	Outputs []MOutput `xml:"output"`
}

type MOutput struct {
	XMLName  xml.Name `xml:"output"`
	Name     string   `xml:"name,attr"`
	Vendor   string   `xml:"vendor,omitempty"`
	Product  string   `xml:"product,omitempty"`
	Serial   string   `xml:"serial,omitempty"`
	Width    string   `xml:"width,omitempty"`
	Height   string   `xml:"height,omitempty"`
	Rate     string   `xml:"rate,omitempty"`
	X        string   `xml:"x,omitempty"`
	Y        string   `xml:"y,omitempty"`
	Rotation string   `xml:"rotation,omitempty"`
	ReflectX string   `xml:"reflect_x,omitempty"`
	ReflectY string   `xml:"reflect_y,omitempty"`
	Primary  string   `xml:"primary,omitempty"`
}

type Rect struct {
	x, y, w, h uint
}

func NewRect(x, y, w, h uint) *Rect {
	return &Rect{x, y, w, h}
}

func (r *Rect) String() string {
	return fmt.Sprintf("(%d, %d, %d, %d)", r.x, r.y, r.w, r.h)
}

// ParseMonitors can parse monitor XML files,
// like the one that typically exists in ~/.config/monitors.xml
func ParseMonitorFile(filename string) (*MonitorConfiguration, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var monitors MonitorConfiguration
	if err = xml.Unmarshal(data, &monitors); err != nil {
		return nil, fmt.Errorf("Could not parse %s as XML: error: %s", filename, err)
	}

	return &monitors, nil
}

func NewMonitorConfiguration() (*MonitorConfiguration, error) {
	// Check if there are overlapping monitors (overlapping rectangles)
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	filename := filepath.Join(homedir, ".config/monitors.xml")
	return ParseMonitorFile(filename)
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

// Overlapping can check if two monitors in monitors.xml have overlapping
// areas. This is useful to know, because it may cause artifacts when setting
// the desktop wallpapers in Gnome3, Cinnamon and MATE.
func (mc *MonitorConfiguration) Overlapping() bool {
	mc, err := NewMonitorConfiguration()
	if err != nil {
		return false
	}
	// Run a check per <configuration> section in the XML file
	for _, conf := range mc.Configurations {
		rects := make([]*Rect, 0)
		for _, output := range conf.Outputs {
			if output.X != "" && output.Y != "" && output.Width != "" && output.Height != "" {
				x, err := strconv.Atoi(output.X)
				if err != nil {
					continue
				}
				y, err := strconv.Atoi(output.Y)
				if err != nil {
					continue
				}
				width, err := strconv.Atoi(output.Width)
				if err != nil {
					continue
				}
				height, err := strconv.Atoi(output.Height)
				if err != nil {
					continue
				}
				r := NewRect(uint(x), uint(y), uint(width), uint(height))
				rects = append(rects, r)
			}
		}
		if overlaps(rects) {
			return true
		}
	}
	return false
}

// MonConfOverlap is a convenience function for checking if the
// x,y,w,h monitor configurations in ie. ~/.config/monitors.xml are
// overlapping or not. If monitors.xml can not be parsed or read,
// false is returned.
func MonConfOverlap(filename string) bool {
	// Replace ~ with the home directory
	if strings.HasPrefix(filename, "~") {
		homedir, err := os.UserHomeDir()
		if err == nil {
			filename = filepath.Join(homedir, filename[1:])
		}
	}
	if mc, err := ParseMonitorFile(filename); err != nil {
		return mc.Overlapping()
	}
	return false
}
