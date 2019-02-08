package monitor

import (
	"encoding/xml"
	"io/ioutil"
)

// Handle the GNOME wallpaper animation XML format

type GBackground struct {
	XMLName     xml.Name      `xml:"background"`
	StartTime   GStartTime    `xml:"starttime"`
	Statics     []GStatic     `xml:"static"`
	Transitions []GTransition `xml:"transition"`
}

type GStartTime struct {
	XMLName xml.Name `xml:"starttime"`
	Year    int      `xml:"year"`
	Month   int      `xml:"month"`
	Day     int      `xml:"day"`
	Hour    int      `xml:"hour"`
	Minute  int      `xml:"minute"`
	Second  int      `xml:"second"`
}

type GStatic struct {
	XMLName  xml.Name `xml:"static"`
	Seconds  float64  `xml:"duration"`
	Filename string   `xml:"file"`
}

type GTransition struct {
	XMLName      xml.Name `xml:"transition"`
	Type         string   `xml:"type,attr,omitempty"`
	Seconds      float64  `xml:"duration"`
	FromFilename string   `xml:"from"`
	ToFilename   string   `xml:"to"`
}

func Parse(XMLFilename string) (*GBackground, error) {
	data, err := ioutil.ReadFile(XMLFilename)
	if err != nil {
		return nil, err
	}
	var background GBackground
	xml.Unmarshal(data, &background)
	return &background, nil
}

func (background *GBackground) String() string {
	data, err := xml.MarshalIndent(background, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}
