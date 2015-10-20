package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func RegenerateConf() (bool, error) {
	cmd := exec.Command("docker-gen", "vhosts.tmpl", "vhosts.conf")
	out, err := cmd.CombinedOutput()
	fmt.Printf("Output:\n%s", out)
	if err != nil {
		return false, err
	}
	return true, nil
}

func NginxPID() (int, error) {
	b, err := ioutil.ReadFile("nginx.pid")
	if err == nil {
		p, err := strconv.Atoi(strings.TrimSpace(string(b)))
		if err == nil {
			return p, nil
		}
	}
	return -1, err
}
func nginxProc() (*os.Process, error) {
	pid, err := NginxPID()
	if pid == -1 {
		return nil, err
	}
	if proc, err := os.FindProcess(pid); err == nil {
		return proc, nil
	}
	return nil, err
}

func ReloadNginx() error {
	cmd := exec.Command("sh", "-c", "nginx -s reload")
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	return err
}
