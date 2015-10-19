package wharfmaster

import (
	"fmt"
	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/middleware"
	"github.com/gorilla/mux"
	"github.com/stretchr/graceful"
	. "github.com/telamon/wharfmaster/models"
	_ "github.com/telamon/wharfmaster/util"
	"log"
	"net/http"
	"time"
	//	"path"
)

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
