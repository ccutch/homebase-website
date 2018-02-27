package pages

import "net/http"

// Creating a instance of a page for the homepage, this can be done in a
// different file for more complexity or code organization
var Homepage = Page{
	Template: "homepage.tmpl",
	Params: func(r *http.Request) interface{} {
		return "Homebase Homepage"
	},
}
