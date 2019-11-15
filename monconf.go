package wallutils

// Functions and structs for dealing with ~/.config/monitors.xml, which is used by GNOME, Cinnamon and MATE

import (
	"encoding/xml"
	"fmt"
	"image"
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

// ParseMonitorFile can parse monitor XML files,
// like the one that typically exists in ~/.config/monitors.xml
func ParseMonitorFile(filename string) (*MonitorConfiguration, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var monitors MonitorConfiguration
	if err = xml.Unmarshal(data, &monitors); err != nil {
		return nil, fmt.Errorf("could not parse %s as XML: error: %s", filename, err)
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

// Overlapping can check if configurations in monitors.xml have overlapping areas.
func (mc *MonitorConfiguration) Overlapping() bool {
	mc, err := NewMonitorConfiguration()
	if err != nil {
		return false
	}
	// Run a check per <configuration> section in the XML file
	for _, conf := range mc.Configurations {
		rects := make([]image.Rectangle, 0)
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
		if Overlaps(rects) {
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
