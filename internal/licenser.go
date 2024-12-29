package internal

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type Licenser struct {
	LicenseID      string
	LicenseContext LicenseContext
	OutputFile     string
	licenseBody    string
}

func (l *Licenser) ParseTemplate() (LicenseTemplate, error) {
	tmplPath := "templates/" + l.LicenseID + ".tmpl"
	data, err := EmbedFS.ReadFile(tmplPath)
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
		return errors.New("usupported license")
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
