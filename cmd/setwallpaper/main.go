package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

func setWallpaperAction(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("please specify an image filename")
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
	app.UsageText = "setwallpaper [options] [path to JPEG or PNG image]"

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

	app.Action = setWallpaperAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
