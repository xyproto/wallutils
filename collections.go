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
)

const (
	// Minimum dimensions for qualifying as a "wallpaper"
	minimumWidth  = 640
	minimumHeight = 480
)

// TODO: Move these global variables into structs

var (
	numCPU = runtime.NumCPU()

	images          sync.Map // stores the full path -> *Wallpaper struct, for png + jpeg files
	gnomeWallpapers sync.Map // stores the full path -> *GnomeWallpaper struct, for xml files

	searchComplete bool // will only search once, until the search is reset
)

// Reset the search, prepare to search again
func ResetSearch() {
	images = sync.Map{}
	gnomeWallpapers = sync.Map{}
	searchComplete = false
}

func collectionName(path string) string {
	dir := filepath.Dir(path)
	// Strip away the latest directory of the path until it is not a generic
	// folder name, but may be the name of the wallpaper collection.
	for {
		switch filepath.Base(dir) {
		case "pixmaps", "contents", "images", "wallpapers":
			dir = filepath.Dir(dir)
		default:
			return filepath.Base(dir)
		}
	}
}

// partOfCollection checks if it is likely that a given filename is part of a wallpaper collection
func partOfCollection(filename string) bool {
	// filename contains width x height and is preceeded by either a "_" or nothing
	_, err := FilenameToRes(filename)
	return err == nil
}

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

func largeEnough(width, height uint) bool {
	return (width >= minimumWidth) && (height >= minimumHeight)
}

// visit is called per file that is found, and will be called concurrently
func visit(path string, f os.FileInfo, err error) error {
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
		images.Store(path, wp)
	case ".jpg", ".jpeg":
		width, height, err := jpegSize(path)
		if err != nil {
			return err
		}
		if !largeEnough(width, height) {
			return nil
		}
		wp := &Wallpaper{collectionName(path), path, width, height, partOfCollection(path)}
		images.Store(path, wp)
	case ".svg":
		// TODO: Consider supporting SVG wallpapers in the future
		//fmt.Println("SVG ", path)
		return nil
	case ".xpm", ".xbm":
		// TODO: Consider supporting XPM and/or XBM wallpapers in the future
		//fmt.Println("X bitmap", path)
		return nil
	case ".xml":
		gb, err := Parse(path)
		if err != nil {
			return err
		}
		// Use the name of the XML, before the filename extension, as the collection name
		name := firstname(filepath.Base(path))
		gw := &GnomeWallpaper{name, path, gb}
		gnomeWallpapers.Store(path, gw)
	}
	//fmt.Println("@" + path)
	return nil
}

func searchPath(path string) {
	err := powerwalk.WalkLimit(path, visit, numCPU)
	if err != nil {
		panic(err)
	}
	searchComplete = true
}

func foundWallpapers() []*Wallpaper {
	var collected []*Wallpaper
	images.Range(func(_, value interface{}) bool {
		wp, ok := value.(*Wallpaper)
		if !ok {
			// internal error
			panic("a value in the images map is not a pointer to a Wallpaper struct")
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
	return collected
}

func foundGnomeWallpapers() []*GnomeWallpaper {
	var collected []*GnomeWallpaper
	gnomeWallpapers.Range(func(_, value interface{}) bool {
		gw, ok := value.(*GnomeWallpaper)
		if !ok {
			// internal error
			panic("a value in the gnomeWallpapers map is not a pointer to a GnomeWallpaper struct")
		}
		collected = append(collected, gw)
		return true
	})
	// Now sort the collected GNOME wallpapers by the collection name
	sort.Slice(collected, func(i, j int) bool {
		return collected[i].CollectionName < collected[j].CollectionName
	})
	return collected
}

// SearchPaths will concurrently collect all wallpapers that are large enough.
// Also parse Gnome Background XML files.
func SearchPaths(paths []string) ([]*Wallpaper, []*GnomeWallpaper) {
	for _, path := range paths {
		searchPath(path)
	}
	return foundWallpapers(), foundGnomeWallpapers()
}

// FindWallpapers will collect and parse wallpapers and GNOME background XML files in all default wallpaper directories
func FindWallpapers() ([]*Wallpaper, []*GnomeWallpaper) {
	return SearchPaths(DefaultWallpaperDirectories)
}

// FindCollectionNames gathers all the names of all available wallpaper packs or GNOME timed backgrounds
func FindCollectionNames() []string {
	wallpapers, gnomeWallpapers := FindWallpapers()
	var collectionNames []string
	for _, wp := range wallpapers {
		if wp.PartOfCollection {
			collectionNames = append(collectionNames, wp.CollectionName)
		}
	}
	for _, gw := range gnomeWallpapers {
		collectionNames = append(collectionNames, gw.CollectionName)
	}
	return unique(collectionNames)
}
