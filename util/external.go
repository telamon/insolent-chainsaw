package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
func NginxWorkerPids() ([]int, error) {
	master, err := NginxPID()
	if err != nil {
		return nil, err
	}
	pids, err := psPids()
	if err != nil {
		return nil, err
	}
	workers := make([]int, 0, len(pids))
	for _, pid := range pids {
		if pid != master {
			workers = append(workers, pid)
		}
	}
	return workers, nil
}
func psPids() ([]int, error) {
	cmd := exec.Command("sh", "-c", "pgrep nginx")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	spids := strings.Split(strings.TrimSpace(string(out)), "\n")
	pids := make([]int, 0, len(spids))
	for _, pid := range spids {
		p, err := strconv.Atoi(strings.TrimSpace(pid))
		if err != nil {
			return nil, err
		}
		pids = append(pids, p)
	}
	return pids, err
}

func NginxReloaded(workers []int) (bool, error) {
	workersNew, err := NginxWorkerPids()
	if err != nil {
		return false, err
	}
	for _, op := range workers {
		for _, np := range workersNew {
			if op == np {
				return false, nil
			}
		}
	}
	return true, nil
}

func ReloadNginx() error {
	workers, err := NginxWorkerPids()
	if err != nil {
		return err
	}
	cmd := exec.Command("sh", "-c", "nginx -s reload")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Errorf("%s", out)
		return err
	}
	for r := false; !r; {
		time.Sleep(200 * time.Millisecond)
		r, err = NginxReloaded(workers)
		if err != nil {
			return err
		}
	}

	return err
}
