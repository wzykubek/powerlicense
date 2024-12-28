package main

import (
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

type LicensingData struct {
	AuthorName  string
	AuthorEmail string
	Year        int
}

type LicenseData struct {
	FullName    string   `yaml:"title"`
	ID          string   `yaml:"spdx-id"`
	Description string   `yaml:"description"` // TODO
	Permissions []string `yaml:"permissions"` // TODO
	Limitations []string `yaml:"limitations"` // TODO
	Conditions  []string `yaml:"conditions"`  // TODO
}

//go:embed all:templates
var Templates embed.FS

var GitConfigError = errors.New("Can't read Git config")
var NotSupportedError = errors.New("Not supported license")
var InvalidFrontMatter = errors.New("Template front matter is invalid")

func getGitUserData() (string, string, error) {
	var userData [2]string
	for i, v := range []string{"user.name", "user.email"} {
		cmd := exec.Command("git", "config", "--get", v)
		out, err := cmd.Output()
		if err != nil {
			return "", "", GitConfigError
		}

		userData[i] = strings.TrimSpace(string(out))
	}

	return userData[0], userData[1], nil
}

func getTemplateList() []string {
	files, err := fs.ReadDir(Templates, "templates")
	if err != nil {
		panic(err)
	}

	var tmplList []string
	for _, v := range files {
		tmplList = append(tmplList, strings.Replace(v.Name(), ".tmpl", "", 1))
	}

	return tmplList
}

func listTemplates() {
	tmplList := getTemplateList()
	fmt.Println(strings.Join(tmplList, ", "))
}

func parseFrontMatter(tmplPath string) (LicenseData, string, error) {
	data, err := Templates.ReadFile(tmplPath)
	if err != nil {
		panic(err)
	}

	parts := strings.SplitN(string(data), "---", 3)
	if len(parts) < 3 {
		return LicenseData{}, "", InvalidFrontMatter
	}

	var licenseData LicenseData
	err = yaml.Unmarshal([]byte(parts[1]), &licenseData)
	if err != nil {
		return LicenseData{}, "", InvalidFrontMatter
	}

	return licenseData, strings.TrimSpace(parts[2]), nil
}

func genLicense(lcnsName string, lcnsData LicensingData, outFileName string) error {
	tmplFile := "templates/" + lcnsName + ".tmpl"
  _, lcnsBody, err := parseFrontMatter(tmplFile)
  if err != nil {
    panic(err)
  }

	tmpl, err := template.New(lcnsName).Parse(lcnsBody)
	if err != nil {
		return NotSupportedError
	}

	outFile, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, lcnsData)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	OutputFile := flag.String("output", "LICENSE", "Specify different output file")
	License := flag.String("license", "", "Choose a license")
	authorName := flag.String("name", "", "Set the author name")
	authorEmail := flag.String("email", "", "Set the author email")
	ListTemplates := flag.Bool("list", false, "List available licenses")
	flag.Parse()

	*License = strings.ToUpper(*License)

	if *ListTemplates {
		listTemplates()
		os.Exit(0)
	}

	if *License == "" {
		fmt.Printf("Error: No license specified\n\nUse --license LICENSE\n\nAvailable licenses:\n")
		listTemplates()
		os.Exit(1)
	}

	if *authorName == "" || *authorEmail == "" {
		var err error
		*authorName, *authorEmail, err = getGitUserData()
		if err != nil {
			if errors.Is(err, GitConfigError) {
				fmt.Printf(
					"Error: Can't read Git config.\n\nUse --name \"NAME\" and --email EMAIL instead.\n",
				)
				os.Exit(3)
			}
		}
	}

	lcnsData := LicensingData{
		AuthorName:  *authorName,
		AuthorEmail: *authorEmail,
		Year:        time.Now().Year(),
	}

	err := genLicense(*License, lcnsData, *OutputFile)
	if err != nil {
		if errors.Is(err, NotSupportedError) {
			fmt.Printf("Error: There is no '%s' license\n\nAvailable licenses:\n", *License)
			listTemplates()
			os.Exit(2)
		}
	}
}
