package views

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
)

type Page struct {
	Title   string
	Styles  []string
	Scripts []string
	Content string
	Base    string
}

func (p *Page) Render(wr io.Writer) error {
	tmpl, err := template.ParseFiles(
		filepath.Join(templatesFolder, p.Base),
		filepath.Join(templatesFolder, p.Content),
	)
	if err != nil {
		fmt.Println("Error parsing template files")
		return err
	}
	err = tmpl.ExecuteTemplate(wr, p.Base, p)
	if err != nil {
		fmt.Printf("Error executing template: %s\n", p.Base)
	}
	err = tmpl.ExecuteTemplate(wr, p.Content, p)
	if err != nil {
		fmt.Printf("Error executing content template: %s\n", p.Content)
	}
	return nil
}
