package osipsclasses

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

// WritePidFile Write a pid file, but first make sure it doesn't exist with a running pid.
func WritePidFile(pidPath string, pidFile string) error {
	if _, err := os.Stat(pidPath); os.IsNotExist(err) {
		os.Mkdir(pidPath, 0700)
	}

	// Read in the pid file as a slice of bytes.
	if piddata, err := ioutil.ReadFile(pidPath + pidFile); err == nil {
		// Convert the file contents to an integer.
		if pid, err := strconv.Atoi(string(piddata)); err == nil {
			// Look for the pid in the process list.
			if process, err := os.FindProcess(pid); err == nil {
				// Send the process a signal zero kill.
				if err := process.Signal(syscall.Signal(0)); err == nil {
					// We only get an error if the pid isn't running, or it's not ours.
					return fmt.Errorf("pid already running: %d", pid)
				}
			}
		}
	}
	// If we get here, then the pidfile didn't exist,
	// or the pid in it doesn't belong to the user running this app.
	return ioutil.WriteFile((pidPath + pidFile), []byte(fmt.Sprintf("%d", os.Getpid())), 0664)
}

