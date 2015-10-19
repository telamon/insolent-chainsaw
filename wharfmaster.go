package wharfmaster

import (
	"fmt"
	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/middleware"
	"github.com/gorilla/mux"
	"github.com/stretchr/graceful"
	. "github.com/telamon/wharfmaster/models"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
	//	"path"
	"io/ioutil"
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

type WharfMaster struct {
	Services   []Service
	Middle     *interpose.Middleware
	HttpServer *graceful.Server
}

func (this *WharfMaster) Start(port int) {
	this.HttpServer = &graceful.Server{
		Timeout: 5 * time.Second,
		Server:  &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: this.Middle},
	}
	fmt.Println("Server listening on port:", port)
	err := this.HttpServer.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (this *WharfMaster) Stop() {
	defer fmt.Println("Server stopped")
	this.HttpServer.Stop(5 * time.Second)
}

func New() *WharfMaster {
	middle := interpose.New()

	// Logger
	middle.Use(middleware.GorillaLog())
	// Header modification example
	middle.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("X-Server-Name", "The Germinator")
			next.ServeHTTP(rw, req)
		})
	})
	// Inject router
	router := mux.NewRouter()
	middle.UseHandler(router)

	// Setup Routes
	router.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rw, "Ther once was a happy roach, he liked to do roachy stuff")
	})
	//router.HandleFunc("/v1/entities", ctrl.EntityCreate).Methods("POST")
	//router.HandleFunc("/v1/entities", ctrl.EntityList).Methods("GET")

	return &WharfMaster{Middle: middle}
}
