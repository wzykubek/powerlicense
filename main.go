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

type InputData struct {
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
var TemplatesDir embed.FS

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
	licList := getTemplateList()
	fmt.Println(strings.Join(licList, ", "))
}

func parseFrontMatter(tmplPath string) (LicenseData, string, error) {
	data, err := TemplatesDir.ReadFile(tmplPath)
	if err != nil {
		panic(err)
	}

	parts := strings.SplitN(string(data), "---", 3)
	if len(parts) < 3 {
		return LicenseData{}, "", InvalidFrontMatter
	}

	var licData LicenseData
	err = yaml.Unmarshal([]byte(parts[1]), &licData)
	if err != nil {
		return LicenseData{}, "", InvalidFrontMatter
	}

	return licData, strings.TrimSpace(parts[2]), nil
}

func genLicense(licName string, inputData InputData, outFileName string) error {
	tmplPath := "templates/" + licName + ".tmpl"
	_, lcnsBody, err := parseFrontMatter(tmplPath)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New(licName).Parse(lcnsBody)
	if err != nil {
		return NotSupportedError
	}

	outFile, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, inputData)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	OutputFile := flag.String("output", "LICENSE", "Specify different output file")
	LicenseName := flag.String("license", "", "Choose a license")
	AuthorName := flag.String("name", "", "Set the author name")
	AuthorEmail := flag.String("email", "", "Set the author email")
	ListLicenses := flag.Bool("list", false, "List available licenses")
	flag.Parse()

	*LicenseName = strings.ToUpper(*LicenseName)

	if *ListLicenses {
		listLicenses()
		os.Exit(0)
	}

	if *LicenseName == "" {
		fmt.Printf("Error: No license specified\n\nUse --license LICENSE\n\nAvailable licenses:\n")
		listLicenses()
		os.Exit(1)
	}

	if *AuthorName == "" || *AuthorEmail == "" {
		var err error
		*AuthorName, *AuthorEmail, err = getGitUserData()
		if err != nil {
			if errors.Is(err, GitConfigError) {
				fmt.Printf(
					"Error: Can't read Git config.\n\nUse --name \"NAME\" and --email EMAIL instead.\n",
				)
				os.Exit(3)
			}
		}
	}

	inputData := InputData{
		AuthorName:  *AuthorName,
		AuthorEmail: *AuthorEmail,
		Year:        time.Now().Year(),
	}

	err := genLicense(*LicenseName, inputData, *OutputFile)
	if err != nil {
		if errors.Is(err, NotSupportedError) {
			fmt.Printf("Error: There is no '%s' license\n\nAvailable licenses:\n", *LicenseName)
			listLicenses()
			os.Exit(2)
		}
	}
}
