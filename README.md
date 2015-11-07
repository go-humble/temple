Humble/Temple
=============

[![Version](https://img.shields.io/badge/version-X.X.X%20develop-5272B4.svg)](https://github.com/go-humble/temple/releases)
[![GoDoc](https://godoc.org/github.com/go-humble/temple?status.svg)](https://godoc.org/github.com/go-humble/temple)


A library and a command line tool for sanely managing go templates, with the ability
to share them between client and server. The library and code generated by the cli are
both compatible with [gopherjs](https://github.com/gopherjs/gopherjs), so you can compile
to javascript and run in the browser. Temple works great as a stand-alone package or in
combination with other packages in the [Humble Framework](https://github.com/go-humble/humble).

This README is specific to the command line tool, which reads .tmpl files in your project and
generates go source code. The command line tool uses the library (mainly the Build function).
The README for the library can be found at
[github.com/go-humble/temple/temple](https://github.com/go-humble/temple/tree/develop/temple).


What Does the Command Line Tool Do?
-----------------------------------

The command line tool first reads the contents of all the .tmpl files in the specified src directory
(and it's sub-directories). Then it generates go code at the specified dest file which compiles the
templates. The generated code uses strings which contain the content of source files as arguments to
[`template.Parse`](http://golang.org/pkg/html/template/#Template.Parse).

This is useful for a few reasons. If you're building and deploying binary executables, it means you
no longer have to ship a templates folder with your binary distribution. It also makes it possible
for you to use templates in code that has been compiled to javascript with gopherjs and is running
in a browser. Since code generated by temple works exactly the same on the server and client, it also
enables you to share templates between them. Check out
[go-humble/examples/people](https://github.com/go-humble/examples/tree/master/people) for an example of 
how shared templates can work.


Browser Support
---------------

Temple is regularly tested with IE9+ and the latest versions of Chrome, Safari, and Firefox.

Javascript code generated with gopherjs uses typed arrays, so in order to work with IE9,
you will need a
[polyfill for typed arrays](https://github.com/inexorabletash/polyfill/blob/master/typedarray.js). 


Installation
------------

Install the temple command line tool with the following.

```bash
go get -u github.com/go-humble/temple
```

You may also need to install gopherjs. The latest version is recommended. Install
gopherjs with:

```bash
go get -u github.com/gopherjs/gopherjs
```


Usage Guide
-----------

### Basic Usage

The temple command line tool has one main subcommand called `build`. Typical usage will look
something like this:

`temple build templates templates/templates.go`

The first argument, in this case `templates`, is a directory which contains your .tmpl files.
The second argument, in this case `templates/templates.go` is the name of a file where generated
code will be written. If the file does not exist, temple will create it for you, and any previous
content will be overwritten.

You can run `temple help` to learn more about the possible commands and `temple help build` to
learn more about the build command specifically.

### Generated Code

The code generated by the temple command line tool will look something like this:

```go
package templates

// This package has been automatically generated with temple.
// Do not edit manually!

import (
	"github.com/go-humble/temple/temple"
)

var (
	GetTemplate     func(name string) (*temple.Template, error)
	GetPartial      func(name string) (*temple.Partial, error)
	GetLayout       func(name string) (*temple.Layout, error)
	MustGetTemplate func(name string) *temple.Template
	MustGetPartial  func(name string) *temple.Partial
	MustGetLayout   func(name string) *temple.Layout
)

func init() {
	var err error
	g := temple.NewGroup()

	if err = g.AddPartial("head", `...`); err != nil {
		panic(err)
	}

	if err = g.AddLayout("app", `...`); err != nil {
		panic(err)
	}

	if err = g.AddTemplate("people/index", `...`); err != nil {
		panic(err)
	}

	GetTemplate = g.GetTemplate
	GetPartial = g.GetPartial
	GetLayout = g.GetLayout
	MustGetTemplate = g.MustGetTemplate
	MustGetPartial = g.MustGetPartial
	MustGetLayout = g.MustGetLayout
}
```

The code creates a single template [`Group`](http://godoc.org/github.com/go-humble/temple/temple/#Group)
and adds the templates, partials, and layouts (if applicable) to the group. Then it exposes the methods of
the group for getting templates, partials, and layouts as exported global functions.

The [temple.Template](http://godoc.org/github.com/go-humble/temple/temple/#Template),
[temple.Partial](http://godoc.org/github.com/go-humble/temple/temple/#Partial), and
[temple.Layout](http://godoc.org/github.com/go-humble/temple/temple/#Layout) types all inherit from
the builtin
[Template type from the html/template package](http://golang.org/pkg/html/template/#Template). That
means you can render them like regular templates with the
[`Execute`](http://golang.org/pkg/html/template/#Template.Execute) method. Temple also provides an
additional method for rendering templates in the dom called
[`ExecuteEl`](http://godoc.org/github.com/go-humble/temple/temple/#ExecuteEl).

### Naming conventions

In go, every template needs to have a name. temple assigns a name to each template based on its
filename and location relative to the src directory. So for example, if you had a template file
located at `templates/people/show.tmpl` and `templates` was your src directory, the name assigned
to the template would be `"people/show"`.

### Partials and Layouts

Temple uses two optional groups called "partials" and "layouts" to help organize template files.
You can specify a directory that contains partials with the `--partials` flag, and a directory
that contains layouts with the `--layouts` flag. Any .tmpl files found in these directories will
be treated specially, and they should not overlap with each other or the src directory with all
your regular templates ("regular templates" is the name we'll use to refer to .tmpl files in the
src directory that are neither partials or layouts). This organization feature is completely optional,
so if you don't want to use it, you omit the `--partials` and `--layouts` flags and manage your
templates any way you want.

Before continuing, it is recommended that you read the documentation for the
[text/template](http://golang.org/pkg/text/template/) and
[html/template](http://golang.org/pkg/html/template/) packages. In addition, this
article about
[template inheritence in go](https://elithrar.github.io/article/approximating-html-template-inheritance/)
will help explain some of the concepts that temple uses.

#### Partials

Partials are templates which typically represent only part of a full page. For example, you might
have a partial for rendering a single model or a partial for the head section of your html. Partials
are associated with (i.e., added to the parse tree of) all other partials, in addition to layouts and
regular templates. That means you can render a partial inside of a regular template or layout with the
`template` action. The name of the partial templates is based on its filename and location relative to the
partials directory. [`PartialsPrefix`](http://godoc.org/github.com/albrow/prtty#pkg-variables) is added
to the template name of all partials, which by default is simply `"partials/"`.

So if your partials directory is `my-partials`, and you have the following partial template file located
at `my-partials/head.tmpl`:

```handlebars
<head>
	<title>Example Humble Application</title>	
	<link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/css/bootstrap.min.css" rel="stylesheet">	
</head>
```

You could render it inside of another template like so:

```handlebars
<!doctype html>
<html>
	{{ template "partials/head" }}
	<body>
	</body>
</html>
```

Check out the partials in
[go-humble/examples/people](https://github.com/go-humble/examples/tree/master/people/shared/templates/partials)
for a more in-depth example.

#### Layouts

Layouts are templates which define the structure for a page and require another template to fill in
the details. Typically, layouts will use one or more `template` actions to render partials or regular
templates inside of some structure. Layouts are associated with (i.e., added to the parse tree of) all
regular templates and have access to partials. That means you can render layouts inside of a regular template,
typically after declaring some sub-template that the layout expects to be defined. The name of the layout
templates is based on its filename and location relative to the layouts directory.
[`LayoutsPrefix`](http://godoc.org/github.com/go-humble/temple#pkg-variables) is added to the template
name of all layouts, which by default is simply `"layouts/"`.

For example, if your layouts directory is `my-layouts`, and you have the following layout template file
located at `my-layouts/app.tmpl`:

```handlebars
<!doctype html>
<html>
	<head>
		<title>Example Humble Application</title>
	</head>
	<body>
		{{ template "content" }}
	</body>
</html>
```

You could then render a template inside of the layout by first defining the "content" sub-template and
then rendering the layout:

```handlebars
{{ define "content" }}
	Hello, Content!
{{ end }}
{{ template "layouts/app" }}
```

If you rendered the template (not the layout), the output would look like this:

```html
<!doctype html>
<html>
	<head>
		<title>Example Humble Application</title>
	</head>
	<body>
		Hello, Content!
	</body>
</html>
```

Check out the layouts in
[go-humble/examples/people](https://github.com/go-humble/examples/tree/master/people/shared/templates/layouts)
for a more in-depth example.

Testing
-------

Temple uses regular go testing, so you can run the all the tests with `go test ./...`.


Contributing
------------

See [CONTRIBUTING.md](https://github.com/go-humble/temple/blob/master/CONTRIBUTING.md)


License
-------

Temple is licensed under the MIT License. See the [LICENSE](https://github.com/go-humble/temple/blob/master/LICENSE)
file for more information.
