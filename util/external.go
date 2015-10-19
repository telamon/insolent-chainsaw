package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func cmdWrapper(line string) *exec.Cmd {
	var cmd *exec.Cmd
	if os.Getenv("GO_ENV") == "testing" {
		cmd = exec.Command("./exec-wrapper.sh", line)
	} else {
		s := strings.SplitN(line, " ", 2)
		cmd = exec.Command(s[0], s[1])
	}
	return cmd
}
func RegenerateConf() (bool, error) {
	cmd := cmdWrapper("docker-gen nginx.tmpl nginx.conf")
	out, err := cmd.CombinedOutput()
	fmt.Printf("Output:\n%s", out)
	if err != nil {
		return false, err
	}
	return true, nil
}
func StartNginx() error {
	cmd := cmdWrapper("nginx -c /app/nginx.conf")
	out, err := cmd.CombinedOutput()
	fmt.Printf("Output:\n%s", out)
	if err != nil {
		return err
	}
	return nil
}
func NginxPID() int {
	b, err := ioutil.ReadFile("nginx.pid")
	if err == nil {
		p, err := strconv.Atoi(strings.TrimSpace(string(b)))
		if err == nil {
			return p
		}
	}
	return -1
}
func StopNginx() {
	pid := NginxPID()
	if pid == -1 {
		return
	}
	if os.Getenv("GO_ENV") == "testing" {
		cmd := exec.Command("./exec-wrapper.sh", fmt.Sprintf("kill %d", pid))
		out, _ := cmd.CombinedOutput()
		fmt.Printf("Output:\n%s", out)
	} else {
		if proc, err := os.FindProcess(pid); err == nil {
			fmt.Printf("Sending kill signal to NGinX (%d)", pid)
			proc.Kill()
		}
	}
}
