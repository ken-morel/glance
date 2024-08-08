package handler

import (
	// "context"
	"fmt"
	"github.com/glanceapp/glance/intern/assets"
	"github.com/glanceapp/glance/intern/glance"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

func MuxServer(a *glance.Application) (*http.ServeMux, error) {
	// TODO: add gzip support, static files must have their gzipped contents cached
	// TODO: add HTTPS support
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", a.HandlePageRequest)
	mux.HandleFunc("GET /{page}", a.HandlePageRequest)
	mux.HandleFunc("GET /api/pages/{page}/content/{$}", a.HandlePageContentRequest)
	mux.Handle("GET /static/{path...}", http.StripPrefix("/static/", glance.FileServerWithCache(http.FS(assets.PublicFS), 2*time.Hour)))

	if a.Config.Server.AssetsPath != "" {
		absAssetsPath, err := filepath.Abs(a.Config.Server.AssetsPath)

		if err != nil {
			return mux, fmt.Errorf("invalid assets path: %s", a.Config.Server.AssetsPath)
		}

		slog.Info("Serving assets", "path", absAssetsPath)
		assetsFS := glance.FileServerWithCache(http.Dir(a.Config.Server.AssetsPath), 2*time.Hour)
		mux.Handle("/assets/{path...}", http.StripPrefix("/assets/", assetsFS))
	}
	a.Config.Server.StartedAt = time.Now()

	slog.Info("Created mux server server", "host", a.Config.Server.Host, "port", a.Config.Server.Port)
	return mux, nil
}
