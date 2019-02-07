package main

import (
	"fmt"
	"github.com/xyproto/monitor"
)

func main() {
	for _, name := range monitor.FindCollectionNames() {
		fmt.Println(name)
	}
}
