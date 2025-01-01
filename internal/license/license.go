package license

import (
	"bytes"
	"errors"
	"os"
	"text/template"

	t "go.wzykubek.xyz/licensmith/internal/template"
)

type License struct {
	ID      string
	Context Context
	Body    string
}

func (l *License) Gen() error {
	tmplPath := l.ID + ".tmpl"
	tmpl, err := t.Parse(tmplPath)
	if err != nil {
		return errors.New("usupported license")
	}

	body, err := template.New(l.ID).Parse(tmpl.Body)
  if err != nil {
    return err
  }

	var output bytes.Buffer
	body.Execute(&output, l.Context)

	l.Body = output.String()

	return nil
}

func (l *License) Write(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(l.Body); err != nil {
		return err
	}

	return nil
}
