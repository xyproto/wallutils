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
	var (
		dir = c.Args().Get(0)

		// Retrieve flags from the context
		verbose   = c.IsSet("verbose")
		recursive = c.IsSet("recursive")
		mode      = c.String("mode")
		onlyLarge = c.IsSet("onlylarge")

		// Prepare to search for images
		matches []string
		err     error
	)

	if recursive {

		if verbose {
			fmt.Printf("Searching %s recursively: ", dir)
		}
		// onlyLarge means >= 640x480
		matches, err = wallutils.FindImagesAt(dir, []string{".png", ".jpg", ".jpeg"}, onlyLarge)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return err
		}
		if verbose {
			fmt.Printf("found %d images\n", len(matches))
		}

	} else {

		pngMatches, err := filepath.Glob(filepath.Join(dir, "/*.png"))
		if err != nil {
			return err
		}
		matches = append(matches, pngMatches...)

		jpgMatches, err := filepath.Glob(filepath.Join(dir, "/*.jpg"))
		if err != nil {
			return err
		}
		matches = append(matches, jpgMatches...)

		jpegMatches, err := filepath.Glob(filepath.Join(dir, "/*.jpeg"))
		if err != nil {
			return err
		}
		matches = append(matches, jpegMatches...)
	}

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
	return wallutils.SetWallpaperCustom(imageFilename, mode, verbose)
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
		cli.BoolFlag{
			Name:  "recursive, r",
			Usage: "search recursively",
		},
		cli.BoolFlag{
			Name:  "onlylarge, l",
			Usage: "only images >= 640x480 pixels",
		},
		cli.StringFlag{
			Name:  "mode, m",
			Value: "stretch", // the default value
			Usage: "wallpaper mode (stretch | center | tile | scale) \n\t+ modes specific to the currently running DE/WM",
		},
	}

	app.Action = setRandomWallpaperAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
