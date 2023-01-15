package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/urfave/cli"
	"github.com/xyproto/heic"
	"github.com/xyproto/wallutils"
)

// Convert tries to convert a HEIC file to the STW format
func Convert(filename string) (string, error) {
	ctx, err := heic.NewContext()
	if err != nil {
		return "", err
	}
	if err := ctx.ReadFromFile(filename); err != nil {
		return "", err
	}
	if count := ctx.GetNumberOfTopLevelImages(); count == 0 {
		return "", errors.New("0 top level images")
	}
	if ids := ctx.GetListOfTopLevelImageIDs(); len(ids) == 0 {
		return "", errors.New("0 top level image IDs")
	}
	if _, err := ctx.GetPrimaryImageID(); err != nil {
		return "", errors.New("no primary image")
	}
	handle, err := ctx.GetPrimaryImageHandle()
	if err != nil {
		return "", fmt.Errorf("could not get primary image handle: %s", err)
	}
	if !handle.IsPrimaryImage() {
		return "", errors.New("primary image handle is not for the primary image")
	}
	if metadataCount := handle.MetadataCount(); metadataCount == 0 {
		return "", errors.New("no dynamic wallpaper metadata in " + filename)
	}
	metadataIDs := handle.MetadataIDs()
	if len(metadataIDs) == 0 {
		return "", errors.New("no metadata IDs")
	}
	firstMetadataID := metadataIDs[0]

	timeTable, err := handle.ImageTimes(firstMetadataID)
	if err != nil {
		return "", err
	}

	name := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	if len(name) == 0 {
		return "", errors.New("image name is empty")
	}

	s := fmt.Sprintf("stw: 1.0\nname: %s\n", name)
	s += fmt.Sprintf("format: /usr/share/backgrounds/%s/%%s.jpg\n", name)

	var lines []string
	for i, t := range timeTable {
		lines = append(lines, fmt.Sprintf("@%02d:%02d: %02d", t.Hour(), t.Minute(), i))
	}
	sort.Strings(lines)

	return s + strings.Join(lines, "\n"), nil
}

func conversionAction(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("please give the path to a HEIC dynamic wallpaper file as the first argument")
	}
	filename := c.Args().Get(0)

	simpleTimedWallpaperString, err := Convert(filename)
	if err != nil {
		return err
	}

	// Output the result of the conversion
	fmt.Println(simpleTimedWallpaperString)
	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "heic2stw"
	app.Usage = "convert from GNOME Timed Wallpaper to the Simple Timed Wallpaper format"
	app.UsageText = "heic2stw [options] [HEIC file]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Action = conversionAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
