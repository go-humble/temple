package temple

import (
	"html/template"
	"io"
	"strings"
)

var (
	Templates     = map[string]Template{}
	Partials      = map[string]Partial{}
	Layouts       = map[string]Layout{}
	PartialPrefix = "partials/"
	LayoutPrefix  = "layouts/"
)

type Template struct {
	*template.Template
}

type Partial struct {
	*template.Template
}

type Layout struct {
	*template.Template
}

type Executor interface {
	Execute(wr io.Writer, data interface{}) error
}

func (p Partial) PrefixedName() string {
	if strings.HasPrefix(p.Name(), PartialPrefix) {
		return p.Name()
	} else {
		return PartialPrefix + p.Name()
	}
}

func (l Layout) PrefixedName() string {
	if strings.HasPrefix(l.Name(), LayoutPrefix) {
		return l.Name()
	} else {
		return LayoutPrefix + l.Name()
	}
}

func AddTemplate(tmpl *template.Template) error {
	tmpl.Funcs(Funcs)
	template := Template{
		Template: tmpl,
	}
	Templates[tmpl.Name()] = template
	// Associate each partial with this template
	for _, partial := range Partials {
		if template.Lookup(partial.PrefixedName()) == nil {
			if _, err := template.AddParseTree(partial.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	// Associate each layout with this template
	for _, layout := range Layouts {
		if template.Lookup(layout.PrefixedName()) == nil {
			if _, err := template.AddParseTree(layout.PrefixedName(), layout.Tree); err != nil {
				return err
			}
		}
	}
	return nil
}

func AddPartial(tmpl *template.Template) error {
	tmpl.Funcs(Funcs)
	partial := Partial{
		Template: tmpl,
	}
	Partials[tmpl.Name()] = partial
	// Associate this partial with every template
	for _, template := range Templates {
		if template.Lookup(partial.PrefixedName()) == nil {
			if _, err := template.AddParseTree(partial.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	// Associate this partial with every other partial
	for _, other := range Partials {
		if other.Lookup(partial.PrefixedName()) == nil {
			if _, err := other.AddParseTree(partial.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	// Associate every other partial with this partial
	for _, other := range Partials {
		if partial.Lookup(partial.PrefixedName()) == nil {
			if _, err := partial.AddParseTree(partial.PrefixedName(), other.Tree); err != nil {
				return err
			}
		}
	}
	// Associate this partial with every layout
	for _, layout := range Layouts {
		if layout.Lookup(partial.PrefixedName()) == nil {
			if _, err := layout.AddParseTree(partial.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	return nil
}

func AddLayout(tmpl *template.Template) error {
	tmpl.Funcs(Funcs)
	layout := Layout{
		Template: tmpl,
	}
	Layouts[tmpl.Name()] = layout
	// Associate this layout with every template
	for _, template := range Templates {
		if template.Lookup(layout.PrefixedName()) == nil {
			if _, err := template.AddParseTree(layout.PrefixedName(), tmpl.Tree); err != nil {
				return err
			}
		}
	}
	// Associate each partial with this layout
	for _, partial := range Partials {
		if layout.Lookup(layout.PrefixedName()) == nil {
			if _, err := layout.AddParseTree(layout.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	return nil
}

func AddFunc(name string, f interface{}) {
	// Add the func the global list of funcs
	Funcs[name] = f
	// Add the func to each template, partial, and layout
	newFuncs := map[string]interface{}{
		name: f,
	}
	for _, template := range Templates {
		template.Funcs(newFuncs)
	}
	for _, partial := range Partials {
		partial.Funcs(newFuncs)
	}
	for _, layout := range Layouts {
		layout.Funcs(newFuncs)
	}
}
