package monitor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"time"
)

// Handle the GNOME timed wallpaper XML format

type GBackground struct {
	XMLName         xml.Name      `xml:"background"`
	StartTime       GStartTime    `xml:"starttime"`
	Statics         []GStatic     `xml:"static"`
	Transitions     []GTransition `xml:"transition"`
	staticOrder     StaticMap
	transitionOrder TransitionMap
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

type TransitionMap map[int]int
type StaticMap map[int]int

// Duration returns how long a static wallpaper should last
func (s *GStatic) Duration() time.Duration {
	return time.Duration(s.Seconds) * time.Second
}

// Duration returns how long a transition should last
func (s *GTransition) Duration() time.Duration {
	return time.Duration(s.Seconds) * time.Second
}

func Parse(XMLFilename string) (*GBackground, error) {
	data, err := ioutil.ReadFile(XMLFilename)
	if err != nil {
		return nil, err
	}
	var background GBackground
	xml.Unmarshal(data, &background)

	// After parsing the XML, find the order of the <static> and <transition>
	// tags. This is needed later, when calculating the event times.
	background.staticOrder, background.transitionOrder, err = findOrder(data)
	if err != nil {
		return nil, err
	}

	//log.Println("staticOrder", background.staticOrder)
	//log.Println("transitionOrder", background.transitionOrder)

	return &background, nil
}

// Given a transition, find the index in the total collection of elements with
// a duration.
func (gb *GBackground) TransitionOrder(t *GTransition) (int, error) {
	for i, tElement := range gb.Transitions {
		if t == &tElement {
			return gb.transitionOrder[i], nil
		}
	}
	return -1, errors.New("Could not find the given GTransition in the collection")
}

// Given a static, find the index in the total collection of elements with
// a duration.
func (gb *GBackground) StaticOrder(s *GStatic) (int, error) {
	for i, sElement := range gb.Statics {
		if s == &sElement {
			return gb.staticOrder[i], nil
		}
	}
	return -1, errors.New("Could not find the given GStatic in the collection")
}

// The order in the XML matters when calculating the timing.
// This function returns two maps. One that maps transition
// index to overall transition+static index and one that maps
// static index to overall transition+static index.
func findOrder(XMLData []byte) (StaticMap, TransitionMap, error) {
	staticTag := []byte("<static")
	transitionTag := []byte("<transition")

	staticCount := bytes.Count(XMLData, staticTag)
	transitionCount := bytes.Count(XMLData, transitionTag)

	if staticCount == 0 && transitionCount == 0 {
		// No static and no transition tags, OK
		return nil, nil, nil
	}

	// Start off by searching from the very start of the data
	offset := 0

	// Keep track of encountered "<static" strings
	staticCounter := 0

	// Keep track of encountered <transition" strings
	transitionCounter := 0

	//staticIndex -> totalIndex
	staticOrder := make(StaticMap, staticCount)

	//transitionIndex -> totalIndex
	transitionOrder := make(TransitionMap, transitionCount)

	// TODO: Strip away all comments before processing the XML data

	// Iterate one time per static or transition tag
	for count := 0; count < (staticCount + transitionCount); count++ {
		sPos := bytes.Index(XMLData[offset:], staticTag)
		tPos := bytes.Index(XMLData[offset:], transitionTag)
		// Use the smallest found index
		if sPos < tPos && sPos != -1 {
			// Found static tag
			pos := offset + sPos
			//log.Println("STATIC", pos, staticCounter, "->", count)
			// Record the static index and the total transition/static index
			staticOrder[staticCounter] = count
			// Increase the static tag counter
			staticCounter++
			// Increase the offset with the found position
			offset = pos + len(staticTag)
		} else if tPos != -1 {
			pos := offset + tPos
			//log.Println("TRANSITION", pos, transitionCounter, "->", count)
			// Record the transition index and the total transition/static index
			transitionOrder[transitionCounter] = count
			// Increase the transition tag counter
			transitionCounter++
			// Increase the offset with the found position
			offset = pos + len(transitionTag)
		} else {
			// No more matches
			break
		}
	}
	return staticOrder, transitionOrder, nil
}

func (background *GBackground) String() string {
	data, err := xml.MarshalIndent(background, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}
