package temple

import (
	"bytes"
	"fmt"
	"html/template"
)

var (
	Funcs = map[string]interface{}{
		"partial": PartialHelper,
		"layout":  LayoutHelper,
	}
)

func PartialHelper(name string, context interface{}) (template.HTML, error) {
	partial, found := Partials[name]
	if !found {
		return "", fmt.Errorf("temple: Error in func partial: Could not find partial with name: %s", name)
	}
	// TODO: use a buffer pool
	buf := bytes.NewBuffer([]byte{})
	if err := partial.Execute(buf, context); err != nil {
		return "", fmt.Errorf("temple: Error in func partial: %s", err.Error())
	}
	return template.HTML(buf.String()), nil
}

func LayoutHelper(name string, context interface{}) (template.HTML, error) {
	layout, found := Layouts[name]
	if !found {
		return "", fmt.Errorf("temple: Error in func layout: Could not find layout with name: %s", name)
	}
	// TODO: use a buffer pool
	buf := bytes.NewBuffer([]byte{})
	if err := layout.Execute(buf, context); err != nil {
		return "", fmt.Errorf("temple: Error in func layout: %s", err.Error())
	}
	return template.HTML(buf.String()), nil
}
