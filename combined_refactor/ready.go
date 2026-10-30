package main

import (
	"fmt"
	"os"
)

// Publish only after binding. Android supplies a unique private path per launch.
func writeReadyPort(path string, port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port: %d", port)
	}
	temp := path + ".tmp"
	if err := os.WriteFile(temp, []byte(fmt.Sprint(port)), 0600); err != nil {
		return err
	}
	defer os.Remove(temp)
	return os.Rename(temp, path)
}
