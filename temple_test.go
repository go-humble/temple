package temple

import (
	"bytes"
	"html/template"
	"testing"
)

func reset() {
	Templates = map[string]Template{}
	Partials = map[string]Partial{}
	Layouts = map[string]Layout{}
}

func TestAddTemplate(t *testing.T) {
	defer reset()
	if err := AddTemplate(template.Must(template.New("test").Parse(`Hello, {{ . }}!`))); err != nil {
		t.Fatalf("Unexpected error in AddTemplate: %s", err.Error())
	}
	// Get the template from the map
	testTmpl, found := Templates["test"]
	if !found {
		t.Fatal(`Template named "test" was not added to map of Templates`)
	}
	// Execute the template and check the result
	buf := bytes.NewBuffer([]byte{})
	if err := testTmpl.Execute(buf, "world"); err != nil {
		t.Errorf("Unexpected error executing template: %s", err.Error())
	}
	expected := "Hello, world!"
	got := buf.String()
	if expected != got {
		t.Errorf("Template was not executed correctly. Expected `%s` but got `%s`.", expected, got)
	}
}

func TestAddPartial(t *testing.T) {
	defer reset()
	// Add some Partials and make sure each is added to the map
	partials := []*template.Template{
		template.Must(template.New("foo").Parse("foo")),
		template.Must(template.New("bar").Parse("bar")),
		template.Must(template.New("baz").Parse("baz")),
		// foobarbaz calls on each of the other partials. This tests that
		// partials are associated with all other partials.
		template.Must(template.New("foobarbaz").Parse(`{{ template "partials/foo" }}{{ template "partials/bar" }}{{ template "partials/baz" }}`)),
	}
	for _, partial := range partials {
		if err := AddPartial(partial); err != nil {
			t.Fatalf("Unexpected error in AddPartial: %s", err.Error())
		}
		if _, found := Partials[partial.Name()]; !found {
			t.Errorf(`Partial named "%s" was not added to the map of Templates`, partial.Name())
		}
	}
	// The test template calls on each of the four partials. This tests that
	// partials are associated with templates.
	if err := AddTemplate(template.Must(template.New("test").Parse(`{{ template "partials/foo" }} {{ template "partials/bar" }} {{ template "partials/baz" }} {{ template "partials/foobarbaz" }}`))); err != nil {
		t.Fatalf("Unexpected error in AddTemplate: %s", err.Error())
	}
	testTmpl, found := Templates["test"]
	if !found {
		t.Fatal(`Template named "test" was not added to map of Partials`)
	}
	// Execute the template and check the result
	buf := bytes.NewBuffer([]byte{})
	if err := testTmpl.Execute(buf, nil); err != nil {
		t.Errorf("Unexpected error executing template: %s", err.Error())
	}
	expected := "foo bar baz foobarbaz"
	got := buf.String()
	if expected != got {
		t.Errorf("Template was not executed correctly. Expected `%s` but got `%s`.", expected, got)
	}
}

func TestAddLayout(t *testing.T) {
	defer reset()
	// The foo partial will be called on by the header layout, which tests that
	// partials are associated with layouts.
	if err := AddPartial(template.Must(template.New("foo").Parse("foo"))); err != nil {
		t.Fatalf("Unexpected error in AddPartial: %s", err.Error())
	}
	if _, found := Partials["foo"]; !found {
		t.Errorf(`Partial named "%s" was not added to the map of Partials`, "foo")
	}
	// The header layout renders a content template (which must be defined by a template using
	// the layout) and calls for the foo partial.
	if err := AddLayout(template.Must(template.New("header").Parse(`<h2>{{ template "content" }} {{ template "partials/foo" }}</h2>`))); err != nil {
		t.Fatalf("Unexpected error in AddLayout: %s", err.Error())
	}
	if _, found := Layouts["header"]; !found {
		t.Errorf(`Layout named "%s" was not added to the map of Layouts`, "header")
	}
	// The test template defines a content template and attempts to render itself inside the
	// header layout.
	if err := AddTemplate(template.Must(template.New("test").Parse(`{{ define "content"}}test{{end}}{{ template "layouts/header" }}`))); err != nil {
		t.Fatalf("Unexpected error in AddTemplate: %s", err.Error())
	}
	testTmpl, found := Templates["test"]
	if !found {
		t.Fatal(`Template named "test" was not added to map of Templates`)
	}
	// Execute the template and check the result
	buf := bytes.NewBuffer([]byte{})
	if err := testTmpl.Execute(buf, nil); err != nil {
		t.Errorf("Unexpected error executing template: %s", err.Error())
	}
	expected := "<h2>test foo</h2>"
	got := buf.String()
	if expected != got {
		t.Errorf("Template was not executed correctly. Expected `%s` but got `%s`.", expected, got)
	}
}
