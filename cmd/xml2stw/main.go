package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils/pkg/gnometimed"
	"github.com/xyproto/wallutils"
)

func conversionAction(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("please give the path to a GNOME timed wallpaper XML file as the first argument")
	}
	filename := c.Args().Get(0)

	simpleTimedWallpaperString, err := gnometimed.GnomeFileToSimpleString(filename)
	if err != nil {
		return err
	}

	// Output the result of the conversion
	fmt.Println(simpleTimedWallpaperString)

	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "xml2stw"
	app.Usage = "convert from GNOME Timed Wallpaper to the Simple Timed Wallpaper format"
	app.UsageText = "xml2stw [options] [XML file]"

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
		//cli.StringFlag{
		//	Name:  "output, o",
		//	Value: "output.stw",
		//	Usage: "verbose output",
		//},
	}

	app.Action = conversionAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
