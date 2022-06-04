package api

import (
	"fmt"
	"os/exec"
)

const shell = "/bin/sh"

func execCommand(command string, args ...interface{}) error {
	cmd := exec.Command(shell, "-c", fmt.Sprintf(command, args))
	_, err := cmd.Output()
	return err
}
