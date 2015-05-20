// Copyright 2015 Alex Browne.
// All rights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

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

// A Group represents a set of associated Templates, Partials, and Layouts.
// Each Group also gets its own template.FuncMap called Funcs.
type Group struct {
	Templates map[string]Template
	Partials  map[string]Partial
	Layouts   map[string]Layout
	// Funcs is a map of function names to functions. All functions in the
	// FuncMap are accessible by all Templates, Partials, and Layouts for
	// this Group.
	Funcs template.FuncMap
}

// Executor represents some type of template that is capable of executing (i.e. rendering)
// to an io.Writer with some data. It is satisfied by Template, Partial, and Layout as well
// as the builtin template.Template.
type Executor interface {
	Execute(wr io.Writer, data interface{}) error
}

// NewGroup creates, initializes, and returns a new Group
func NewGroup() *Group {
	return &Group{
		Templates: map[string]Template{},
		Partials:  map[string]Partial{},
		Layouts:   map[string]Layout{},
		Funcs:     template.FuncMap{},
	}
}

func (g *Group) AddFunc(name string, f interface{}) {
	g.Funcs[name] = f
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

func (g *Group) AddTemplate(name, src string) error {
	tmpl, err := template.New(name).Funcs(g.Funcs).Parse(src)
	if err != nil {
		return err
	}
	template := Template{
		Template: tmpl,
	}
	g.Templates[tmpl.Name()] = template
	return g.associateTemplate(template)
}

func (g *Group) AddTemplateFile(name, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return g.AddTemplate(name, string(src))
}

func (g *Group) AddTemplateFiles(dir string) error {
	return collectTemplateFiles(dir, g.AddTemplateFile)
}

func (g *Group) associateTemplate(template Template) error {
	// Associate each partial with this template
	for _, partial := range g.Partials {
		if template.Lookup(partial.PrefixedName()) == nil {
			if _, err := template.AddParseTree(partial.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	// Associate each layout with this template
	for _, layout := range g.Layouts {
		if template.Lookup(layout.PrefixedName()) == nil {
			if _, err := template.AddParseTree(layout.PrefixedName(), layout.Tree); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Group) AddPartial(name, src string) error {
	tmpl, err := template.New(name).Funcs(g.Funcs).Parse(src)
	if err != nil {
		return err
	}
	partial := Partial{
		Template: tmpl,
	}
	g.Partials[tmpl.Name()] = partial
	return g.associatePartial(partial)
}

func (g *Group) AddPartialFile(name, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return g.AddPartial(name, string(src))
}

func (g *Group) AddPartialFiles(dir string) error {
	return collectTemplateFiles(dir, g.AddPartialFile)
}

func (g *Group) associatePartial(partial Partial) error {
	// Associate this partial with every template
	for _, template := range g.Templates {
		if template.Lookup(partial.PrefixedName()) == nil {
			if _, err := template.AddParseTree(partial.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	// Associate this partial with every other partial
	for _, other := range g.Partials {
		if other.Lookup(partial.PrefixedName()) == nil {
			if _, err := other.AddParseTree(partial.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	// Associate every other partial with this partial
	for _, other := range g.Partials {
		if partial.Lookup(partial.PrefixedName()) == nil {
			if _, err := partial.AddParseTree(partial.PrefixedName(), other.Tree); err != nil {
				return err
			}
		}
	}
	// Associate this partial with every layout
	for _, layout := range g.Layouts {
		if layout.Lookup(partial.PrefixedName()) == nil {
			if _, err := layout.AddParseTree(partial.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Group) AddLayout(name, src string) error {
	tmpl, err := template.New(name).Funcs(g.Funcs).Parse(src)
	if err != nil {
		return err
	}
	layout := Layout{
		Template: tmpl,
	}
	g.Layouts[tmpl.Name()] = layout
	return g.associateLayout(layout)
}

func (g *Group) AddLayoutFile(name, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return g.AddLayout(name, string(src))
}

func (g *Group) AddLayoutFiles(dir string) error {
	return collectTemplateFiles(dir, g.AddLayoutFile)
}

func (g *Group) AddAllFiles(templatesDir, partialsDir, layoutsDir string) error {
	for dir, f := range map[string]func(string) error{
		templatesDir: g.AddTemplateFiles,
		partialsDir:  g.AddPartialFiles,
		layoutsDir:   g.AddLayoutFiles,
	} {
		if err := f(dir); err != nil {
			return err
		}
	}
	return nil
}

func (g *Group) associateLayout(layout Layout) error {
	// Associate this layout with every template
	for _, template := range g.Templates {
		if template.Lookup(layout.PrefixedName()) == nil {
			if _, err := template.AddParseTree(layout.PrefixedName(), layout.Tree); err != nil {
				return err
			}
		}
	}
	// Associate each partial with this layout
	for _, partial := range g.Partials {
		if layout.Lookup(layout.PrefixedName()) == nil {
			if _, err := layout.AddParseTree(layout.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	return nil
}

func collectTemplateFiles(dir string, handler func(name, filename string) error) error {
	dir = filepath.Clean(dir)
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".tmpl") {
			name := strings.TrimSuffix(strings.TrimPrefix(path, dir+string(os.PathSeparator)), ".tmpl")
			if err := handler(name, path); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
