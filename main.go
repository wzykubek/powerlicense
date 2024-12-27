package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

type LicensingData struct {
	AuthorName  string
	AuthorEmail string
	Year        int
}

var GitConfigError = errors.New("Can't read Git config")
var NotSupportedError = errors.New("Not supported license")

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
	d, err := os.Open("templates")
	if err != nil {
		panic(err)
	}
	files, err := d.Readdir(0)
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

func genLicense(lcnsName string, lcnsData LicensingData, outFileName string) error {
	tmplFile := "templates/" + lcnsName + ".tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
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
