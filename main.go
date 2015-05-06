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
	var cmdBuild = &cobra.Command{
		Use:   "build <src> <dest>",
		Short: "Compile the templates in the src directory into a package in the dest directory.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				prtty.Error.Fatal("temple build requires exactly 2 arguments: the src directory and the dest directory.")
			}
			if err := build(args[0], args[1]); err != nil {
				prtty.Error.Fatal(err)
			}
		},
	}

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print the current version number.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	var rootCmd = &cobra.Command{
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
	Templates   []TemplateFile
}

type TemplateFile struct {
	Name   string
	Source string
}

func build(src, dest string) error {
	prtty.Info.Println("--> building...")
	prtty.Default.Printf("    src: %s", src)
	prtty.Default.Printf("    dest: %s", dest)

	// packageName is everything after the last path separator
	packageName := dest[strings.LastIndex(dest, string(os.PathSeparator)):]
	templateData := TemplateData{
		PackageName: packageName,
	}
	srcFiles, err := filepath.Glob(src + string(os.PathSeparator) + "*.tmpl")
	if err != nil {
		return err
	}
	for _, file := range srcFiles {
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		templateData.Templates = append(templateData.Templates, TemplateFile{
			Name:   strings.Title(strings.TrimPrefix(strings.TrimSuffix(file, ".tmpl"), src+"/")),
			Source: string(contents),
		})
	}

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
