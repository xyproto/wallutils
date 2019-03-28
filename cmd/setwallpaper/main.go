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

// download can download a file to the given filename.
// Set redownload to true for downloading again even if it exists.
func download(url, filename string, verbose, redownload bool) error {
	// Check if the file exists (and that force is not enabled)
	if exists(filename) && !redownload {
		// The file already exists. This is fine, skip the download
		return nil
	}
	// Prepare the client
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if verbose {
		if verbose {
			fmt.Println("Downloading " + url)
		}
	}
	// Download the file
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
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
	downloadDir := c.String("download")

	if !exists(downloadDir) {
		return errors.New("could not find: " + downloadDir)
	}

	// Check if the argument is an URL that uses the http or https protocol
	if strings.HasPrefix(imageFilename, "http://") || strings.HasPrefix(imageFilename, "https://") {
		u, err := url.Parse(imageFilename)
		if err == nil { // no error
			downloadFilename := filepath.Join(downloadDir, filepath.Base(imageFilename))
			if err := download(u.String(), downloadFilename, verbose, false); err != nil {
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
		cli.StringFlag{
			Name:  "download, d",
			Value: "/tmp",
			Usage: "download directory",
		},
	}

	app.Action = setWallpaperAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
