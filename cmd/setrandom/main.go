package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func setRandomWallpaperAction(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("please specify a directory to choose wallpapers from")
	}
	dir := c.Args().Get(0)

	// Check if the verbose flag is set
	verbose := c.IsSet("verbose")

	pngMatches, err := filepath.Glob(filepath.Join(dir, "/*.png"))
	if err != nil {
		return err
	}

	jpgMatches, err := filepath.Glob(filepath.Join(dir, "/*.jpg"))
	if err != nil {
		return err
	}

	jpegMatches, err := filepath.Glob(filepath.Join(dir, "/*.jpeg"))
	if err != nil {
		return err
	}

	matches := pngMatches
	matches = append(matches, jpgMatches...)
	matches = append(matches, jpegMatches...)

	if len(matches) == 0 {
		return fmt.Errorf("found no .png, .jpg or .jpeg files in %s", dir)
	}

	imageFilename := matches[rand.Int()%len(matches)]
	if absImageFilename, err := filepath.Abs(imageFilename); err == nil {
		imageFilename = absImageFilename
	}

	if verbose {
		fmt.Printf("Setting background image to: %s\n", imageFilename)
	}
	return wallutils.SetWallpaperVerbose(imageFilename, verbose)
}

func main() {
	app := cli.NewApp()

	app.Name = "setrandom"
	app.Usage = "choose a wallpaper and set it"
	app.UsageText = "setrandom [options] [directory with images]"

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

	app.Action = setRandomWallpaperAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
