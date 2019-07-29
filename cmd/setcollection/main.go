package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

// Select the wallpaper that is closest to the current monitor resolution and set that as the wallpaper
func SelectAndSetWallpaper(wallpapers []*wallutils.Wallpaper) error {
	// Gather a slice of filenames
	var filenames []string
	for _, wp := range wallpapers {
		filenames = append(filenames, wp.Path)
	}

	// Select the image filename that is closest to the current monitor resolution
	imageFilename, err := wallutils.Closest(filenames)
	if err != nil {
		return err
	}

	// Find the absolute path
	absImageFilename, err := filepath.Abs(imageFilename)
	if err == nil {
		imageFilename = absImageFilename
	}

	// Set the desktop wallpaper
	if err := wallutils.SetWallpaper(imageFilename); err != nil {
		return fmt.Errorf("could not set wallpaper: %s", err)
	}

	return nil
}

func setWallpaperCollectionAction(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("please specify a wallpaper collection name")
	}
	collectionName := c.Args().Get(0)

	verbose := c.IsSet("verbose")

	if verbose {
		fmt.Printf("Setting wallpaper collection \"%s\"\n", collectionName)
		fmt.Print("Searching for wallpapers...")
	}

	searchResults, err := wallutils.FindWallpapers()
	if err != nil {
		return err
	}

	if searchResults.Empty() {
		return errors.New("could not find any wallpapers on the system")
	}

	if verbose {
		fmt.Println("ok")
		fmt.Print("Filtering wallpapers by collection name...")
	}

	wallpapers := searchResults.WallpapersByName(collectionName)
	gnomeTimedWallpapers := searchResults.GnomeTimedWallpapersByName(collectionName)
	simpleTimedWallpapers := searchResults.SimpleTimedWallpapersByName(collectionName)

	if verbose {
		fmt.Println("ok")
	}

	if len(wallpapers) == 0 && (len(gnomeTimedWallpapers) > 0 || len(simpleTimedWallpapers) > 0) {
		return errors.New("timed wallpapers are not supported by this utility, please try \"settimed\" instead")
	}

	if len(wallpapers) == 0 {
		return fmt.Errorf("no such collection: %s", collectionName)
	}

	return SelectAndSetWallpaper(wallpapers)
}

func main() {
	app := cli.NewApp()

	app.Name = "setcollection"
	app.Usage = "change the desktop wallpaper"
	app.UsageText = "setcollection [options] [name of wallpaper collection]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "verbose output",
		},
	}

	app.Action = setWallpaperCollectionAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
