// Copyright 2015 Alex Browne.
// All rights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package temple

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	// PartialPrefix is added to the name of all partials.
	PartialPrefix = "partials/"
	// LayoutPrefix is added to the name of all layouts.
	LayoutPrefix = "layouts/"
)

// Template is a lightweight wrapper around template.Template
// from the builtin html/template package. It has all the same
// methods, plus some additional ones. All partials and layouts
// are associated with all templates, i.e. added to the parse
// tree for all templates. So you can render partials and layouts
// inside of a template with the `template` action. See
// http://golang.org/pkg/text/template/ for more information
// about nested templates and the `template` action.
type Template struct {
	*template.Template
}

// Partial is a lightweight wrapper around template.Template
// from the builtin html/template package. It has all the same
// methods, plus some additional ones. Partials are associated
// with all other partials, as well as with regular templates
// and layouts. That means you can render partials inside of
// templates, layouts, or other partials with the
// `template` action. See http://golang.org/pkg/text/template/
// for more information about nested templates and the `template`
// action.
type Partial struct {
	*template.Template
}

// Layout is a lightweight wrapper around template.Template
// from the builtin html/template package. It has all the same
// methods, plus some additional ones. Layouts are assiciated with
// all regular temlates, and partials are associated with layouts.
// That means you can render a partial inside of a lyout with
// the `template` action. You can also render a layout from inside
// a regular template with the `template` action. See
// http://golang.org/pkg/text/template/ for more information
// about nested templates and the `template` action.
type Layout struct {
	*template.Template
}

// A Group represents a set of associated templates, partials, and layouts.
// Each Group also gets its own template.FuncMap called Funcs.
type Group struct {
	templates map[string]*Template
	partials  map[string]*Partial
	layouts   map[string]*Layout
	// Funcs is a map of function names to functions. All functions in the
	// FuncMap are accessible by all templates, partials, and layouts for
	// this Group.
	Funcs template.FuncMap
}

// GetTemplate returns the template identified by name, or an error if
// the template could not be found.
func (g Group) GetTemplate(name string) (*Template, error) {
	template, found := g.templates[name]
	if !found {
		return nil, fmt.Errorf("Could not find template named %s", name)
	}
	return template, nil
}

// GetPartial returns the partial identified by name, or an error if
// the partial could not be found.
func (g Group) GetPartial(name string) (*Partial, error) {
	partial, found := g.partials[name]
	if !found {
		return nil, fmt.Errorf("Could not find partial named %s", name)
	}
	return partial, nil
}

// GetLayout returns the layout identified by name, or an error if
// the layout could not be found.
func (g Group) GetLayout(name string) (*Layout, error) {
	layout, found := g.layouts[name]
	if !found {
		return nil, fmt.Errorf("Could not find layout named %s", name)
	}
	return layout, nil
}

// MustGetTemplate works like GetTemplate, except that it panics
// instead of returning an error if the template could not be
// found.
func (g Group) MustGetTemplate(name string) *Template {
	template, found := g.templates[name]
	if !found {
		panic("Could not find template named " + name)
	}
	return template
}

// MustGetPartial works like GetPartial, except that it panics
// instead of returning an error if the partial could not be
// found.
func (g Group) MustGetPartial(name string) *Partial {
	partial, found := g.partials[name]
	if !found {
		panic("Could not find partial named " + name)
	}
	return partial
}

// MustGetLayout works like GetLayout, except that it panics
// instead of returning an error if the layout could not be
// found.
func (g Group) MustGetLayout(name string) *Layout {
	layout, found := g.layouts[name]
	if !found {
		panic("Could not find layout named " + name)
	}
	return layout
}

// NewGroup creates, initializes, and returns a new Group
func NewGroup() *Group {
	return &Group{
		templates: map[string]*Template{},
		partials:  map[string]*Partial{},
		layouts:   map[string]*Layout{},
		Funcs:     template.FuncMap{},
	}
}

// Executor represents some type of template that is capable of executing (i.e. rendering)
// to an io.Writer with some data. It is satisfied by Template, Partial, and Layout as well
// as the builtin template.Template.
type Executor interface {
	Execute(wr io.Writer, data interface{}) error
}

// AddFunc adds f to the FuncMap for the group under the given
// name. You must call AddFunc before adding any templates, partials,
// or layouts. Once added, all templates can call the function directly.
// See http://golang.org/pkg/text/template/ for more information
// about the FuncMap type and how to call functions from inside
// templates.
func (g *Group) AddFunc(name string, f interface{}) {
	g.Funcs[name] = f
}

// PrefixedName returns the name of the partial with PartialsPrefix
// included. By default, PartialsPrefix is "partials/". Each partial
// can be rendered inside any template, layout, or other partial using
// the `template` action with the prefixed name as the first argument.
func (p Partial) PrefixedName() string {
	if strings.HasPrefix(p.Name(), PartialPrefix) {
		return p.Name()
	} else {
		return PartialPrefix + p.Name()
	}
}

// PrefixedName returns the name of the layout with LayoutsPrefix
// included. By default, LayoutsPrefix is "layouts/". Each layout
// can be rendered by any template using the `template` action with
// the prefixed name as the first argument.
func (l Layout) PrefixedName() string {
	if strings.HasPrefix(l.Name(), LayoutPrefix) {
		return l.Name()
	} else {
		return LayoutPrefix + l.Name()
	}
}

// AddTemplate adds a regular template to the group with the
// given name and source.
func (g *Group) AddTemplate(name, src string) error {
	tmpl, err := template.New(name).Funcs(g.Funcs).Parse(src)
	if err != nil {
		return err
	}
	template := Template{
		Template: tmpl,
	}
	g.templates[tmpl.Name()] = &template
	return g.associateTemplate(template)
}

// AddTemplateFile reads the contents of filename and adds
// a template to the group using the given name and the contents
// of the file as the source.
func (g *Group) AddTemplateFile(name, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return g.AddTemplate(name, string(src))
}

// AddTemplateFiles recursively adds all the .tmpl files in dir
// and its subdirectories to the group. The name assigned to each
// template is based on the filename and the path relative to dir.
// So if dir is `templates`, the file located at
// `templates/people/show.tmpl` will be given the name "people/show".
func (g *Group) AddTemplateFiles(dir string) error {
	return collectTemplateFiles(dir, g.AddTemplateFile)
}

// associateTemplate adds all the needed associations to the template.
// Namely, it associates all partials and layouts with the template.
func (g *Group) associateTemplate(template Template) error {
	// Associate each partial with this template
	for _, partial := range g.partials {
		if template.Lookup(partial.PrefixedName()) == nil {
			if _, err := template.AddParseTree(partial.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	// Associate each layout with this template
	for _, layout := range g.layouts {
		if template.Lookup(layout.PrefixedName()) == nil {
			if _, err := template.AddParseTree(layout.PrefixedName(), layout.Tree); err != nil {
				return err
			}
		}
	}
	return nil
}

// AddPartial adds a partial to the group with the given name
// and source.
func (g *Group) AddPartial(name, src string) error {
	tmpl, err := template.New(name).Funcs(g.Funcs).Parse(src)
	if err != nil {
		return err
	}
	partial := Partial{
		Template: tmpl,
	}
	g.partials[tmpl.Name()] = &partial
	return g.associatePartial(partial)
}

// AddPartialFile reads the contents of filename and adds
// a partial to the group using the given name and the contents
// of the file as the source.
func (g *Group) AddPartialFile(name, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return g.AddPartial(name, string(src))
}

// AddPartialFiles recursively adds all the .tmpl files in dir
// and its subdirectories to the group. The name assigned to each
// partial is based on the filename and the path relative to dir.
// PartialsPrefix is prepended to the name, which by default is
// "partials/". So if dir is `my-partials`, the file located at
// `my-partials/head.tmpl` will be given the name "partials/head".
func (g *Group) AddPartialFiles(dir string) error {
	return collectTemplateFiles(dir, g.AddPartialFile)
}

// associatePartial adds all the needed associations to the partial.
// Namely, it associates it with all templates, layouts, and other
// partials.
func (g *Group) associatePartial(partial Partial) error {
	// Associate this partial with every template
	for _, template := range g.templates {
		if template.Lookup(partial.PrefixedName()) == nil {
			if _, err := template.AddParseTree(partial.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	// Associate this partial with every other partial
	for _, other := range g.partials {
		if other.Lookup(partial.PrefixedName()) == nil {
			if _, err := other.AddParseTree(partial.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	// Associate every other partial with this partial
	for _, other := range g.partials {
		if partial.Lookup(partial.PrefixedName()) == nil {
			if _, err := partial.AddParseTree(partial.PrefixedName(), other.Tree); err != nil {
				return err
			}
		}
	}
	// Associate this partial with every layout
	for _, layout := range g.layouts {
		if layout.Lookup(partial.PrefixedName()) == nil {
			if _, err := layout.AddParseTree(partial.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	return nil
}

// AddLayout adds a layout to the group with the given name
// and source.
func (g *Group) AddLayout(name, src string) error {
	tmpl, err := template.New(name).Funcs(g.Funcs).Parse(src)
	if err != nil {
		return err
	}
	layout := Layout{
		Template: tmpl,
	}
	g.layouts[tmpl.Name()] = &layout
	return g.associateLayout(layout)
}

// AddLayoutFile reads the contents of filename and adds
// a layout to the group using the given name and the contents
// of the file as the source.
func (g *Group) AddLayoutFile(name, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return g.AddLayout(name, string(src))
}

// AddLayoutFiles recursively adds all the .tmpl files in dir
// and its subdirectories to the group. The name assigned to each
// layout is based on the filename and the path relative to dir.
// LayoutsPrefix is prepended to the name, which by default is
// "layouts/". So if dir is `my-layouts`, the file located at
// `my-layouts/app.tmpl` will be given the name "layouts/app".
func (g *Group) AddLayoutFiles(dir string) error {
	return collectTemplateFiles(dir, g.AddLayoutFile)
}

// associateLayout adds all the needed associations to the layout.
// Namely, it associates the layout with all templates and associates
// all partials with the layout.
func (g *Group) associateLayout(layout Layout) error {
	// Associate this layout with every template
	for _, template := range g.templates {
		if template.Lookup(layout.PrefixedName()) == nil {
			if _, err := template.AddParseTree(layout.PrefixedName(), layout.Tree); err != nil {
				return err
			}
		}
	}
	// Associate each partial with this layout
	for _, partial := range g.partials {
		if layout.Lookup(layout.PrefixedName()) == nil {
			if _, err := layout.AddParseTree(layout.PrefixedName(), partial.Tree); err != nil {
				return err
			}
		}
	}
	return nil
}

// AddAllFiles adds the .tmpl files located in templatesDir, partialsDir,
// and layoutsDir to the group as regular templates, partials, and layouts,
// respectively. It also adds the needed associations. The name assigned to
// each template, partial, or layout is based on the filename and the path
// relative to dir, just as it is in AddTemplateFiles, AddPartialFiles, and
// AddLayoutFiles, respectively.
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

// collectTemplateFiles is a function which navigates recursively through
// dir and its subdirectories and finds any files with the .tmpl file extension.
// Then it calls the given handler func with the filename and template name (without
// the PartialsPrefix or LayoutsPrefix added).
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
