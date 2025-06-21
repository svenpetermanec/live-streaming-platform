package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

func main() {
	// TODO: move to config

	http.Handle(
		"/", http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				fmt.Println("handling request to", req.URL)

				server := http.FileServer(http.Dir("/app/streams"))

				cleanPath, err := sanitizePath(req.URL.Path)
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

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

				// ignore preflight
				if req.Method == "OPTIONS" {
					return
				}

				fmt.Println(cleanPath, ext)

				server.ServeHTTP(w, req)
			},
		),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
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
