package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

const shell = "/bin/sh"

func execCmd(format string, args ...interface{}) error {
	cmd := exec.Command(shell, "-c", fmt.Sprintf(format, args...))
	_, err := cmd.Output()
	return err
}

func execCmdGetOutput(format string, args ...interface{}) (string, error) {
	cmd := exec.Command(shell, "-c", fmt.Sprintf(format, args...))
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

func getInterface(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["itf"]
}

func responseError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("error %v", err)})
	return
}

func ipRouteAdd(itf string, ip string) error {
	ipm := fmt.Sprintf("%s/32", ip)
	return execCmd("ip -4 route add %s dev %s", ipm, itf)
}
