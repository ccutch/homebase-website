package pages

import (
	html "html/template"
)

type User struct {
	Name string
}

var funcs = html.FuncMap{
	"getUser": func() User {
		return User{
			Name: "Connor",
		}
	},
}
