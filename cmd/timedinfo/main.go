package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

// Indent all lines with the given prefix.
// Will trim the right side of the string for newlines before indenting.
func Indent(s string, prefix string) string {
	return prefix + strings.Replace(strings.TrimRight(s, "\n"), "\n", "\n"+prefix, -1)
}

func timedInfoAction(c *cli.Context) error {
	searchResults, err := wallutils.FindWallpapers()
	if err != nil {
		return err
	}
	white := color.New(color.FgWhite, color.Bold)
	blue := color.New(color.FgBlue, color.Bold)
	gray := color.New(color.Reset)
	const prefix = "\t"

	nameFilter := ""
	if c.NArg() > 0 {
		nameFilter = c.Args().Get(0)
	}

	first := true
	for _, stw := range searchResults.SimpleTimedWallpapers() {
		if nameFilter == "" || stw.Name == nameFilter {
			if first {
				first = false
			} else {
				fmt.Println()
			}
			white.Print("Simple Timed Wallpaper: ")
			blue.Print(stw.Name)
			fmt.Println()
			gray.Println("\n" + Indent("path: "+stw.Path+"\n"+stw.String(), prefix))
		}
	}
	for _, gtw := range searchResults.GnomeTimedWallpapers() {
		if nameFilter == "" || gtw.Name == nameFilter {
			if first {
				first = false
			} else {
				fmt.Println()
			}
			white.Print("GNOME Timed Wallpaper: ")
			blue.Print(gtw.Name)
			fmt.Println()
			gray.Println("\n" + Indent(gtw.String(), prefix))
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "timedinfo"
	app.Usage = "show information about timed wallpapers on the system"
	app.UsageText = "timedinfo [options] [name]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Action = timedInfoAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
