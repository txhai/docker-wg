package api

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const shell = "/bin/sh"

func execCmd(command string, args ...interface{}) error {
	cmd := exec.Command(shell, "-c", fmt.Sprintf(command, args))
	_, err := cmd.Output()
	return err
}

func execCmdGetOutput(command string, args ...interface{}) (string, error) {
	cmd := exec.Command(shell, "-c", fmt.Sprintf(command, args))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func parseStr(s string) string {
	if s == "" {
		return s
	}
	if s == "(none)" {
		return ""
	}
	return strings.TrimSpace(s)
}

func parseInt64(s string) int64 {
	ps := parseStr(s)
	if ps == "" {
		return 0
	}
	i, err := strconv.ParseInt(ps, 10, 64)
	if err != nil {
		return 0
	}
	return i
}
