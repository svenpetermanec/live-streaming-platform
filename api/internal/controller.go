package internal

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"transcoder/api/cmd/config"
	"transcoder/internal/logging"
)

type Controller struct {
	fileServer http.Handler
	config     config.Config
	logger     *logging.Logger
}

func NewController(config config.Config, logger *logging.Logger) *Controller {
	return &Controller{
		fileServer: http.StripPrefix("/streams", http.FileServer(http.Dir(config.HLS.OutputDir))),
		config:     config,
		logger:     logger,
	}
}

func (c *Controller) ServeStream(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Handling ServeStream", logging.Data{"path": r.URL.Path}) // TODO: middleware

	cleanPath, err := sanitizePath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(cleanPath)
	switch ext {
	case ".m3u8":
		w.Header().Set("Content-Type", "application/cnd.apple.mpegurl")
	case ".ts":
		w.Header().Set("Content-Type", "video/mp2t")
	default:
		http.Error(w, "File type not allowed", http.StatusForbidden)
		return
	}

	// middleware
	w.Header().Set("Access-Control-Allow-Origin", c.config.HTTPServer.CORSOrigin)
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.config.HTTPServer.CORSMethods, ","))

	// ignore preflight
	if r.Method == "OPTIONS" {
		return
	}

	c.fileServer.ServeHTTP(w, r)
}

func (c *Controller) CreateAccount(w http.ResponseWriter, r *http.Request) {
	// TODO: generate stream key
}

func sanitizePath(rawPath string) (string, error) {
	decoded, err := url.QueryUnescape(rawPath)
	if err != nil {
		return "", err
	}

	clean := filepath.Clean(decoded)

	if strings.Contains(clean, "..") || strings.Contains(decoded, "..") {
		return "", fmt.Errorf("invalid path: %s", rawPath)
	}

	return clean, nil
}
