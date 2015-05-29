// Copyright 2015 Alex Browne.
// All rights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package temple

import (
	"testing"
)

func TestAddTemplate(t *testing.T) {
	g := NewGroup()
	if err := g.AddTemplate("test", `Hello, {{ . }}!`); err != nil {
		t.Fatalf("Unexpected error in AddTemplate: %s", err.Error())
	}
	// Get the template from the map
	testTmpl, found := g.templates["test"]
	if !found {
		t.Fatal(`Template named "test" was not added to map of Templates`)
	}
	expectExecutorOutputs(t, testTmpl, "world", "Hello, world!")
}

func TestAddPartial(t *testing.T) {
	g := NewGroup()
	// Add some Partials and make sure each is added to the map
	partials := map[string]string{
		"foo": "foo",
		"bar": "bar",
		"baz": "baz",
		// foobarbaz calls on each of the other partials. This tests that
		// partials are associated with all other partials.
		"foobarbaz": `{{ template "partials/foo" }}{{ template "partials/bar" }}{{ template "partials/baz" }}`,
	}
	for name, src := range partials {
		if err := g.AddPartial(name, src); err != nil {
			t.Fatalf("Unexpected error in AddPartial: %s", err.Error())
		}
		if _, found := g.partials[name]; !found {
			t.Errorf(`Partial named "%s" was not added to the map of Partials`, name)
		}
	}
	// The test template calls on each of the four partials. This tests that
	// partials are associated with templates.
	if err := g.AddTemplate("test", `{{ template "partials/foo" }} {{ template "partials/bar" }} {{ template "partials/baz" }} {{ template "partials/foobarbaz" }}`); err != nil {
		t.Fatalf("Unexpected error in AddTemplate: %s", err.Error())
	}
	testTmpl, found := g.templates["test"]
	if !found {
		t.Fatal(`Template named "test" was not added to map of Templates`)
	}
	expectExecutorOutputs(t, testTmpl, nil, "foo bar baz foobarbaz")
}

func TestAddLayout(t *testing.T) {
	g := NewGroup()
	// The foo partial will be called on by the header layout, which tests that
	// partials are associated with layouts.
	if err := g.AddPartial("foo", "foo"); err != nil {
		t.Fatalf("Unexpected error in AddPartial: %s", err.Error())
	}
	if _, found := g.partials["foo"]; !found {
		t.Errorf(`Partial named "%s" was not added to the map of Partials`, "foo")
	}
	// The header layout renders a content template (which must be defined by a template using
	// the layout) and calls for the foo partial.
	if err := g.AddLayout("header", `<h2>{{ template "content" }} {{ template "partials/foo" }}</h2>`); err != nil {
		t.Fatalf("Unexpected error in AddLayout: %s", err.Error())
	}
	if _, found := g.layouts["header"]; !found {
		t.Errorf(`Layout named "%s" was not added to the map of Layouts`, "header")
	}
	// The test template defines a content template and attempts to render itself inside the
	// header layout.
	if err := g.AddTemplate("test", `{{ define "content"}}test{{end}}{{ template "layouts/header" }}`); err != nil {
		t.Fatalf("Unexpected error in AddTemplate: %s", err.Error())
	}
	testTmpl, found := g.templates["test"]
	if !found {
		t.Fatal(`Template named "test" was not added to map of Templates`)
	}
	expectExecutorOutputs(t, testTmpl, nil, "<h2>test foo</h2>")
}

func TestAddAllFiles(t *testing.T) {
	g := NewGroup()
	// Load all the files from the test_files directory
	if err := g.AddAllFiles("test_files/templates", "test_files/partials", "test_files/layouts"); err != nil {
		t.Fatalf("Unexpected error in AddAllFiles: %s", err.Error())
	}
	// Try rendering the todos/index template with some data
	type Todo struct {
		Title string
	}
	todos := []Todo{
		{Title: "One"},
		{Title: "Two"},
		{Title: "Three"},
	}
	todosTmpl, found := g.templates["todos/index"]
	if !found {
		t.Fatal(`Template named "todos/index" was not added to map of Templates`)
	}
	expectExecutorOutputs(t, todosTmpl, todos, "<html><head><title>Todos</title></head><body><ul><li>One</li><li>Two</li><li>Three</li></ul></body></html>")
}
