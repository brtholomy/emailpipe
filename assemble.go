package main

import (
	"bytes"
	"io"
	"os"
	"text/template"
)

func CreateTemplate(name, t string) *template.Template {
	return template.Must(template.New(name).Parse(t))
}

func Assemble(post *Post, tmpl_path *string) (*string, error) {
	dat, err := os.ReadFile(*tmpl_path)
	if err != nil {
		return nil, err
	}

	t := CreateTemplate("t", string(dat))
	var buf bytes.Buffer
	if err := t.Execute(io.Writer(&buf), post); err != nil {
		return nil, err
	}
	c := buf.String()
	return &c, nil
}
