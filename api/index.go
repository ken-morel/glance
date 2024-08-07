package handler

import (
	// "context"
	"fmt"
	"github.com/glanceapp/glance/internal/glance"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	// "regexp"
	// "strings"
	// "sync"
	"time"

	"github.com/glanceapp/glance/internal/assets"
	// "github.com/glanceapp/glance/internal/widget"
)

var serverCreated bool
var server *http.ServeMux

func GetServer() *http.ServeMux {
	configFile, err := os.Open("glance.yml")

	if err != nil {
		fmt.Printf("failed opening config file: %v\n", err)
	}
	config, err := glance.NewConfigFromYml(configFile)
	configFile.Close()
	if err != nil {
		fmt.Printf("failed parsing config file: %v\n", err)
	}
	a, err := glance.NewApplication(config)

	if err != nil {
		fmt.Printf("failed creating application: %v\n", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", a.HandlePageRequest)
	mux.HandleFunc("GET /{page}", a.HandlePageRequest)
	mux.HandleFunc("GET /api/pages/{page}/content/{$}", a.HandlePageContentRequest)
	mux.Handle("GET /static/{path...}", http.StripPrefix("/static/", glance.FileServerWithCache(http.FS(assets.PublicFS), 2*time.Hour)))

	if a.Config.Server.AssetsPath != "" {
		absAssetsPath, err := filepath.Abs(a.Config.Server.AssetsPath)

		if err != nil {
			fmt.Errorf("invalid assets path: %s", a.Config.Server.AssetsPath)
		}

		slog.Info("Serving assets", "path", absAssetsPath)
		assetsFS := glance.FileServerWithCache(http.Dir(a.Config.Server.AssetsPath), 2*time.Hour)
		mux.Handle("/assets/{path...}", http.StripPrefix("/assets/", assetsFS))
	}

	// server := http.Server{
	// 	Addr:    fmt.Sprintf("%s:%d", a.Config.Server.Host, a.Config.Server.Port),
	// 	Handler: mux,
	// }
	// a.Config.Server.StartedAt = time.Now()
	// fmt.Println("Starting server", "host", a.Config.Server.Host, "port", a.Config.Server.Port)
	return mux
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if !serverCreated {
		server = GetServer()
		serverCreated = true
	}
	server.ServeHTTP(w, r)
}
