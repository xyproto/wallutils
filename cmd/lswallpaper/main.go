package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
	"os"
)

func listWallpapersAction(c *cli.Context) error {
	searchResults, err := wallutils.FindWallpapers()
	if err != nil {
		return err
	}
	withDetails := c.IsSet("long")
	collectionMarker := c.IsSet("star")
	for _, wp := range searchResults.Wallpapers() {
		if withDetails && collectionMarker {
			fmt.Println(wp)
		} else if withDetails {
			fmt.Printf("%dx%d\t%16s\t%s\n",
				wp.Width, wp.Height, wp.CollectionName, wp.Path)
		} else if collectionMarker {
			star := " "
			if wp.PartOfCollection {
				star = "*"
			}
			fmt.Printf("(%s) %s\n", star, wp.Path)
		} else {
			fmt.Println(wp.Path)
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "lswallpaper"
	app.Usage = "list all wallpapers on the system"
	app.UsageText = "lswallpaper [options]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "long, l",
			Usage: "also list paths",
		},
		cli.BoolFlag{
			Name:  "star, s",
			Usage: "prefix wallpapers with a star if they are part of a collection",
		},
	}

	app.Action = listWallpapersAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
