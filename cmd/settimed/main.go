package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xyproto/gnometimed"
	"github.com/xyproto/simpletimed"
	"github.com/xyproto/wallutils"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func main() {
	verbose := true

	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Please give a timed wallpaper name as the first argument.")
		os.Exit(1)
	}

	collectionName := os.Args[1]

	// Ok, it was a filename
	if strings.Contains(collectionName, ".") && exists(collectionName) {
		filename := collectionName
		switch filepath.Ext(filename) {
		case ".stw":
			stw, err := simpletimed.ParseSTW(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if verbose {
				fmt.Printf("Using: %s\n", stw.Path)
			}
			// Start endless event loop
			if err := stw.EventLoop(verbose, func(path string) error { return wallutils.SetWallpaperVerbose(path, verbose) }); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		case ".xml":
			gtw, err := gnometimed.ParseXML(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if verbose {
				fmt.Printf("Using: %s\n", gtw.Path)
			}
			// Start endless event loop
			if err := gtw.EventLoop(verbose, func(path string) error { return wallutils.SetWallpaperVerbose(path, verbose) }); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		default:
			fmt.Fprintln(os.Stderr, "Unrecognized file extension:", filepath.Ext(filename))
			os.Exit(1)
		}
	}

	if verbose {
		fmt.Printf("Setting wallpaper collection: %s\n", collectionName)
		fmt.Println("Searching for wallpapers...")
	}
	searchResults, err := wallutils.FindWallpapers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if searchResults.NoTimedWallpapers() {
		fmt.Fprintln(os.Stderr, "Could not find any timed wallpapers on the system.")
		os.Exit(1)
	}
	if verbose {
		fmt.Println("Filtering wallpapers by name...")
	}
	simpleTimedWallpapers := searchResults.SimpleTimedWallpapersByName(collectionName)
	gnomeTimedWallpapers := searchResults.GnomeTimedWallpapersByName(collectionName)

	// gnomeTimedWallpapers and simpleTimedWallpapers have now been filtered so that they only contain elements with matching collection names

	if (len(gnomeTimedWallpapers) == 0) && (len(simpleTimedWallpapers) == 0) {
		fmt.Fprintln(os.Stderr, "No such timed wallpaper: "+collectionName)
		os.Exit(1)
	}

	if (len(gnomeTimedWallpapers) > 1) || (len(simpleTimedWallpapers) > 1) {
		fmt.Fprintln(os.Stderr, "Found several timed backgrounds, with the same name.")
		os.Exit(1)
	}

	if len(simpleTimedWallpapers) == 1 {
		stw := simpleTimedWallpapers[0]
		if verbose {
			fmt.Printf("Using: %s\n", stw.Path)
		}
		// Start endless event loop
		if err := stw.EventLoop(verbose, func(path string) error { return wallutils.SetWallpaperVerbose(path, verbose) }); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else if len(gnomeTimedWallpapers) == 1 {
		gtw := gnomeTimedWallpapers[0]
		if verbose {
			fmt.Printf("Using: %s\n", gtw.Path)
		}
		// Start endless event loop
		if err := gtw.EventLoop(verbose, func(path string) error { return wallutils.SetWallpaperVerbose(path, verbose) }); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

}
