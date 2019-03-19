package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

func setWallpaperAction(c *cli.Context) error {
	if c.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Please specify an image filename.")
		os.Exit(1)
	}
	imageFilename := c.Args().Get(0)

	// Find the absolute path
	absImageFilename, err := filepath.Abs(imageFilename)
	if err == nil {
		imageFilename = absImageFilename
	}

	// Check if the verbose flag is set
	verbose := c.IsSet("verbose")

	// Set the desktop wallpaper
	if err := wallutils.SetWallpaperVerbose(imageFilename, verbose); err != nil {
		return fmt.Errorf("could not set wallpaper: %s", err)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "setwallpaper"
	app.Usage = "change the desktop wallpaper"
	app.Version = wallutils.VersionString
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose output",
		},
	}
	app.Action = setWallpaperAction
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
