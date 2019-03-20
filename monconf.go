package wallutils

// Read ~/.config/monitors.xml, if available
// This will make it possible to detect overlapping coordinates when
// changing the desktop background.

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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

// overlaps checks if a slice of rectangles overlaps
func overlaps(rects []*Rect) bool {
	// TODO: Implement
	println("TO IMPLEMENT")
	return false
}

func (mc *MonitorConfiguration) Overlapping() bool {
	mc, err := NewMonitorConfiguration()
	if err != nil {
		return false
	}
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
