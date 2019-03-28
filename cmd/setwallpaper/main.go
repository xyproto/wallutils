package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

// exists checks if the given path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// download can download a file to the given filename
// Use force if existing files should be overwritten.
func download(url, filename string, force bool) error {
	// Prepare the client
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Download the file
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// Check if the file exists (and that force is not enabled)
	if exists(filename) && !force {
		return fmt.Errorf("%s already exists", filename)
	}
	// Write the file
	return ioutil.WriteFile(filename, b, 0644)
}

func setWallpaperAction(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("please specify an image filename or URL")
	}
	imageFilename := c.Args().Get(0)

	// Retrieve flags from the context
	verbose := c.IsSet("verbose")
	mode := c.String("mode")

	// Check if the argument is an URL that uses the http or https protocol
	if strings.HasPrefix(imageFilename, "http://") || strings.HasPrefix(imageFilename, "https://") {
		u, err := url.Parse(imageFilename)
		if err == nil { // no error
			// TODO: Use a function for getting the temp directory
			downloadFilename := filepath.Join("/tmp/", filepath.Base(imageFilename))
			if verbose {
				fmt.Println("Downloading " + u.String())
			}
			if err := download(u.String(), downloadFilename, true); err != nil {
				return err
			}
			// Use the downloaded image
			imageFilename = downloadFilename
		}
	}

	// Find the absolute path
	absImageFilename, err := filepath.Abs(imageFilename)
	if err == nil {
		imageFilename = absImageFilename
	}

	// Set the desktop wallpaper
	if err := wallutils.SetWallpaperCustom(imageFilename, verbose, mode); err != nil {
		return fmt.Errorf("could not set wallpaper: %s", err)
	}
	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "setwallpaper"
	app.Usage = "set the desktop wallpaper"
	app.UsageText = "setwallpaper [options] [path or URL to JPEG or PNG image]"

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
		cli.StringFlag{
			Name:  "mode",
			Value: "fill",
			Usage: "wallpaper mode (fill | center | scale | tile) + modes specific to the currently running DE/WM",
		},
	}

	app.Action = setWallpaperAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
