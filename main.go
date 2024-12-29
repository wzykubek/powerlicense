package main

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

//go:embed all:templates
var TemplatesDir embed.FS

type LicenseTemplate struct {
	Title       string   `yaml:"title"`
	ID          string   `yaml:"spdx-id"`
	Description string   `yaml:"description"` // TODO
	Permissions []string `yaml:"permissions"` // TODO
	Limitations []string `yaml:"limitations"` // TODO
	Conditions  []string `yaml:"conditions"`  // TODO
	Body        string
}

type LicenseContext struct {
	AuthorName  string
	AuthorEmail string
	Year        int
}

type Licenser struct {
	LicenseID      string
	LicenseContext LicenseContext
	OutputFile     string
	licenseBody    string
}

func NewLicenseContext(authorName string, authorEmail string) (LicenseContext, error) {
	var err error
	if authorName == "" {
		authorName, err = gitUserData("user.name")
	}
	if authorEmail == "" {
		authorEmail, err = gitUserData("user.email")
	}
	if err != nil {
		return LicenseContext{}, err
	}

	return LicenseContext{
		AuthorName:  authorName,
		AuthorEmail: authorEmail,
		Year:        time.Now().Year(),
	}, nil
}

func (l *Licenser) ParseTemplate() (LicenseTemplate, error) {
	licenseID := strings.ToUpper(l.LicenseID)
	tmplPath := "templates/" + licenseID + ".tmpl"
	data, err := TemplatesDir.ReadFile(tmplPath)
	if err != nil {
		return LicenseTemplate{}, err
	}

	parts := strings.SplitN(string(data), "---", 3)

	var licenseTmpl LicenseTemplate
	yaml.Unmarshal([]byte(parts[1]), &licenseTmpl)
	licenseTmpl.Body = strings.TrimSpace(parts[2])

	return licenseTmpl, nil
}

func (l *Licenser) Generate() error {
	license, err := l.ParseTemplate()
	if err != nil {
		return errors.New("Not supported license")
	}

	tmpl, _ := template.New(l.LicenseID).Parse(license.Body)

	var output bytes.Buffer
	tmpl.Execute(&output, l.LicenseContext)

	l.licenseBody = output.String()

	return nil
}

func (l *Licenser) WriteFile() error {
	outFile, err := os.Create(l.OutputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := outFile.WriteString(l.licenseBody); err != nil {
		return err
	}

	return nil
}

func gitUserData(key string) (string, error) {
	cmd := exec.Command("git", "config", "--get", key)
	out, err := cmd.Output()
	if err != nil {
		return "", errors.New("Can't read Git config")
	}

	value := strings.TrimSpace(string(out))
	return value, nil
}

func templateList() []string {
	files, err := fs.ReadDir(TemplatesDir, "templates")
	if err != nil {
		panic(err)
	}

	var tmplList []string
	for _, v := range files {
		tmplList = append(tmplList, strings.Replace(v.Name(), ".tmpl", "", 1))
	}

	return tmplList
}

func listLicenses() {
	tmplList := templateList()
	fmt.Println(strings.Join(tmplList, ", "))
}

func main() {
	OutputFile := flag.String("output", "LICENSE", "Specify different output file")
	LicenseID := flag.String("license", "", "Specify license by SPDX ID (e.g. BSD-3-Clause)")
	AuthorName := flag.String("name", "", "Set the author name (read from Git by default)")
	AuthorEmail := flag.String("email", "", "Set the author email (read from Git by default)")
	ListLicenses := flag.Bool("list", false, "List available licenses")
	flag.Parse()

	if *ListLicenses {
		listLicenses()
		os.Exit(0)
	}

	if *LicenseID == "" {
		fmt.Printf("Error: No license specified\n\nUse --license LICENSE\n\nAvailable licenses:\n")
		listLicenses()
		os.Exit(1)
	}

	licenseCtx, err := NewLicenseContext(*AuthorName, *AuthorEmail)
	if err != nil && err.Error() == "Can't read Git config" {
		fmt.Printf(
			"Error: Can't read Git config.\n\nUse --name \"NAME\" and --email EMAIL instead.\n",
		)
		os.Exit(3)
	}

	licenser := Licenser{
		LicenseID:      *LicenseID,
		LicenseContext: licenseCtx,
		OutputFile:     *OutputFile,
	}

	err = licenser.Generate()
	if err != nil && err.Error() == "Not supported license" {
		fmt.Printf("Error: There is no '%s' license\n\nAvailable licenses:\n", *LicenseID)
		listLicenses()
		os.Exit(2)
	}

	if err = licenser.WriteFile(); err != nil {
		panic(err)
	}
}
