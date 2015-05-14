package temple

import (
	"bytes"
	"errors"
	"github.com/albrow/prtty"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

func Build(src, dest, partials, layouts string) error {
	prtty.Info.Println("--> building...")
	prtty.Default.Printf("    src: %s", src)
	prtty.Default.Printf("    dest: %s", dest)
	if partials != "" {
		prtty.Default.Printf("    partials: %s", partials)
	}
	if layouts != "" {
		prtty.Default.Printf("    layouts: %s", layouts)
	}
	dirs := sourceDirGroup{
		templates: src,
		partials:  partials,
		layouts:   layouts,
	}
	if err := checkCompileTemplates(dirs); err != nil {
		return err
	}
	if err := generateFile(dirs, dest); err != nil {
		return err
	}
	prtty.Info.Println("--> done!")
	return nil
}

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

type templateData struct {
	PackageName string
	Templates   []sourceFile
	Partials    []sourceFile
	Layouts     []sourceFile
}

type sourceFile struct {
	Name string
	Src  string
}

type sourceDirGroup struct {
	templates string
	partials  string
	layouts   string
}

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

func generateFile(dirs sourceDirGroup, dest string) error {
	prtty.Info.Println("--> generating go code...")
	packageName := filepath.Base(filepath.Dir(dest))
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

//go:generate go-bindata --pkg=temple templates/...

func (data *templateData) writeToFile(dest string) error {
	tmplAsset, err := Asset("templates/generated.go.tmpl")
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
