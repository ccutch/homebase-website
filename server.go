// server.go
package main

import (
	"log"
	"net/http"

	"github.com/ccutch/goparcel"
	"github.com/ccutch/homebase-website/pages"
	"github.com/gobuffalo/packr"
)

// Finally setup and run the server using go's http basic package
func main() {
	worker, err := goparcel.Open("./dist", "frontend/entry.js")
	worker.EnsureWorkspace(err)
	worker.PublicURL = "/dist/"
	go worker.Start()

	pages.LoadPartials()

	fs := http.FileServer(packr.NewBox("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/dist/", worker.FileServer())
	http.Handle("/", pages.Homepage)

	log.Println("Server running at http://127.0.0.1:8000")
	http.ListenAndServe("0.0.0.0:8000", nil)
}
