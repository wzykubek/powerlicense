package main

import (
	"embed"
	"io/fs"
	"strings"
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

func listTemplates() []string {
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
	tmplList := listTemplates()
	println(strings.Join(tmplList, ", "))
}
