package monitor

import (
	"fmt"
)

func ExampleFilenameToRes() {
	fn := "hello_there_320x200.png"
	res, err := FilenameToRes(fn)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	// Output:
	// 320x200
}
