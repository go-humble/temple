// Copyright 2015 Alex Browne.
// All rights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package temple

import (
	"bytes"
	"honnef.co/go/js/dom"
)

// ExecuteEl executes an Executor with the given data and then
// writes the result to  the innerHTML of el. It only works if
// you have compiled this code to javascript with gopherjs and
// it is running in a browser.
func ExecuteEl(e Executor, el dom.Element, data interface{}) error {
	// TODO: use a buffer pool
	buf := bytes.NewBuffer([]byte{})
	if err := e.Execute(buf, data); err != nil {
		return err
	}
	el.SetInnerHTML(buf.String())
	return nil
}

// ExecuteEl executes the template with the given data and then
// writes the result to  the innerHTML of el. It only works if
// you have compiled this code to javascript with gopherjs and
// it is running in a browser.
func (t *Template) ExecuteEl(el dom.Element, data interface{}) error {
	return ExecuteEl(t, el, data)
}

// ExecuteEl executes the partial with the given data and then
// writes the result to  the innerHTML of el. It only works if
// you have compiled this code to javascript with gopherjs and
// it is running in a browser.
func (p *Partial) ExecuteEl(el dom.Element, data interface{}) error {
	return ExecuteEl(p, el, data)
}

// ExecuteEl executes the layout with the given data and then
// writes the result to  the innerHTML of el. It only works if
// you have compiled this code to javascript with gopherjs and
// it is running in a browser.
func (l *Layout) ExecuteEl(el dom.Element, data interface{}) error {
	return ExecuteEl(l, el, data)
}

// ParseInlineTemplates scans the DOM for inline templates which
// must be script tags with the type "text/template". The id property
// will be used for the name of each template, and the special
// property "data-kind" can be used to distinguish between regular
// templates, partials, and layouts. So, to declare an inline partial
// for use with the ParseInlineTemplates method, use an opening script
// tag that looks like:
//   <script type="text/template" id="todo" data-kind="partial">
func (g *Group) ParseInlineTemplates() error {
	document := dom.GetWindow().Document()
	elements := document.QuerySelectorAll(`script[type="text/template"]`)
	for _, el := range elements {
		switch el.GetAttribute("data-kind") {
		case "template":
			if err := g.AddInlineTemplate(el); err != nil {
				return err
			}
		case "partial":
			if err := g.AddInlinePartial(el); err != nil {
				return err
			}
		case "layout":
			if err := g.AddInlineLayout(el); err != nil {
				return err
			}
		default:
			if err := g.AddInlineTemplate(el); err != nil {
				return err
			}
		}
	}
	return nil
}

// AddInlineTemplate adds the inline template el to the
// group as a regular template. It uses the id property as
// the template name and the innerHTML as the template source.
// Typically inline templates will be in script tags that look
// like:
//   <script type="text/template" id="home">
func (g *Group) AddInlineTemplate(el dom.Element) error {
	return g.AddTemplate(el.ID(), el.InnerHTML())
}

// AddInlinePartial adds the inline template el to the
// group as a partial. It uses the id property as the template
// name and the innerHTML as the template source. PartialsPrefix
// will be added to the name, which by default is "partials/".
// Typically
// inline templates will be in script tags that look like:
//   <script type="text/template" id="todo">
// which would correspond to a template with the name "partials/todo".
func (g *Group) AddInlinePartial(el dom.Element) error {
	return g.AddPartial(el.ID(), el.InnerHTML())
}

// AddInlineLayout adds the inline template el to the
// group as a layout. It uses the id property as the template
// name and the innerHTML as the template source. LayoutsPrefix
// will be added to the name, which by default is "layouts/".
// Typically inline templates will be in script tags that look like:
//   <script type="text/template" id="app">
// which would correspond to a template with the name "layouts/app".
func (g *Group) AddInlineLayout(el dom.Element) error {
	return g.AddLayout(el.ID(), el.InnerHTML())
}
