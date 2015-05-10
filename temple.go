package temple

import (
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

// Executor represents some type of template that is capable of executing (i.e. rendering)
// to an io.Writer with some data. It is satisfied by Template, Partial, and Layout as well
// as the builtin template.Template.
type Executor interface {
	Execute(wr io.Writer, data interface{}) error
}

func reset() {
	Templates = map[string]Template{}
	Partials = map[string]Partial{}
	Layouts = map[string]Layout{}
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

func AddTemplate(name, src string) error {
	tmpl, err := template.New(name).Funcs(Funcs).Parse(src)
	if err != nil {
		return err
	}
	template := Template{
		Template: tmpl,
	}
	Templates[tmpl.Name()] = template
	return associateTemplate(template)
}

func AddTemplateFile(name, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return AddTemplate(name, string(src))
}

func AddTemplateFiles(dir string) error {
	return collectTemplateFiles(dir, AddTemplateFile)
}

func associateTemplate(template Template) error {
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

func AddPartial(name, src string) error {
	tmpl, err := template.New(name).Funcs(Funcs).Parse(src)
	if err != nil {
		return err
	}
	partial := Partial{
		Template: tmpl,
	}
	Partials[tmpl.Name()] = partial
	return associatePartial(partial)
}

func AddPartialFile(name, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return AddPartial(name, string(src))
}

func AddPartialFiles(dir string) error {
	return collectTemplateFiles(dir, AddPartialFile)
}

func associatePartial(partial Partial) error {
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

func AddLayout(name, src string) error {
	tmpl, err := template.New(name).Funcs(Funcs).Parse(src)
	if err != nil {
		return err
	}
	layout := Layout{
		Template: tmpl,
	}
	Layouts[tmpl.Name()] = layout
	return associateLayout(layout)
}

func AddLayoutFile(name, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return AddLayout(name, string(src))
}

func AddLayoutFiles(dir string) error {
	return collectTemplateFiles(dir, AddLayoutFile)
}

func AddAllFiles(templatesDir, partialsDir, layoutsDir string) error {
	for dir, f := range map[string]func(string) error{
		templatesDir: AddTemplateFiles,
		partialsDir:  AddPartialFiles,
		layoutsDir:   AddLayoutFiles,
	} {
		if err := f(dir); err != nil {
			return err
		}
	}
	return nil
}

func associateLayout(layout Layout) error {
	// Associate this layout with every template
	for _, template := range Templates {
		if template.Lookup(layout.PrefixedName()) == nil {
			if _, err := template.AddParseTree(layout.PrefixedName(), layout.Tree); err != nil {
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

func collectTemplateFiles(dir string, add func(name, filename string) error) error {
	dir = filepath.Clean(dir)
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".tmpl") {
			name := strings.TrimSuffix(strings.TrimPrefix(path, dir+string(os.PathSeparator)), ".tmpl")
			if err := add(name, path); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
