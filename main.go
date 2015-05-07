package main

import (
	"fmt"
	"github.com/albrow/prtty"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	version = "temple version X.X.X (develop)"
)

var (
	// NOTE: GOPATH might consist of multiple paths. If that is the case, we look in the first path.
	gopath        = strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))[0]
	templePath    = filepath.Join(gopath, "src", "github.com", "albrow", "temple")
	generatedTmpl = template.Must(template.ParseFiles(filepath.Join(templePath, "templates.go.tmpl")))
)

func main() {
	cmdBuild := &cobra.Command{
		Use:   "build <src> <dest>",
		Short: "Compile the templates in the src directory and write generated go code to the dest file.",
		Long: `The build command will compile the .tmpl files found in the src directory,

along with the .tmpl files found in the layouts and includes directories (if
provided). It generates go source code containing the compiled templates and 
writes it to the dest file.

The build command works best if your templates are organized to approximate
template inheritance, as described in this article:
https://elithrar.github.io/article/approximating-html-template-inheritance/.
However, if you don't want organize your templates this way, the generated go
file will give you direct access to the builtin html templates
(*template.Template objects from the html/template package), so you can combine
parse trees manually. You also have the option of not combining parse trees at
all, and simply having each .tmpl file represent a stand-alone template.

The generated go file is designed to be fairly readable for humans, so feel
free to take a look. (Just don't edit it directly!)

## Includes

If provided, all .tmpl files in the includes directory are referred to as
"includes templates" or simply "includes". Includes are parsed first (before
layouts and regular templates). Includes should contain .tmpl files for things
like the <head> section, which are shared between different layouts, or other
components that are shared between different regular templates. No .tmpl file in
the includes directory can conflict with any other .tmpl file (e.g. they cannot
declare sub-templates of the same name). All the includes will be added to the
parse tree for the layouts and all other templates via the template.AddParseTree
method. It is safe for includes to reference each other, as long as they don't
conflict or create cyclical references.

## Layouts

If provided, all .tmpl files in the layouts directory are referred to as "layout
templates" or simply "layouts". Layouts are parsed after includes and before
regular templates. Typically, layouts will be referenced by a regular template,
and will expect the regular template to define certain sub-templates (e.g.
"content" or "title"), which will then be inserted into the layout. An
application will almost always want to have at least one layout, conventially
called "app.tmpl", which regular templates will use. No .tmpl file in the
layouts directory can conflict with any other .tmpl file (e.g. they cannot
declare sub-templates of the same name). If includes were also provided, all
includes will be added to the parse tree for each layout via the
template.AddParseTree method. Therefore a layout can reference any template in
includes. Layouts can also reference each other, as long as they don't conflict
or create cylclical references. All layouts will be added to the parse tree for
the regular templates in the src directory via the template.AddParseTree method.

## Regular Templates

All the .tmpl files found in the src directory are referred to as "regular
templates", or simply "templates", and are parsed last. All layouts and includes
(if any) are added to the parse tree for each template via the
template.AddParseTree method. Therefore templates can reference both layouts and
includes. Since regular templates will never be parsed together, they can
conflict with eachother (e.g. they can declare sub-templates of the same name).
As a consequence, regular templates also cannot reference eachother.

`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				prtty.Error.Fatal("temple build requires exactly 2 arguments: the src directory and the dest file.")
			}
			includes := cmd.Flag("includes").Value.String()
			layouts := cmd.Flag("layouts").Value.String()
			if err := build(args[0], args[1], includes, layouts); err != nil {
				prtty.Error.Fatal(err)
			}
		},
	}
	cmdBuild.Flags().String("includes", "", "(optional) The directory to look for includes. Includes are .tmpl files that are shared between layouts and all templates.")
	cmdBuild.Flags().String("layouts", "", "(optional) The directory to look for layouts. Layouts are .tmpl shared between all templates and have access to includes.")

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Print the current version number.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	rootCmd := &cobra.Command{
		Use:   "temple",
		Short: "A command line tool for sharing go templates between client and server.",
	}
	rootCmd.AddCommand(cmdBuild, cmdVersion)
	if err := rootCmd.Execute(); err != nil {
		prtty.Error.Fatal(err)
	}
}

type TemplateData struct {
	PackageName string
	Templates   []*TemplateFile
	Includes    []*TemplateFile
	Layouts     []*TemplateFile
}

type TemplateFile struct {
	VarName string
	Name    string
	Source  string
}

func build(src, dest, includes, layouts string) error {
	prtty.Info.Println("--> building...")
	prtty.Default.Printf("    src: %s", src)
	prtty.Default.Printf("    dest: %s", dest)
	if includes != "" {
		prtty.Default.Printf("    includes: %s", includes)
	}
	if layouts != "" {
		prtty.Default.Printf("    layouts: %s", layouts)
	}
	templateData, err := generateTemplateData(src, dest, includes, layouts)
	if err != nil {
		return err
	}
	if err := writeFile(templateData, dest); err != nil {
		return err
	}
	if err := formatFile(dest); err != nil {
		return err
	}
	prtty.Info.Println("--> done")
	return nil
}

func NewTemplateFile(filename string) (*TemplateFile, error) {
	// name is everything after the last slash, not including the file extension
	name := strings.TrimSuffix(filepath.Base(filename), ".tmpl")
	// varName is just the name titlized so it is an exported variable
	varName := strings.Title(name)
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &TemplateFile{
		VarName: varName,
		Name:    name,
		Source:  string(fileContents),
	}, nil
}

func ParseTemplateFiles(dir string) ([]*TemplateFile, error) {
	templateFiles := []*TemplateFile{}
	files, err := filepath.Glob(filepath.Join(dir, "*.tmpl"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		prtty.Warn.Printf("    WARNING: No .tmpl files found in %s", dir)
	}
	for _, filename := range files {
		prtty.Default.Printf("    %s", filename)
		tf, err := NewTemplateFile(filename)
		if err != nil {
			return nil, err
		}
		templateFiles = append(templateFiles, tf)
	}
	return templateFiles, nil
}

func generateTemplateData(src, dest, includes, layouts string) (TemplateData, error) {
	packageName := filepath.Base(filepath.Dir(dest))
	templateData := TemplateData{
		PackageName: packageName,
	}

	if includes != "" {
		prtty.Info.Println("--> parsing includes...")
		includes, err := ParseTemplateFiles(includes)
		if err != nil {
			return templateData, err
		}
		templateData.Includes = includes
	}
	if layouts != "" {
		prtty.Info.Println("--> parsing layouts...")
		layouts, err := ParseTemplateFiles(layouts)
		if err != nil {
			return templateData, err
		}
		templateData.Layouts = layouts
	}
	prtty.Info.Println("--> parsing templates...")
	templates, err := ParseTemplateFiles(src)
	if err != nil {
		return templateData, err
	}
	templateData.Templates = templates
	return templateData, nil
}

func writeFile(data TemplateData, dest string) error {
	prtty.Info.Println("--> generating go code...")
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}
	if !strings.HasSuffix(dest, ".go") {
		dest += ".go"
	}
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	prtty.Success.Printf("    CREATE %s", dest)
	if err := generatedTmpl.Execute(destFile, data); err != nil {
		return err
	}
	return nil
}

func formatFile(dest string) error {
	if _, err := exec.LookPath("gofmt"); err != nil {
		// gofmt is not installed or is not in PATH
		return nil
	}
	prtty.Default.Println("    formatting with gofmt...")
	output, err := exec.Command("gofmt", "-w", dest).CombinedOutput()
	if err != nil {
		return err
	}
	if output != nil && len(output) > 0 {
		prtty.Default.Printf("    %s", string(output))
	}
	return nil
}
