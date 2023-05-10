package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

func listGPUs() error {
	gpus, err := wallutils.GPUs(true)
	if err != nil {
		return err
	}
	for _, gpu := range gpus {
		fmt.Printf("[%s] %s, %d MiB\n", gpu.Bus, gpu.Name, gpu.VRAM)
	}
	return nil
}

func getMinimumVRAM(gpus []wallutils.GPU) uint {
	minimum := uint(0)
	for _, gpu := range gpus {
		if (minimum == 0 || gpu.VRAM < minimum) && gpu.VRAM > 0 {
			minimum = gpu.VRAM
		}
	}
	return minimum
}

func getVRAMAction(c *cli.Context) error {
	if c.IsSet("list") {
		return listGPUs()
	}

	// Get only the non-integrated GPUs (unless the integrated flag is given, then all GPUs are retrieved)
	gpus, err := wallutils.GPUs(c.IsSet("integrated"))
	if err != nil {
		return err
	}

	// If no non-integrated GPUs, try to also get the integrated GPUs, unless they are already considered
	if len(gpus) == 0 && !c.IsSet("integrated") {
		gpus, err = wallutils.GPUs(true)
		if err != nil {
			return err
		}
		if len(gpus) == 0 {
			fmt.Fprintln(os.Stderr, "error: could not find any available GPUs")
			return errors.New("could not find any available GPUs")
		}
	}

	minimum := getMinimumVRAM(gpus)

	fmt.Printf("%d MiB\n", minimum)
	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "vram"
	app.Usage = "get the minimum amount of VRAM for all non-integrated GPUs.\n          If only integrated GPUs are available, the minimum amount of VRAM for these are returned instead."
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
			Usage: "list bus ID, description and the amount of VRAM for each GPU",
		},
		cli.BoolFlag{
			Name:  "integrated, i",
			Usage: "find the minimum amount of VRAM for all GPUs, including integrated ones",
		},
	}

	app.Action = getVRAMAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
