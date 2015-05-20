// Copyright 2015 Alex Browne.
// All rights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package temple

import (
	"bytes"
	"honnef.co/go/js/dom"
)

func ExecuteEl(e Executor, el dom.Element, data interface{}) error {
	// TODO: use a buffer pool
	buf := bytes.NewBuffer([]byte{})
	if err := e.Execute(buf, data); err != nil {
		return err
	}
	el.SetInnerHTML(buf.String())
	return nil
}

func (t *Template) ExecuteEl(el dom.Element, data interface{}) error {
	return ExecuteEl(t, el, data)
}

func (p *Partial) ExecuteEl(el dom.Element, data interface{}) error {
	return ExecuteEl(p, el, data)
}

func (l *Layout) ExecuteEl(el dom.Element, data interface{}) error {
	return ExecuteEl(l, el, data)
}

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

func (g *Group) AddInlineTemplate(el dom.Element) error {
	return g.AddTemplate(el.ID(), el.InnerHTML())
}

func (g *Group) AddInlinePartial(el dom.Element) error {
	return g.AddPartial(el.ID(), el.InnerHTML())
}

func (g *Group) AddInlineLayout(el dom.Element) error {
	return g.AddLayout(el.ID(), el.InnerHTML())
}
