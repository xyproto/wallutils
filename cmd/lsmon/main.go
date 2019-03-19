package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

func listMonitorAction(c *cli.Context) error {
	// Retrieve a slice of Monitor structs, or exit with an error
	monitors, err := wallutils.Monitors()
	if err != nil {
		return err
	}
	// For every monitor, output the ID, width and height
	alsoDPI := c.IsSet("dpi")
	for _, mon := range monitors {
		if alsoDPI {
			fmt.Printf("%d: %dx%d (DPI: %dx%d)\n", mon.ID, mon.Width, mon.Height, mon.DPIw, mon.DPIh)
		} else {
			fmt.Printf("%d: %dx%d\n", mon.ID, mon.Width, mon.Height)
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "lsmon"
	app.Usage = "output ID, width and height for all connected monitors"
	app.UsageText = "lsmon [options]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "l, d, dpi",
			Usage: "also output the monitor DPI",
		},
	}

	app.Action = listMonitorAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
