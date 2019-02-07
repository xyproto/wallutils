package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xyproto/monitor"
)

const versionString = "Random Wallpaper Changer 1.0.0"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	fmt.Println(versionString)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please type \"yes\" if you want to change the desktop wallpaper to a random image: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if strings.TrimSpace(text) != "yes" {
		fmt.Println("OK, made no changes.")
		os.Exit(0)
	}

	matches, err := filepath.Glob("/usr/share/pixmaps/*.png")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	imageFilename := matches[rand.Int()%len(matches)]
	if absImageFilename, err := filepath.Abs(imageFilename); err == nil {
		imageFilename = absImageFilename
	}
	if _, err := os.Stat(imageFilename); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Setting background image to: " + imageFilename)
	if err := monitor.SetWallpaper(imageFilename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
