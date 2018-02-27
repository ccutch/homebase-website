// server.go
package main

import (
	"log"
	"net/http"

	"github.com/ccutch/homebase-website/pages"
)

// Finally setup and run the server using go's http basic package
func main() {
	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", pages.Homepage)

	log.Println("Server running at 127.0.0.1:8000")
	http.ListenAndServe(":8000", nil)
}
