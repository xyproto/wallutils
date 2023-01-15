package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

func getDPIAction(c *cli.Context) error {
	// Retrieve a slice of Monitor structs, or exit with an error
	monitors, err := wallutils.Monitors()
	if err != nil {
		return err
	}

	if c.IsSet("all") {
		for i, monitor := range monitors {
			fmt.Printf("[%d] %dx%d\n", i, monitor.DPIw, monitor.DPIh)
		}
		return nil
	}

	// Output the average DPI
	DPIw, DPIh := uint(0), uint(0)
	for _, monitor := range monitors {
		DPIw += monitor.DPIw
		DPIh += monitor.DPIh
	}
	DPIw /= uint(len(monitors))
	DPIh /= uint(len(monitors))

	// Check if both numbers should be outputted
	if c.IsSet("both") {
		fmt.Printf("%dx%d\n", DPIw, DPIh)
		return nil
	}

	// Only the horizontal number
	fmt.Println(DPIw)
	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "getdpi"
	app.Usage = "get the average horizontal DPI"
	app.UsageText = "getdpi [options]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a, l",
			Usage: "output DPI information for all available monitors",
		},
		cli.BoolFlag{
			Name:  "both, b",
			Usage: "output both the horizontal and vertical average DPI",
		},
	}

	app.Action = getDPIAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
