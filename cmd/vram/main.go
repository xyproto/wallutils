package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

func getVRAMAction(c *cli.Context) error {
	if c.IsSet("list") {
		// Retrieve a slice of GPU structs, or exit with an error, including integrated graphic cards ("VGA" in lspci output)
		gpus, err := wallutils.GPUs(true)
		if err != nil {
			return err
		}

		for i, gpu := range gpus {
			fmt.Printf("[%d] %s, %d MiB\n", i, gpu.Name, gpu.VRAM) // in MiB
		}
		return nil
	}

	// Retrieve a slice of GPU structs, or exit with an error, excluding integrated graphic cards ("VGA" in lspci output)
	gpus, err := wallutils.GPUs(false)
	if err != nil {
		return err
	}

	// Output the average VRAM in MiB
	VRAM := uint(0)
	for _, gpu := range gpus {
		VRAM += gpu.VRAM
	}
	l := uint(len(gpus))
	if l > 0 {
		VRAM /= l
	}

	// Output the average about of VRAM for all GPUs, in MiB
	fmt.Printf("%d MiB\n", VRAM)
	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "vram"
	app.Usage = "get the average VRAM for all available GPUs"
	app.UsageText = "vram [options]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "list the VRAM information for all available GPUs",
		},
	}

	app.Action = getVRAMAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
