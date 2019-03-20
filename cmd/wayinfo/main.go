package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

func wayInfoAction(c *cli.Context) error {
	// Fetch the info string
	info, err := wallutils.WaylandInfo()
	if err != nil {
		return err
	}

	// Output the info
	fmt.Println(info)

	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "wayinfo"
	app.Usage = "output information about the current Wayland setup"
	app.UsageText = "wayinfo [options]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Action = wayInfoAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
