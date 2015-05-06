package main

import (
	"fmt"
	"github.com/albrow/prtty"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
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

func init() {

}

func main() {
	cmdBuild := &cobra.Command{
		Use:   "build <src> <dest>",
		Short: "Compile the templates in the src directory into a package in the dest directory.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				prtty.Error.Fatal("temple build requires exactly 2 arguments: the src directory and the dest directory.")
			}
			includes := cmd.Flag("includes").Value.String()
			if err := build(args[0], args[1], includes); err != nil {
				prtty.Error.Fatal(err)
			}
		},
	}
	cmdBuild.Flags().String("includes", "", "A directory which includes *.tmpl files to be shared between all templates.")

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
}

type TemplateFile struct {
	VarName string
	Name    string
	Source  string
}

func NewTemplateFile(filename string) (*TemplateFile, error) {
	// name is everything after the last slash, not including the file extension
	name := strings.TrimSuffix(filename[strings.LastIndex(filename, string(os.PathSeparator))+1:], ".tmpl")
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

func build(src, dest, includes string) error {
	prtty.Info.Println("--> building...")
	prtty.Default.Printf("    src: %s", src)
	prtty.Default.Printf("    dest: %s", dest)
	if includes != "" {
		prtty.Default.Printf("    includes: %s", includes)
	}

	// packageName is everything after the last path separator
	packageName := dest[strings.LastIndex(dest, string(os.PathSeparator))+1:]
	templateData := TemplateData{
		PackageName: packageName,
	}

	if includes != "" {
		prtty.Info.Println("--> parsing includes...")
		includes, err := ParseTemplateFiles(includes)
		if err != nil {
			return err
		}
		templateData.Includes = includes
	}
	prtty.Info.Println("--> parsing template files...")
	templates, err := ParseTemplateFiles(src)
	if err != nil {
		return err
	}
	templateData.Templates = templates

	prtty.Info.Println("--> generating go code...")
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}
	destFilename := filepath.Join(dest, packageName+".go")
	destFile, err := os.Create(destFilename)
	if err != nil {
		return err
	}
	prtty.Success.Printf("    CREATE %s", destFilename)
	if err := generatedTmpl.Execute(destFile, templateData); err != nil {
		return err
	}

	prtty.Info.Println("--> done")
	return nil
}
