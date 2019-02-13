package monitor

import (
	"github.com/stretchr/powerwalk"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"
)

const (
	// Minimum dimensions for qualifying as a "wallpaper"
	minimumWidth  = 640
	minimumHeight = 480
)

var (
	DefaultWallpaperDirectories = []string{"/usr/share/pixmaps", "/usr/share/wallpapers", "/usr/share/backgrounds", "/usr/local/share/pixmaps", "/usr/local/share/wallpapers", "/usr/local/share/backgrounds", "/usr/share/archlinux"}
)

type SearchResults struct {
	wallpapers                  sync.Map                // stores the full path -> *Wallpaper struct, for png + jpeg files
	gnomeWallpapers             sync.Map                // stores the full path -> *GnomeTimedWallpaper struct, for xml files
	simpleTimedWallpapers       sync.Map                // stores the full path -> *SimpleTimedWallpaper struct, for stw files
	sortedWallpapers            []*Wallpaper            // holds sorted wallpapers
	sortedGnomeTimedWallpapers  []*GnomeTimedWallpaper  // holds sorted Gnome Timed Wallpapers
	sortedSimpleTimedWallpapers []*SimpleTimedWallpaper // holds sorted Simple Timed Wallpapers
}

var (
	numCPU = runtime.NumCPU()

	defaultLoopTime = 5 * time.Second

	// Global search results struct, to be filled by the visit function that is passed to the powerwalk.WalkLimit function
	sr *SearchResults
)

// Reset the search, prepare to search again
func NewSearchResults() *SearchResults {
	return &SearchResults{
		wallpapers:                  sync.Map{},
		gnomeWallpapers:             sync.Map{},
		simpleTimedWallpapers:       sync.Map{},
		sortedWallpapers:            []*Wallpaper{},
		sortedGnomeTimedWallpapers:  []*GnomeTimedWallpaper{},
		sortedSimpleTimedWallpapers: []*SimpleTimedWallpaper{},
	}
}

// collectionName will strip away the last part of the path, until the remaining last word is no "pixmaps", "contents", "images", "backgrounds", or "wallpapers".
// This is usually the name of the wallpaper collection.
func collectionName(path string) string {
	dir := filepath.Dir(path)
	for {
		switch filepath.Base(dir) {
		case "pixmaps", "contents", "images", "wallpapers", "backgrounds":
			dir = filepath.Dir(dir)
		default:
			return filepath.Base(dir)
		}
	}
}

// partOfCollection checks if it is likely that a given filename is part of a wallpaper collection
func partOfCollection(filename string) bool {
	// filename contains width x height and is preceded by either a "_" or nothing
	_, err := FilenameToRes(filename)
	return err == nil
}

// pngSize returns the with and height of a PNG file, without reading the entire file
func pngSize(path string) (uint, uint, error) {
	pngFile, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	ic, err := png.DecodeConfig(pngFile)
	if err != nil {
		return 0, 0, err
	}
	return uint(ic.Width), uint(ic.Height), nil
}

// jpegSize returns the with and height of a JPEG file, without reading the entire file
func jpegSize(path string) (uint, uint, error) {
	jpegFile, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	ic, err := jpeg.DecodeConfig(jpegFile)
	if err != nil {
		return 0, 0, err
	}
	return uint(ic.Width), uint(ic.Height), nil
}

// largeEnough checks if the given size is equal to or larger than the global minimum size
func largeEnough(width, height uint) bool {
	return (width >= minimumWidth) && (height >= minimumHeight)
}

// visit is called per file that is found, and will be called concurrently by powerwalk.WalkLimit
func (sr *SearchResults) visit(path string, f os.FileInfo, err error) error {
	switch filepath.Ext(path) {
	case ".png":
		width, height, err := pngSize(path)
		if err != nil {
			return err
		}
		if !largeEnough(width, height) {
			return nil
		}
		wp := &Wallpaper{collectionName(path), path, width, height, partOfCollection(path)}
		sr.wallpapers.Store(path, wp)
	case ".jpg", ".jpeg":
		width, height, err := jpegSize(path)
		if err != nil {
			return err
		}
		if !largeEnough(width, height) {
			return nil
		}
		wp := &Wallpaper{collectionName(path), path, width, height, partOfCollection(path)}
		sr.wallpapers.Store(path, wp)
	case ".svg":
		// TODO: Consider supporting SVG wallpapers in the future
		return nil
	case ".xpm", ".xbm":
		// TODO: Consider supporting XPM and/or XBM wallpapers in the future
		return nil
	case ".stw": // Simple Timed Wallpaper
		stw, err := ParseSTW(path)
		if err != nil {
			return err
		}
		sr.simpleTimedWallpapers.Store(path, stw)
	case ".xml":
		gw, err := ParseXML(path)
		if err != nil {
			return err
		}
		sr.gnomeWallpapers.Store(path, gw)
	}
	return nil
}

// sortWallpapers sorts the found wallpapers
func (sr *SearchResults) sortWallpapers() {
	var collected []*Wallpaper
	sr.wallpapers.Range(func(_, value interface{}) bool {
		wp, ok := value.(*Wallpaper)
		if !ok {
			// internal error
			panic("a value in the wallpapers map is not a pointer to a Wallpaper struct")
		}
		collected = append(collected, wp)
		return true
	})
	// Now sort the collected wallpapers by the collection name, and then by the size
	sort.Slice(collected, func(i, j int) bool {
		if collected[i].CollectionName == collected[j].CollectionName {
			return (collected[i].Width * collected[i].Height) < (collected[j].Width * collected[i].Height)
		}
		return collected[i].CollectionName < collected[j].CollectionName
	})
	sr.sortedWallpapers = collected
}

// sortGnomeTimedWallpapers sorts the Found gnome Timed Wallpapers
func (sr *SearchResults) sortGnomeTimedWallpapers() {
	var collected []*GnomeTimedWallpaper
	sr.gnomeWallpapers.Range(func(_, value interface{}) bool {
		gw, ok := value.(*GnomeTimedWallpaper)
		if !ok {
			// internal error
			panic("a value in the gnomeWallpapers map is not a pointer to a GnomeTimedWallpaper struct")
		}
		collected = append(collected, gw)
		return true
	})
	// Now sort the collected GNOME wallpapers by the collection name
	sort.Slice(collected, func(i, j int) bool {
		return collected[i].Name < collected[j].Name
	})
	sr.sortedGnomeTimedWallpapers = collected
}

// sortSimpleTimedWallpapers sorts the found Simple Timed Wallpapers
func (sr *SearchResults) sortSimpleTimedWallpapers() {
	var collected []*SimpleTimedWallpaper
	sr.simpleTimedWallpapers.Range(func(_, value interface{}) bool {
		stw, ok := value.(*SimpleTimedWallpaper)
		if !ok {
			// internal error
			panic("a value in the simpleTimedWallpapers map is not a pointer to a SimpleTimedWallpaper struct")
		}
		collected = append(collected, stw)
		return true
	})
	// Now sort the collected Simple Timed Wallpapers by the collection name
	sort.Slice(collected, func(i, j int) bool {
		return collected[i].Name < collected[j].Name
	})
	sr.sortedSimpleTimedWallpapers = collected
}

// FindWallpapers will collect and parse wallpapers and GNOME background XML files in all default wallpaper directories
func FindWallpapers() (*SearchResults, error) {
	sr := NewSearchResults()
	for _, path := range DefaultWallpaperDirectories {
		// Search the given path, using the sr.visit function
		if err := powerwalk.WalkLimit(path, sr.visit, numCPU); err != nil {
			return nil, err
		}
	}
	sr.sortWallpapers()
	sr.sortSimpleTimedWallpapers()
	sr.sortGnomeTimedWallpapers()
	return sr, nil
}

// FindCollectionNames gathers all the names of all available wallpaper packs or GNOME timed backgrounds
func (sr *SearchResults) CollectionNames() []string {
	var collectionNames []string
	for _, wp := range sr.sortedWallpapers {
		if wp.PartOfCollection {
			collectionNames = append(collectionNames, wp.CollectionName)
		}
	}
	for _, gw := range sr.sortedGnomeTimedWallpapers {
		collectionNames = append(collectionNames, gw.Name)
	}
	for _, stw := range sr.sortedSimpleTimedWallpapers {
		collectionNames = append(collectionNames, stw.Name)
	}
	return unique(collectionNames)
}

// Wallpapers returns a sorted slice of all found wallpapers
func (sr *SearchResults) Wallpapers() []*Wallpaper {
	return sr.sortedWallpapers
}

// GnomeTimedWallpapers returns a sorted slice of all found gnome timed wallpapers
func (sr *SearchResults) GnomeTimedWallpapers() []*GnomeTimedWallpaper {
	return sr.sortedGnomeTimedWallpapers
}

// SimpleTimedWallpapers returns a sorted slice of all found simple timed wallpapers
func (sr *SearchResults) SimpleTimedWallpapers() []*SimpleTimedWallpaper {
	return sr.sortedSimpleTimedWallpapers
}

// WallpapersByName will return simple timed wallpapers that match with the collection name
func (sr *SearchResults) WallpapersByName(name string) []*Wallpaper {
	var collection []*Wallpaper
	for _, wp := range sr.sortedWallpapers {
		if wp.PartOfCollection && wp.CollectionName == name {
			collection = append(collection, wp)
		}
	}
	return collection
}

// GnomeTimedWallpapersByName will return gnome timed wallpapers that match with the collection name
func (sr *SearchResults) GnomeTimedWallpapersByName(name string) []*GnomeTimedWallpaper {
	var collection []*GnomeTimedWallpaper
	for _, gw := range sr.sortedGnomeTimedWallpapers {
		if gw.Name == name {
			collection = append(collection, gw)
		}
	}
	return collection
}

// SimpleTimedWallpapersByName will return simple timed wallpapers that match with the collection name
func (sr *SearchResults) SimpleTimedWallpapersByName(name string) []*SimpleTimedWallpaper {
	var collection []*SimpleTimedWallpaper
	for _, stw := range sr.sortedSimpleTimedWallpapers {
		if stw.Name == name {
			collection = append(collection, stw)
		}
	}
	return collection
}

func (sr *SearchResults) Empty() bool {
	return len(sr.sortedSimpleTimedWallpapers) == 0 && len(sr.sortedGnomeTimedWallpapers) == 0 && len(sr.sortedWallpapers) == 0
}

func (sr *SearchResults) NoTimedWallpapers() bool {
	return len(sr.sortedSimpleTimedWallpapers) == 0 && len(sr.sortedGnomeTimedWallpapers) == 0
}
