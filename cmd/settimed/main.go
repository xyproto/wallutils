package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
	"github.com/xyproto/gnometimed"
	"github.com/xyproto/simpletimed"
	"github.com/xyproto/wallutils"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func SetTimedWallpaper(collectionOrFilename string, verbose bool, tempImageFilename string) error {
	// Check if it is a timed wallpaper filename
	if strings.Contains(collectionOrFilename, ".") && exists(collectionOrFilename) {
		filename := collectionOrFilename
		switch filepath.Ext(filename) {
		case ".stw":
			stw, err := simpletimed.ParseSTW(filename)
			if err != nil {
				return err
			}
			if verbose {
				fmt.Printf("Launching event loop for: %s\n", stw.Path)
			}
			// Start endless event loop
			if err := stw.EventLoop(verbose,
				func(path string) error {
					return wallutils.SetWallpaperVerbose(path, verbose)
				}, tempImageFilename); err != nil {
				return err
			}
		case ".xml":
			gtw, err := gnometimed.ParseXML(filename)
			if err != nil {
				return err
			}
			if verbose {
				fmt.Printf("Launching event loop for: %s\n", gtw.Path)
			}
			// Start endless event loop
			if err := gtw.EventLoop(verbose,
				func(path string) error {
					return wallutils.SetWallpaperVerbose(path, verbose)
				}, tempImageFilename); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unrecognized file extension: %s", filepath.Ext(filename))
		}
	}

	if verbose {
		fmt.Printf("Setting timed wallpaper: %s\n", collectionOrFilename)
		fmt.Println("Searching for wallpapers...")
	}
	searchResults, err := wallutils.FindWallpapers()
	if err != nil {
		return err
	}
	if searchResults.NoTimedWallpapers() {
		return errors.New("could not find any timed wallpapers on the system")
	}
	if verbose {
		fmt.Println("Filtering wallpapers by name...")
	}
	simpleTimedWallpapers := searchResults.SimpleTimedWallpapersByName(collectionOrFilename)
	gnomeTimedWallpapers := searchResults.GnomeTimedWallpapersByName(collectionOrFilename)

	// gnomeTimedWallpapers and simpleTimedWallpapers have now been filtered so that they only contain elements with matching collection names

	if (len(gnomeTimedWallpapers) == 0) && (len(simpleTimedWallpapers) == 0) {
		return fmt.Errorf("could not find timed wallpaper: %s", collectionOrFilename)
	}

	if (len(gnomeTimedWallpapers) > 1) || (len(simpleTimedWallpapers) > 1) {
		return errors.New("found several timed backgrounds with the same name")
	}

	if len(simpleTimedWallpapers) == 1 {
		stw := simpleTimedWallpapers[0]
		if verbose {
			fmt.Printf("Using: %s\n", stw.Path)
		}
		// Start endless event loop
		if err := stw.EventLoop(verbose, func(path string) error { return wallutils.SetWallpaperVerbose(path, verbose) }, tempImageFilename); err != nil {
			return err
		}
	} else if len(gnomeTimedWallpapers) == 1 {
		gtw := gnomeTimedWallpapers[0]
		if verbose {
			fmt.Printf("Using: %s\n", gtw.Path)
		}
		// Start endless event loop
		if err := gtw.EventLoop(verbose, func(path string) error { return wallutils.SetWallpaperVerbose(path, verbose) }, tempImageFilename); err != nil {
			return err
		}
	}

	// this should never be reached
	return nil
}

func setTimedWallpaperAction(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("please provide a timed wallpaper filename as the first argument")
	}
	collectionOrFilename := c.Args().Get(0)

	// Be verbose unless a silent flag (-s) has been given
	verbose := !c.IsSet("silent")

	tempImageFilename := "/tmp/_settimed.jpg"

	err := SetTimedWallpaper(collectionOrFilename, verbose, tempImageFilename)
	if err != nil {
		// Output the capitalized error message
		msg := err.Error()
		if verbose {
			fmt.Printf("%s%s\n", strings.ToUpper(string(msg[0])), msg[1:])
		}
		// Try again, but with the "-timed" suffix
		err = SetTimedWallpaper(collectionOrFilename+"-timed", verbose, tempImageFilename)
	}
	return err
}

func main() {
	app := cli.NewApp()

	app.Name = "settimed"
	app.Usage = "start an event loop for a timed wallpaper"
	app.UsageText = "settimed [options] [path to a GNOME timed wallpaper or Simple Timed Wallpaper file]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "silent, s",
			Usage: "silence output",
		},
	}

	app.Action = setTimedWallpaperAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
