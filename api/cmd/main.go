package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// TODO: move to config

	http.Handle(
		"/", http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				server := http.FileServer(http.Dir("/app/streams"))
				fmt.Println(req.URL.Path, filepath.Ext(req.URL.Path))

				server.ServeHTTP(w, req)
			},
		),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
