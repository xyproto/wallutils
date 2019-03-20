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
)

type MonitorConfiguration struct {
	XMLName        xml.Name `xml:"monitors"`
	Version        string   `xml:"version,attr"`
	Configurations []MConfiguration
}

type MConfiguration struct {
	XMLName xml.Name `xml:"configuration"`
	Clone   string   `xml:"clone,omitempty"`
	Outputs []MOutput
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

// ParseMonitors can parse ~/.config/monitors.xml
func ParseMonitors() (*MonitorConfiguration, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	filename := filepath.Join(homedir, ".config/monitors.xml")
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

func (mc *MonitorConfiguration) Overlapping() bool {
	// Check if there are overlapping monitors (overlapping rectangles)
	panic("TO IMPLEMENT")
	return false
}
