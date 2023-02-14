package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

// has checks if the given string slice contains the given string
func has(sl []string, s string) bool {
	for _, e := range sl {
		if e == s {
			return true
		}
	}
	return false
}

func listWallpaperCollectionAction(c *cli.Context) error {
	alsoPrintPath := c.IsSet("long")

	// Find all wallpapers
	searchResults, err := wallutils.FindWallpapers()
	if err != nil {
		return err
	}

	if !alsoPrintPath {
		for _, name := range searchResults.CollectionNames() {
			fmt.Println(name)
		}
		return nil
	}

	// Prepare to write text in columns
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', tabwriter.AlignRight)

	// Output all wallpaper collection names and paths (these are directories
	// with files of varying resolutions)
	var collectionNames []string
	for _, wp := range searchResults.Wallpapers() {
		if wp.PartOfCollection {
			name := wp.CollectionName
			dir := filepath.Dir(wp.Path) + "/"
			if alsoPrintPath && !has(collectionNames, name) {
				fmt.Fprintf(w, "%s\t%s\t\t%s\n", name, "Wallpaper Collection", dir)
				collectionNames = append(collectionNames, wp.CollectionName)
			}
		}
	}

	// Timed wallpapers is a collection in the sense that it may point to
	// several wallpaper images.

	// Output all Simple Timed Wallpaper names and paths.
	collectionNames = []string{}
	for _, stw := range searchResults.SimpleTimedWallpapers() {
		name := stw.Name
		path := stw.Path
		if alsoPrintPath && !has(collectionNames, name) {
			fmt.Fprintf(w, "%s\t%s\t\t%s\n", name, "Simple Timed Wallpaper", path)
			collectionNames = append(collectionNames, name)
		}
	}

	// Output all GNOME timed wallpaper names and paths.
	collectionNames = []string{}
	for _, gw := range searchResults.GnomeTimedWallpapers() {
		name := gw.Name
		path := gw.Path
		if alsoPrintPath && !has(collectionNames, name) {
			fmt.Fprintf(w, "%s\t%s\t\t%s\n", name, "GNOME Timed Wallpaper", path)
			collectionNames = append(collectionNames, name)
		}
	}

	// Write the output to stdout
	w.Flush()

	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "lscolletion"
	app.Usage = "list all wallpaper collections on the system"
	app.UsageText = "lscollection [options]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "long, l",
			Usage: "also list collection type and full path",
		},
	}

	app.Action = listWallpaperCollectionAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
