package main

import (
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Serve the "web" folder as static files
	webDir := filepath.Join("web")

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Println("ember2D Editor running at http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
