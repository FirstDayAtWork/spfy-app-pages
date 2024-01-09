package entity

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
)

const (
	templatesFolder  string = "templates"
	TemplateCSS      string = "../static/styles/template.css"
	LoginCSS         string = "../static/styles/login-style.css"
	RegisterCSS      string = "../static/styles/register-style.css"
	TemplateJS       string = "../static/scripts/app-html-template.js"
	LoginJS          string = "../static/scripts/app-login.js"
	RegisterJS       string = "../static/scripts/app-register.js"
	RegisterTitle    string = "Register"
	BaseTemplate     string = "base.html"
	RegisterTemplate string = "register.html"
)

type Page struct {
	Title    string
	Styles   []string
	Scripts  []string
	Content  string
	Template string
}

func (p *Page) Render(wr io.Writer) error {
	tmpl, err := template.ParseFiles(
		filepath.Join(templatesFolder, p.Template),
		filepath.Join(templatesFolder, p.Content),
	)
	if err != nil {
		fmt.Println("Error parsing template files")
		return err
	}
	err = tmpl.ExecuteTemplate(wr, p.Template, p)
	if err != nil {
		fmt.Printf("Error executing template: %s\n", p.Template)
	}
	err = tmpl.ExecuteTemplate(wr, p.Content, p)
	if err != nil {
		fmt.Printf("Error executing content template: %s\n", p.Content)
	}
	return nil
}
