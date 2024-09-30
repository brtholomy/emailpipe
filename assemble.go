package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"text/template"
)

var localhost string = "http://localhost:1313/"
var prodhost string = "https://www.bartholomy.ooo/"

func CreateTemplate(name, t string) *template.Template {
	return template.Must(template.New(name).Parse(t))
}

func Assemble(post *Post, opts *Options) (*string, error) {
	dat, err := os.ReadFile(opts.Template)
	if err != nil {
		return nil, err
	}

	t := CreateTemplate("t", string(dat))
	var buf bytes.Buffer
	if err := t.Execute(io.Writer(&buf), post); err != nil {
		return nil, err
	}

	c := buf.String()
	c = strings.ReplaceAll(c, localhost, prodhost)
	if strings.Contains(c, "localhost") {
		return nil, errors.New("hostname replace failed")
	}
	return &c, nil
}
