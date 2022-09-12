package wallutils

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GPU contains an ID and the VRAM in MiB
type GPU struct {
	ID   uint   // GPU number, from 0 and up
	Name string // GPU name
	VRAM uint   // VRAM, in MiB
	VGA  bool   // Shows up as "VGA" with lspci
}

func collectNVIDIA(gpus *[]GPU) error {
	if nvidiaSMIPath := which("nvidia-smi"); nvidiaSMIPath != "" {
		lines := strings.Split(output(nvidiaSMIPath, []string{"-q"}, false), "\n")
		gpu := new(GPU)
		var lookForTotal bool
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if strings.HasPrefix(trimmedLine, "Product Name") {
				fields := strings.SplitN(trimmedLine, ":", 2)
				gpu.Name = strings.TrimSpace(fields[1])
				lookForTotal = false
			} else if strings.HasPrefix(trimmedLine, "FB Memory Usage") {
				lookForTotal = true
			} else if lookForTotal && strings.HasPrefix(trimmedLine, "Total") {
				fields := strings.SplitN(trimmedLine, ":", 2)
				amount := strings.TrimSpace(fields[1])
				if strings.HasSuffix(amount, "MiB") {
					fields = strings.SplitN(amount, " ", 2)
					if amountInt, err := strconv.Atoi(fields[0]); err == nil { // success
						gpu.VRAM = uint(amountInt)
					}
				} else if strings.HasSuffix(amount, "GiB") {
					fields = strings.SplitN(amount, " ", 2)
					if amountInt, err := strconv.Atoi(fields[0]); err == nil { // success
						gpu.VRAM = uint(amountInt * 1024)
					}
				} else if strings.HasSuffix(amount, "KiB") {
					fields = strings.SplitN(amount, " ", 2)
					if amountInt, err := strconv.Atoi(fields[0]); err == nil { // success
						gpu.VRAM = uint(math.Round(float64(amountInt) / 1024.0))
					}
				} else {
					return fmt.Errorf("unrecognized amount of memory: %s", amount)
				}
				lookForTotal = false

				nextID := uint(0)
				if l := len(*gpus); l > 0 {
					nextID = (*gpus)[l-1].ID + 1
				}
				gpu.ID = nextID

				*gpus = append(*gpus, *gpu)
				gpu = new(GPU)
			}
		}
	}
	return nil
}

func nameIndex(gpus *[]GPU, pciName string) int {
	for index, gpu := range *gpus {
		if gpu.Name == pciName {
			return index
		}
	}
	return -1
}

func collectLSPCI(gpus *[]GPU) error {
	if lspciPath := which("lspci"); lspciPath != "" {
		lines := strings.Split(output(lspciPath, []string{"-v"}, false), "\n")
		gpu := new(GPU)
		var lookForMemory bool
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if strings.Contains(trimmedLine, " VGA ") {
				fields := strings.SplitN(trimmedLine, "VGA", 2)
				pciName := strings.TrimSpace(fields[0])
				if index := nameIndex(gpus, pciName); index != -1 {
					gpu = &((*gpus)[index])
				}
				fields = strings.SplitN(fields[1], ":", 2)
				description := strings.TrimSpace(fields[1])
				if strings.Contains(description, "[") {
					fields = strings.SplitN(description, "[", 2)
					description = strings.TrimSpace(fields[0])
				}
				gpu.Name = description
				lookForMemory = true
			} else if lookForMemory && strings.HasPrefix(trimmedLine, "Memory at") {
				if strings.Contains(trimmedLine, "size=") && !strings.Contains(trimmedLine, "disabled") {
					fields := strings.SplitN(trimmedLine, "size=", 2)
					fields = strings.SplitN(fields[1], "]", 2)
					amount := fields[0]
					if strings.HasSuffix(amount, "M") {
						fields = strings.SplitN(amount, "M", 2)
						if amountInt, err := strconv.Atoi(fields[0]); err == nil { // success
							gpu.VRAM += uint(amountInt)
						}
					} else if strings.HasSuffix(amount, "K") {
						fields = strings.SplitN(amount, "K", 2)
						if amountInt, err := strconv.Atoi(fields[0]); err == nil { // success
							gpu.VRAM += uint(math.Round(float64(amountInt) / 1024.0))
						}
					} else {
						return fmt.Errorf("unrecognized amount of memory: %s", amount)
					}
				}
			} else if trimmedLine == "" {
				lookForMemory = false
				if gpu.Name != "" {
					gpu.VGA = true

					nextID := uint(0)
					if l := len(*gpus); l > 0 {
						nextID = (*gpus)[l-1].ID + 1
					}
					gpu.ID = nextID

					if nameIndex(gpus, gpu.Name) == -1 { // not already in the slice
						*gpus = append(*gpus, *gpu)
					}
					gpu = new(GPU)

				}
			}
		}
	}
	return nil
}

func collectSYSDEV(gpus *[]GPU) error {
	gpu := new(GPU)
	filepath.Walk("/sys/devices", func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		baseFilename := filepath.Base(path)
		if baseFilename != "mem_info_vram_total" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		byteAmount := strings.TrimSpace(string(data))
		if byteAmountInt, err := strconv.Atoi(byteAmount); err == nil { // success
			foundVRAM := uint(math.Round((float64(byteAmountInt) / 1024.0) / 1024.0))

			nextID := uint(0)
			if l := len(*gpus); l > 0 {
				nextID = (*gpus)[l-1].ID + 1
			}

			pciName := filepath.Dir(path)
			pciName = pciName[len(pciName)-7:]

			gpu.ID = nextID
			gpu.Name = pciName
			gpu.VRAM = foundVRAM

			*gpus = append(*gpus, *gpu)
			gpu = new(GPU)
		}

		return nil
	})
	return nil
}

// GPUs returns information about all available GPUs.
// This function will run "nvidia-smi" or any needed utility in order to collect the information.
// If alsoLSPCI is set to true, also integrated graphic cards will be detected (listed as "VGA" in the lspci output)
func GPUs(alsoLSPCI bool) ([]GPU, error) {
	gpus := make([]GPU, 0)
	if err := collectNVIDIA(&gpus); err != nil {
		return gpus, err
	}
	if err := collectSYSDEV(&gpus); err != nil {
		return gpus, err
	}
	if alsoLSPCI {
		if err := collectLSPCI(&gpus); err != nil {
			return gpus, err
		}
	}
	return gpus, nil
}

// String returns a string with GPU ID and VRAM in MiB
func (gpu GPU) String() string {
	return fmt.Sprintf("[%d] %s, %d MiB VRAM", gpu.ID, gpu.Name, gpu.VRAM)
}
