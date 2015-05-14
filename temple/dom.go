package temple

import (
	"bytes"
	"honnef.co/go/js/dom"
)

func ExecuteToEl(e Executor, el dom.Element, data interface{}) error {
	// TODO: use a buffer pool
	buf := bytes.NewBuffer([]byte{})
	if err := e.Execute(buf, data); err != nil {
		return err
	}
	el.SetInnerHTML(buf.String())
	return nil
}

func (t *Template) ExecuteToEl(el dom.Element, data interface{}) error {
	return ExecuteToEl(t, el, data)
}

func (p *Partial) ExecuteToEl(el dom.Element, data interface{}) error {
	return ExecuteToEl(p, el, data)
}

func (l *Layout) ExecuteToEl(el dom.Element, data interface{}) error {
	return ExecuteToEl(l, el, data)
}

func ParseInlineTemplates() error {
	document := dom.GetWindow().Document()
	elements := document.QuerySelectorAll(`script[type="text/template"]`)
	for _, el := range elements {
		switch el.GetAttribute("data-kind") {
		case "template":
			if err := AddInlineTemplate(el); err != nil {
				return err
			}
		case "partial":
			if err := AddInlinePartial(el); err != nil {
				return err
			}
		case "layout":
			if err := AddInlineLayout(el); err != nil {
				return err
			}
		default:
			if err := AddInlineTemplate(el); err != nil {
				return err
			}
		}
	}
	return nil
}

func AddInlineTemplate(el dom.Element) error {
	return AddTemplate(el.ID(), el.InnerHTML())
}

func AddInlinePartial(el dom.Element) error {
	return AddPartial(el.ID(), el.InnerHTML())
}

func AddInlineLayout(el dom.Element) error {
	return AddLayout(el.ID(), el.InnerHTML())
}
