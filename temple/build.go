// Copyright 2015 Alex Browne.
// All rights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package temple

import (
	"bytes"
	"errors"
	"github.com/albrow/prtty"
	"github.com/go-humble/temple/temple/assets"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

// Build is the function called when you run the build sub-command
// in the command line tool. It compiles all the templates in the
// src directory and generates go code in the dest file. If partials
// and/or layouts are provided, it will add them to the generated file
// with calls to AddPartial and AddLayout. If packageName is an empty
// string, the package name will be the directory of the dest file.
func Build(src, dest, partials, layouts, packageName string) error {
	prtty.Info.Println("--> building...")
	prtty.Default.Printf("    src: %s", src)
	prtty.Default.Printf("    dest: %s", dest)
	if partials != "" {
		prtty.Default.Printf("    partials: %s", partials)
	}
	if layouts != "" {
		prtty.Default.Printf("    layouts: %s", layouts)
	}
	if packageName != "" {
		prtty.Default.Printf("    package: %s", packageName)
	}
	dirs := sourceDirGroup{
		templates: src,
		partials:  partials,
		layouts:   layouts,
	}
	if err := checkCompileTemplates(dirs); err != nil {
		return err
	}
	if err := generateFile(dirs, dest, packageName); err != nil {
		return err
	}
	prtty.Info.Println("--> done!")
	return nil
}

// checkCompileTemplates compiles the templates, partials, and layouts
// in dirs with the correct associations to make sure that the templates
// compile. If they don't, we can catch errors early and return them when
// the command line tool is invoked, instead of at runtime.
func checkCompileTemplates(dirs sourceDirGroup) error {
	prtty.Info.Println("--> checking for compilation errors...")
	if dirs.templates == "" {
		return errors.New("temple: templates dir cannot be an empty string.")
	}
	g := NewGroup()
	if dirs.partials != "" {
		prtty.Default.Println("    checking partials...")
		if err := g.AddPartialFiles(dirs.partials); err != nil {
			return err
		}
	}
	if dirs.layouts != "" {
		prtty.Default.Println("    checking layouts...")
		if err := g.AddLayoutFiles(dirs.layouts); err != nil {
			return err
		}
	}
	prtty.Default.Println("    checking templates...")
	if err := g.AddTemplateFiles(dirs.templates); err != nil {
		return err
	}
	return nil
}

// templateData is passed in to the template for the generated code.
type templateData struct {
	PackageName string
	Templates   []sourceFile
	Partials    []sourceFile
	Layouts     []sourceFile
}

// sourceFile represents the source file for a template, partial, or layout.
type sourceFile struct {
	Name string
	Src  string
}

// sourceDirGroup represents a group of source directories, consisting of a
// directory for regular layouts and optionally for partials and layouts.
// The directories for partials and layouts will be empty strings if they
// were not provided.
type sourceDirGroup struct {
	templates string
	partials  string
	layouts   string
}

// collectAllSourceFiles walks recursively through the directories in dirs
// and collects all template, partial, and layout source files, adding them
// to data.
func (data *templateData) collectAllSourceFiles(dirs sourceDirGroup) error {
	if dirs.partials != "" {
		prtty.Info.Println("--> collecting partials...")
		partials, err := collectSourceFiles(dirs.partials)
		if err != nil {
			return err
		}
		data.Partials = partials
	}
	if dirs.layouts != "" {
		prtty.Info.Println("--> collecting layouts...")
		layouts, err := collectSourceFiles(dirs.layouts)
		if err != nil {
			return err
		}
		data.Layouts = layouts
	}
	prtty.Info.Println("--> collecting templates...")
	templates, err := collectSourceFiles(dirs.templates)
	if err != nil {
		return err
	}
	data.Templates = templates
	return nil
}

// generateFile generates go code containing the contents of all the
// files in the sourceDirGroup and writes the code to the dest file. It
// uses the given packageName if it is non-empty, and otherwise falls back
// to the directory that dest is in. If a file already exists at dest, it
// will be overwritten.
func generateFile(dirs sourceDirGroup, dest, packageName string) error {
	prtty.Info.Println("--> generating go code...")
	if packageName == "" {
		packageName = filepath.Base(filepath.Dir(dest))
	}
	data := &templateData{
		PackageName: packageName,
	}
	if err := data.collectAllSourceFiles(dirs); err != nil {
		return err
	}
	if err := data.writeToFile(dest); err != nil {
		return err
	}
	return nil
}

//go:generate go-bindata --pkg=assets -o=assets/bindata.go templates/...

// writeToFile writes the given templateData to the file located
// at dest. It uses the template located at templates/generated.go.tmpl.
// If there is already a file located at dest, it will be overwritten.
func (data *templateData) writeToFile(dest string) error {
	tmplAsset, err := assets.Asset("templates/generated.go.tmpl")
	if err != nil {
		return err
	}
	generatedTmpl := template.Must(template.New("generated").Parse(string(tmplAsset)))
	buf := bytes.NewBuffer([]byte{})
	if err := generatedTmpl.Execute(buf, data); err != nil {
		return err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	prtty.Success.Printf("    created %s", dest)
	return ioutil.WriteFile(dest, formatted, os.ModePerm)
}

// collectSourceFiles recursively walks through dir and its subdirectories
// and returns an array of all the source files (files which end in .tmpl).
func collectSourceFiles(dir string) ([]sourceFile, error) {
	sourceFiles := []sourceFile{}
	if err := collectTemplateFiles(dir, func(name, filename string) error {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		prtty.Default.Printf("    %s", filename)
		sourceFiles = append(sourceFiles, sourceFile{
			Name: name,
			Src:  string(src),
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return sourceFiles, nil
}
