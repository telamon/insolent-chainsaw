package wharfmaster

import (
	"fmt"
	. "github.com/telamon/wharfmaster/models"
	//"os"
	"os/exec"
	//	"path"
)

type WharfMaster struct {
	Services []Service
}

func New() *WharfMaster {

	return &WharfMaster{}
}

func RegenerateConf() (bool, error) {
	template := "./nginx.tmpl"
	target := "./nginx.conf"
	cmd := exec.Command("docker-gen", template, target)
	out, err := cmd.Output()
	fmt.Println(out)
	if err != nil {
		return false, err
	}
	return true, nil
}

var nginx *exec.Cmd

func StartNginx() (*exec.Cmd, error) {
	daemon := "off"
	workers := 2
	pidfile := "/app/nginx.pid"
	cmd := exec.Command("nginx", "-c", "nginx.conf", fmt.Sprintf("-g \"daemon %s;worker_processes %d;pid %s;\"", daemon, workers, pidfile))
	nginx = cmd
	return cmd, cmd.Start()
}

func StopNginx() {
	nginx.Process.Kill()
}
