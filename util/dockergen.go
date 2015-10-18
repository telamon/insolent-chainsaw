package util

import (
	"fmt"
	"os/exec"
	//	"path"
)

func RegenerateConf() (bool, error) {
	cmd := exec.Command("docker-gen", "nginx.tmpl", "nginx.conf")
	out, err := cmd.Output()
	fmt.Println(out)
	if err != nil {
		return false, err
	}
	return true, nil
}
