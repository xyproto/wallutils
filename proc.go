package wallutils

import (
	"os"
	"os/exec"
	"strings"
)

// RunningWithArgs checks if the given command is running with the specified args
func RunningWithArgs(command, args string) (bool, error) {
	// First try reading /proc directly (Linux)
	if found, err := checkProcFSForArgs(command, args); err == nil {
		return found, nil
	}
	// Fallback to pgrep (works on both Linux and FreeBSD)
	output, err := exec.Command("pgrep", "-f", command+".*"+args).Output()
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(string(output))) > 0, nil
}

func checkProcFSForArgs(command, args string) (bool, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Check if directory name is numeric (PID)
		pid := entry.Name()
		if len(pid) == 0 || pid[0] < '0' || pid[0] > '9' {
			continue
		}
		cmdlineFile := "/proc/" + pid + "/cmdline"
		cmdline, err := os.ReadFile(cmdlineFile)
		if err != nil {
			continue
		}
		// cmdline uses null bytes as separators
		cmdArgs := strings.Split(string(cmdline), "\x00")
		if len(cmdArgs) < 2 {
			continue
		}
		// Check if it's the target command with the specified args
		if strings.Contains(cmdArgs[0], command) {
			for _, arg := range cmdArgs[1:] {
				if arg == args {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
