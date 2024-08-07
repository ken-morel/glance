package handler

import (
	// "context"
	"fmt"
	"github.com/glanceapp/glance/intern/glance"
	"net/http"
	"os"
)

var serverCreated bool
var mux *http.ServeMux

func GetServer() int {
	configFile, err := os.Open("glance.yml")

	if err != nil {
		fmt.Printf("failed opening config file: %v\n", err)
		return 1
	}

	config, err := glance.NewConfigFromYml(configFile)
	configFile.Close()

	if err != nil {
		fmt.Printf("failed parsing config file: %v\n", err)
		return 1
	}
	app, err := glance.NewApplication(config)

	if err != nil {
		fmt.Printf("failed creating application: %v\n", err)
		return 1
	}
	mux, err = app.MuxServer()
	if err != nil {
		fmt.Printf("http server error: %v\n", err)
		return 1
	}
	return 0
}
func Handler(w http.ResponseWriter, r *http.Request) {
	if !serverCreated {
		GetServer()
		serverCreated = true
	}
	mux.ServeHTTP(w, r)
}
