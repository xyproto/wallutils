package monitor

import (
	"strconv"
	"strings"
	"time"
)

type GnomeWallpaper struct {

	// The name of the directory containing this XML file, if it's not
	// "pixmaps", "images" or "contents". May use the parent of the parent.
	CollectionName string

	// Path is the full path to the XML file
	Path string

	// Config contains the parsed XML. See: gnomexml.go
	Config *GBackground
}

func (gw *GnomeWallpaper) Time() time.Time {
	st := gw.Config.StartTime
	return time.Date(st.Year, time.Month(st.Month), st.Day, st.Hour, st.Minute, 0, 0, time.Local)
}

func (gw *GnomeWallpaper) TodayTime() time.Time {
	// Get hour, minute and second from the timed wallpaper
	st := gw.Config.StartTime
	hour, min, sec := st.Clock()

	// Get the rest of the fields from the current time
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, now.Nanosecond(), now.Location())
}

func (gw *GnomeWallpaper) Images() []string {
	var filenames []string
	for _, static := range gw.Config.Statics {
		filenames = append(filenames, static.Filename)
	}
	for _, transition := range gw.Config.Transitions {
		filenames = append(filenames, transition.FromFilename)
		filenames = append(filenames, transition.ToFilename)
	}
	return unique(filenames)
}

// String builds a string with various information about this GNOME timed wallpaper
func (gw *GnomeWallpaper) String() string {
	var sb strings.Builder
	sb.WriteString("--- ")
	sb.WriteString(gw.CollectionName)
	sb.WriteString("---\npath\t\t\t= ")
	sb.WriteString(gw.Path)
	sb.WriteString("\nstart time\t\t= ")
	sb.WriteString(gw.Time().String())
	sb.WriteString("\nnumber of static tags\t= ")
	sb.WriteString(strconv.Itoa(len(gw.Config.Statics)))
	sb.WriteString("\nnumber of transitions\t= ")
	sb.WriteString(strconv.Itoa(len(gw.Config.Transitions)))
	sb.WriteString("\nuses these images:\n")
	for _, filename := range gw.Images() {
		sb.WriteString("\t" + filename + "\n")
	}
	return sb.String()
}
