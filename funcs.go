package temple

import (
	"html/template"
)

var (
	// Funcs are helper functions which can be called by all templates, partials, and layouts.
	Funcs = template.FuncMap{}
)

// AddFunc adds the function f to Funcs and makes it callable by name by all templates, partials, and layouts.
// f must follow the conventions described at http://golang.org/pkg/text/template/#FuncMap. Namely, f must have
// either a single return value, or two return values of which the second has type error. AddFunc must be called
// before AddTemplate, AddPartial, or AddLayout in order for the functions to be accessible.
func AddFunc(name string, f interface{}) {
	Funcs[name] = f
}
