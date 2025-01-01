package templates

import (
	"embed"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed *.tmpl
var FS embed.FS

type Template struct {
	Title       string   `yaml:"title"`
	ID          string   `yaml:"spdx-id"`
	Description string   `yaml:"description"` // TODO
	Permissions []string `yaml:"permissions"` // TODO
	Limitations []string `yaml:"limitations"` // TODO
	Conditions  []string `yaml:"conditions"`  // TODO
	Body        string
}

func Parse(path string) (Template, error) {
	data, err := FS.ReadFile(path)
	if err != nil {
		return Template{}, err
	}

	parts := strings.SplitN(string(data), "---", 3)

	var tmpl Template
	yaml.Unmarshal([]byte(parts[1]), &tmpl)
	tmpl.Body = strings.TrimSpace(parts[2])

	return tmpl, nil
}

func List() ([]string, error) {
	files, err := fs.ReadDir(FS, ".")
	if err != nil {
		return []string{}, err
	}

	var templates []string
	for _, v := range files {
		templates = append(templates, strings.Replace(v.Name(), ".tmpl", "", 1))
	}

	return templates, nil
}
