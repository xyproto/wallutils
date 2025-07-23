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
	Name string // GPU name
	Bus  string // ie. 01:00.0
	ID   uint   // GPU number, from 0 and up
	VRAM uint   // VRAM, in MiB
	VGA  bool   // Integrated / shows up as "VGA" with lspci
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
			} else if strings.HasPrefix(trimmedLine, "Bus Id") {
				fields := strings.SplitN(trimmedLine, ":", 2)
				bus := strings.TrimSpace(fields[1])
				if strings.Count(bus, ":") == 2 {
					busFields := strings.SplitN(bus, ":", 2)
					bus = busFields[1]
				}
				gpu.Bus = bus
			} else if strings.HasPrefix(trimmedLine, "FB Memory Usage") {
				lookForTotal = true
			} else if lookForTotal && strings.HasPrefix(trimmedLine, "Total") {
				fields := strings.SplitN(trimmedLine, ":", 2)
				amount := strings.TrimSpace(fields[1])
				if strings.HasSuffix(amount, "KiB") {
					fields = strings.SplitN(amount, " ", 2)
					if amountInt, err := strconv.Atoi(fields[0]); err == nil { // success
						gpu.VRAM = uint(math.Round(float64(amountInt) / 1024.0))
					}
				} else if strings.HasSuffix(amount, "MiB") {
					fields = strings.SplitN(amount, " ", 2)
					if amountInt, err := strconv.Atoi(fields[0]); err == nil { // success
						gpu.VRAM = uint(amountInt)
					}
				} else if strings.HasSuffix(amount, "GiB") {
					fields = strings.SplitN(amount, " ", 2)
					if amountInt, err := strconv.Atoi(fields[0]); err == nil { // success
						gpu.VRAM = uint(amountInt * 1024)
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

func busIndex(gpus *[]GPU, bus string) int {
	for index, gpu := range *gpus {
		if gpu.Bus == bus {
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
		var alreadyThere bool
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if strings.Contains(trimmedLine, " VGA ") {
				fields := strings.SplitN(trimmedLine, "VGA", 2)
				pciName := strings.TrimSpace(fields[0])
				gpu.Bus = pciName
				if index := busIndex(gpus, pciName); index != -1 {
					gpu = &((*gpus)[index])
					gpu.VGA = true
				}
				fields = strings.SplitN(fields[1], ":", 2)
				description := strings.TrimSpace(fields[1])
				if strings.Contains(description, "(") {
					fields = strings.SplitN(description, "(", 2)
					description = strings.TrimSpace(fields[0])
				}
				gpu.Name = description
				alreadyThere = busIndex(gpus, gpu.Bus) != -1
				lookForMemory = true
			} else if lookForMemory && strings.HasPrefix(trimmedLine, "Memory at") && !alreadyThere {
				if strings.Contains(trimmedLine, "size=") && !strings.Contains(trimmedLine, "disabled") {
					fields := strings.SplitN(trimmedLine, "size=", 2)
					fields = strings.SplitN(fields[1], "]", 2)
					amount := fields[0]
					if strings.HasSuffix(amount, "K") {
						fields = strings.SplitN(amount, "K", 2)
						if amountInt, err := strconv.Atoi(fields[0]); err == nil { // success
							gpu.VRAM += uint(math.Round(float64(amountInt) / 1024.0))
						}
					} else if strings.HasSuffix(amount, "M") {
						fields = strings.SplitN(amount, "M", 2)
						if amountInt, err := strconv.Atoi(fields[0]); err == nil { // success
							gpu.VRAM += uint(amountInt)
						}
					} else if strings.HasSuffix(amount, "G") {
						fields = strings.SplitN(amount, "G", 2)
						if amountInt, err := strconv.Atoi(fields[0]); err == nil { // success
							gpu.VRAM += uint(amountInt * 1024)
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

					if busIndex(gpus, gpu.Bus) == -1 { // not already in the slice
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
	const sysDevPath = "/sys/devices"

	// Skip the rest of this function for systems that does not have /sys/devices
	if fi, err := os.Stat(sysDevPath); err != nil || !fi.IsDir() {
		return nil
	}

	gpu := new(GPU)
	filepath.Walk(sysDevPath, func(path string, fi os.FileInfo, err error) error {
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
			gpu.Bus = pciName

			*gpus = append(*gpus, *gpu)
			gpu = new(GPU)
		}

		return nil
	})
	return nil
}

func nonIntegrated(gpus *[]GPU) {
	var igpus []GPU
	for i := range *gpus {
		if strings.Contains((*gpus)[i].Name, "UHD") {
			(*gpus)[i].VGA = true
		}
		if !(*gpus)[i].VGA {
			igpus = append(igpus, (*gpus)[i])
		}
	}
	*gpus = igpus
}

// GPUs returns information about all available GPUs.
// This function will run "nvidia-smi" or any needed utility in order to collect the information.
// If alsoIntegrated is set to true, also integrated graphic cards will be detected (listed as "VGA" in the lspci output, or contains the string "UHD")
func GPUs(alsoIntegrated bool) ([]GPU, error) {
	gpus := make([]GPU, 0)
	if err := collectNVIDIA(&gpus); err != nil {
		if !alsoIntegrated {
			nonIntegrated(&gpus)
		}
		return gpus, err
	}
	if err := collectSYSDEV(&gpus); err != nil {
		if !alsoIntegrated {
			nonIntegrated(&gpus)
		}
		return gpus, err
	}
	if alsoIntegrated {
		if err := collectLSPCI(&gpus); err != nil {
			return gpus, err
		}
	}
	if !alsoIntegrated {
		nonIntegrated(&gpus)
	}
	return gpus, nil
}

// String returns a string with GPU ID and VRAM in MiB
func (gpu GPU) String() string {
	return fmt.Sprintf("[%d] %s, %d MiB VRAM, bus %s", gpu.ID, gpu.Name, gpu.VRAM, gpu.Bus)
}
