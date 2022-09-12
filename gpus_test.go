package wallutils

import (
	"fmt"
	"testing"
)

func TestGPUs(t *testing.T) {
	gpus, err := GPUs(true)
	if err != nil {
		t.Error(err)
	}
	for _, gpu := range gpus {
		fmt.Println(gpu)
	}
}
