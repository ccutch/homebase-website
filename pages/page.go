package pages

import (
	html "html/template"
	"net/http"

	"github.com/gobuffalo/packr"
)

// baseTemplate to base other templates on so partials are preloaded
var baseTemplate *html.Template

func init() {
	// Loading in partial files
	baseTemplate = html.New("base")
	baseTemplate.Funcs(funcs)
}

// LoadPartials absrtacts loading from init function because of docker
// and packr issues
func LoadPartials() {
	box := packr.NewBox("../partials")
	for _, f := range box.List() {
		s := box.String(f)
		baseTemplate.Parse(s)
	}
}

// Page struct can be reference anywhere in the project, if you are in a
// different package you have to add an import but all fields are public.
type Page struct {
	// Template html file for this page
	Template string
	// Params generator function to produce params for the view based on
	// the request
	Params func(r *http.Request) interface{}
}

// ServeHTTP fulfills http.Handler interface, right now just rendering
// a template
func (p Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := html.Must(baseTemplate.Clone())

	box := packr.NewBox("../pages")
	s := box.String(p.Template)
	html.Must(t.Parse(s))

	if p.Params == nil {
		t.Execute(w, nil)
	} else {
		t.Execute(w, p.Params(r))
	}
}
