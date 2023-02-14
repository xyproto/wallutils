package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

func listTimedWallpapersAction(c *cli.Context) error {
	alsoPrintPath := c.IsSet("long")

	// Prepare to write text in columns
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	searchResults, err := wallutils.FindWallpapers()
	if err != nil {
		return err
	}

	for _, stw := range searchResults.SimpleTimedWallpapers() {
		if alsoPrintPath {
			numEvents := len(stw.Statics) + len(stw.Transitions)
			fmt.Fprintf(w, "%s\t%s\t\tevents: %d\n", stw.Name, stw.Path, numEvents)
		} else {
			fmt.Fprintf(w, "%s\n", stw.Name)
		}
	}
	for _, gw := range searchResults.GnomeTimedWallpapers() {
		if alsoPrintPath {
			numEvents := len(gw.Config.Statics) + len(gw.Config.Transitions)
			fmt.Fprintf(w, "%s\t%s\t\tevents: %d\n", gw.Name, gw.Path, numEvents)
		} else {
			fmt.Fprintf(w, "%s\n", gw.Name)
		}
	}
	w.Flush()
	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "lstimed"
	app.Usage = "list all timed wallpapers on the system"
	app.UsageText = "lstimed [options]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "long, l",
			Usage: "also list paths, and the number of timed events",
		},
	}

	app.Action = listTimedWallpapersAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
