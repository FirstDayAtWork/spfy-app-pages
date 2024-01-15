package views

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Parse(filepath ...string) (*Template, error) {
	if len(filepath) == 0 {
		return nil, errors.New(NoFilePathsError)
	}
	htmlTmpl, err := template.ParseFiles(filepath...)
	if err != nil {
		// using errorf for extra context
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	return &Template{
		htmlTmpl: htmlTmpl,
	}, nil
}

type Template struct {
	htmlTmpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTmpl.Execute(w, data) // is it ok?
	if err != nil {
		log.Printf("Error executing template : %v\n", err)
		http.Error(w, TemplateRenderError, http.StatusInternalServerError)
		return
	}
}

type TemplateData struct {
	Title   string
	Styles  []string
	Scripts []string
}
