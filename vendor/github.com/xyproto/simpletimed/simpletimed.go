package simpletimed

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

type Wallpaper struct {
	STWVersion  string
	Name        string
	Format      string
	Path        string // not part of the file data, but handy when parsing
	Statics     []*Static
	Transitions []*Transition
	LoopWait    time.Duration // how long the main event loop should sleep
}

type Static struct {
	At       time.Time
	Filename string
}

type Transition struct {
	From         time.Time
	UpTo         time.Time
	FromFilename string
	ToFilename   string
	Type         string
}

var DefaultLoopTime = 5 * time.Second

func (t *Transition) Duration() time.Duration {
	return mod24(t.UpTo.Sub(t.From))
}

func (t *Transition) String(format string) string {
	if !strings.Contains(format, "%s") {
		// Return the verbose version, where type is always included and the filename is not reduced with a common string format
		if t.Type == "overlay" {
			return fmt.Sprintf("@%s-%s: %s .. %s", cFmt(t.From), cFmt(t.UpTo), t.FromFilename, t.ToFilename)
		}
		return fmt.Sprintf("@%s-%s: %s .. %s | %s", cFmt(t.From), cFmt(t.UpTo), t.FromFilename, t.ToFilename, t.Type)
	}
	fields := strings.SplitN(format, "%s", 2)
	prefix := fields[0]
	suffix := fields[1]
	if t.Type == "overlay" {
		return fmt.Sprintf("@%s-%s: %s .. %s", cFmt(t.From), cFmt(t.UpTo), t.FromFilename[len(prefix):len(t.FromFilename)-len(suffix)], t.ToFilename[len(prefix):len(t.ToFilename)-len(suffix)])
	}
	return fmt.Sprintf("@%s-%s: %s .. %s | %s", cFmt(t.From), cFmt(t.UpTo), t.FromFilename[len(prefix):len(t.FromFilename)-len(suffix)], t.ToFilename[len(prefix):len(t.ToFilename)-len(suffix)], t.Type)
}

func (s *Static) String(format string) string {
	if !strings.Contains(format, "%s") {
		// Return the verbose version, where type is always included and the filename is not reduced with a common string format
		return fmt.Sprintf("@%s: %s", cFmt(s.At), s.Filename)
	}
	fields := strings.SplitN(format, "%s", 2)
	prefix := fields[0]
	suffix := fields[1]
	return fmt.Sprintf("@%s: %s", cFmt(s.At), s.Filename[len(prefix):len(s.Filename)-len(suffix)])
}

// String outputs a valid STW file, where the timestamps are in a sorted order
func (stw *Wallpaper) String() string {
	var lines []string
	for _, s := range stw.Statics {
		lines = append(lines, s.String(stw.Format))
	}
	for _, t := range stw.Transitions {
		lines = append(lines, t.String(stw.Format))
	}
	sort.Strings(lines)
	return fmt.Sprintf("stw: %s\nname: %s\nformat: %s\n", stw.STWVersion, stw.Name, stw.Format) + strings.Join(lines, "\n")
}

func NewWallpaper(version, name, format string) *Wallpaper {
	var (
		statics     []*Static
		transitions []*Transition
	)
	return &Wallpaper{version, name, format, "", statics, transitions, DefaultLoopTime}
}

func (stw *Wallpaper) AddStatic(at time.Time, filename string) {
	var s Static
	s.At = at
	if len(stw.Format) > 0 {
		s.Filename = fmt.Sprintf(stw.Format, filename)
	} else {
		s.Filename = filename
	}
	stw.Statics = append(stw.Statics, &s)
}

func (stw *Wallpaper) AddTransition(from, upto time.Time, fromFilename, toFilename, transitionType string) {
	var t Transition
	t.From = from
	t.UpTo = upto
	if len(stw.Format) > 0 {
		t.FromFilename = fmt.Sprintf(stw.Format, fromFilename)
		t.ToFilename = fmt.Sprintf(stw.Format, toFilename)
	} else {
		t.FromFilename = fromFilename
		t.ToFilename = toFilename
	}
	if len(transitionType) == 0 {
		t.Type = "overlay"
	} else {
		t.Type = transitionType
	}
	stw.Transitions = append(stw.Transitions, &t)
}

func ParseSTW(filename string) (*Wallpaper, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return DataToSimple(filename, data)
}

// DataToSimple converts from the contents of a Simple Timed Wallpaper file to
// a Wallpaper structs. The given path is used in the error messages
// and for setting stw.Path.
func DataToSimple(path string, data []byte) (*Wallpaper, error) {
	var ts []*Transition
	var ss []*Static
	parsed := make(map[string]string)
	for lineCount, byteLine := range bytes.Split(data, []byte("\n")) {
		trimmed := strings.TrimSpace(string(byteLine))
		if strings.HasPrefix(trimmed, "//") {
			fmt.Fprintf(os.Stderr, trimmed[1:])
			continue
		} else if strings.HasPrefix(trimmed, "//") {
			fmt.Fprintf(os.Stderr, trimmed[2:])
			continue
		} else if len(trimmed) == 0 {
			continue
		}
		if strings.HasPrefix(trimmed, "@") {
			if len(trimmed) > 6 && (trimmed[6] == ' ' || trimmed[6] == '-') && (trimmed[7] != ':') {
				if strings.Count(trimmed, "-") < 1 {
					return nil, fmt.Errorf("could not parse %s (no dash), line %d: %s", path, lineCount, trimmed)
				}
				fields := strings.SplitN(trimmed[1:], "-", 2)
				time1 := strings.TrimSpace(fields[0])
				if strings.Count(fields[1], ":") < 2 {
					return nil, fmt.Errorf("could not parse %s (missing colon), line %d: %s", path, lineCount, trimmed)
				}
				fields = strings.SplitN(fields[1], ":", 3)
				time2 := strings.TrimSpace(fields[0] + ":" + fields[1])
				filenames := fields[2]
				if !strings.Contains(filenames, "..") {
					return nil, fmt.Errorf("could not parse %s (missing \"..\"), line %d: %s", path, lineCount, trimmed)
				}
				fields = strings.SplitN(filenames, "..", 2)
				filename1 := strings.TrimSpace(fields[0])
				filename2 := strings.TrimSpace(fields[1])
				transitionType := "overlay"
				if strings.Contains(filename2, "|") {
					fields := strings.SplitN(filename2, "|", 2)
					filename2 = strings.TrimSpace(fields[0])
					transitionType = strings.TrimSpace(fields[1])
				}
				//fmt.Println("TRANSITION", time1, "|", time2, "|", filename1, "|", filename2, "|", transitionType)
				t1, err := time.Parse("15:04", time1)
				if err != nil {
					return nil, fmt.Errorf("could not parse %s (time), line %d: %s", path, lineCount, trimmed)
				}
				t2, err := time.Parse("15:04", time2)
				if err != nil {
					return nil, fmt.Errorf("could not parse %s (time), line %d: %s", path, lineCount, trimmed)
				}
				ts = append(ts, &Transition{t1, t2, filename1, filename2, transitionType})
			} else {
				if strings.Count(trimmed, ":") < 2 {
					return nil, fmt.Errorf("could not parse %s (missing colon), line %d: %s", path, lineCount, trimmed)
				}
				fields := strings.SplitN(trimmed[1:], ":", 3)
				time1 := strings.TrimSpace(fields[0] + ":" + fields[1])
				filename := strings.TrimSpace(fields[2])
				//fmt.Println("STATIC", time1, "|", filename)
				t1, err := time.Parse("15:04", time1)
				if err != nil {
					return nil, fmt.Errorf("could not parse %s (time), line %d: %s", path, lineCount, trimmed)
				}
				ss = append(ss, &Static{t1, filename})
			}
		} else if strings.Contains(trimmed, ":") {
			//fmt.Println("FIELD", trimmed)
			if strings.Count(trimmed, ":") < 1 {
				return nil, fmt.Errorf("could not parse %s (missing colon), line %d: %s", path, lineCount, trimmed)
			}
			fields := strings.SplitN(trimmed, ":", 2)
			key := strings.TrimSpace(fields[0])
			value := strings.TrimSpace(fields[1])
			parsed[key] = value
		} else {
			return nil, fmt.Errorf("could not parse %s (invalid syntax), line %d: %s", path, lineCount, trimmed)
		}
	}
	version, ok := parsed["stw"]
	if !ok {
		return nil, fmt.Errorf("could not find stw field in %s", path)
	}
	name, _ := parsed["name"]     // optional
	format, _ := parsed["format"] // optional
	//pacman, _ := parsed["ILoveCandy"] // optional

	stw := NewWallpaper(version, name, format)
	stw.Path = path
	for _, t := range ts {
		// Adding transitions in a way that make sure the format string is used when interpreting the filenames
		stw.AddTransition(t.From, t.UpTo, t.FromFilename, t.ToFilename, t.Type)
	}
	for _, s := range ss {
		// Adding static images in a way that make sure the format string is used when interpreting the filenames
		stw.AddStatic(s.At, s.Filename)
	}
	//fmt.Println(stw)
	return stw, nil
}
