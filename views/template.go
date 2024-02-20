package views

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
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

func ParseFS(fs fs.FS, patterns ...string) (*Template, error) {
	htmlTmpl, err := template.ParseFS(fs, patterns...)
	if err != nil {
		// using errorf for extra context
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	return &Template{
		htmlTmpl: htmlTmpl,
	}, nil
}

func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

type TemplateData struct {
	Title   string
	Styles  []string
	Scripts []string
	Extra   interface{}
}

type Template struct {
	htmlTmpl     *template.Template
	TemplateData *TemplateData
}

func (t *Template) Execute(w http.ResponseWriter, extra interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if extra != nil {
		t.TemplateData.Extra = extra
	}
	log.Printf("Executing template with data: %+v\n", t.TemplateData)
	err := t.htmlTmpl.Execute(w, t.TemplateData) // is it ok?
	if err != nil {
		log.Printf("Error executing template : %v\n", err)
		http.Error(w, TemplateRenderError, http.StatusInternalServerError)
		return
	}
}
