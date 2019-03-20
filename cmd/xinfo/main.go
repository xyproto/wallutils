package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
	"os"
)

func xInfoAction(c *cli.Context) error {
	// Fetch the info string
	info, err := wallutils.XInfo()
	if err != nil {
		return err
	}

	// Output the info
	fmt.Println(info)

	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "xinfo"
	app.Usage = "output information about the current X setup"
	app.UsageText = "xinfo [options]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Action = xInfoAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
