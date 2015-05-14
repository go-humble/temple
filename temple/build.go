package temple

import (
	"errors"
	"github.com/albrow/prtty"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	if dirs.partials != "" {
		prtty.Default.Println("    checking partials...")
		if err := AddPartialFiles(dirs.partials); err != nil {
			return err
		}
	}
	if dirs.layouts != "" {
		prtty.Default.Println("    checking layouts...")
		if err := AddLayoutFiles(dirs.layouts); err != nil {
			return err
		}
	}
	prtty.Default.Println("    checking templates...")
	if err := AddTemplateFiles(dirs.templates); err != nil {
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
	if err := formatFile(dest); err != nil {
		return err
	}
	return nil
}

func (data *templateData) writeToFile(dest string) error {
	prtty.Success.Printf("    created %s", dest)
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	// NOTE: GOPATH might consist of multiple paths. If that is the case, we look in the first path.
	gopath := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))[0]
	templePath := filepath.Join(gopath, "src", "github.com", "go-humble", "temple", "temple")
	generatedTmpl := template.Must(template.ParseFiles(filepath.Join(templePath, "generated.go.tmpl")))
	return generatedTmpl.Execute(destFile, data)
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
