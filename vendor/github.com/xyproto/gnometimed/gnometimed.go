package gnometimed

import (
	"strconv"
	"strings"
	"time"
)

var DefaultEventLoopDelay = 5 * time.Second

type Wallpaper struct {
	// The name of this timed wallpaper
	Name string

	// Path is the full path to the XML file
	Path string

	// Config contains the parsed XML. See: gnomexml.go
	Config *GBackground

	// LoopWait is for how long the event loop should sleep at every iteration
	LoopWait time.Duration
}

func NewWallpaper(name string, path string, config *GBackground) *Wallpaper {
	return &Wallpaper{name, path, config, DefaultEventLoopDelay}
}

// StartTime returns the timed wallpaper start time, as a time.Time
func (gtw *Wallpaper) StartTime() time.Time {
	// gtw.Config.StartTime is a struct that contains ints,
	// where the values are directly from the parsed XML.
	st := gtw.Config.StartTime
	return time.Date(st.Year, time.Month(st.Month), st.Day, st.Hour, st.Minute, 0, 0, time.Local)
}

func (gtw *Wallpaper) Images() []string {
	var filenames []string
	for _, static := range gtw.Config.Statics {
		filenames = append(filenames, static.Filename)
	}
	for _, transition := range gtw.Config.Transitions {
		filenames = append(filenames, transition.FromFilename)
		filenames = append(filenames, transition.ToFilename)
	}
	return unique(filenames)
}

// String builds a string with various information about this GNOME timed wallpaper
func (gtw *Wallpaper) String() string {
	var sb strings.Builder
	sb.WriteString("path\t\t\t= ")
	sb.WriteString(gtw.Path)
	sb.WriteString("\nstart time\t\t= ")
	sb.WriteString(gtw.StartTime().String())
	sb.WriteString("\nnumber of static tags\t= ")
	sb.WriteString(strconv.Itoa(len(gtw.Config.Statics)))
	sb.WriteString("\nnumber of transitions\t= ")
	sb.WriteString(strconv.Itoa(len(gtw.Config.Transitions)))
	sb.WriteString("\nuses these images:\n")
	for _, filename := range gtw.Images() {
		sb.WriteString("\t" + filename + "\n")
	}
	return strings.TrimSpace(sb.String())
}
