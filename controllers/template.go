package controllers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/FirstDayAtWork/mustracker/templates"
	"github.com/FirstDayAtWork/mustracker/views"
)

type TemplateHandler struct {
	TemplateMap map[string]*views.Template
}

func (th TemplateHandler) getTmpl(path string) (*views.Template, error) {
	tmpl, ok := th.TemplateMap[path]
	if !ok {
		return nil, fmt.Errorf("%s path does not have a template", path)
	}
	return tmpl, nil
}

func (th TemplateHandler) Render(w http.ResponseWriter, r *http.Request) {
	tmpl, err := th.getTmpl(r.URL.Path)
	if err != nil {
		log.Printf("error getting template for %s: %v\n", r.URL.Path, err)
		fmt.Fprintf(w, "Error rendering page, we are so sorry!\n")
		return
	}
	tmpl.Execute(w)
}

// New parses templates needed for the application and constructs the handler.
func NewTemplateHandler() *TemplateHandler {
	// Register
	regTmp := views.Must(
		views.ParseFS(
			templates.FS,
			filepath.Join(views.BaseTemplate),
			filepath.Join(views.RegisterTemplate),
		),
	)
	regTmp.TemplateData = &views.TemplateData{
		Title: views.RegisterTitle,
		Scripts: []string{
			views.RegisterJS,
		},
		Styles: []string{
			views.TemplateCSS,
			views.LoginCSS,
		},
	}
	log.Println("prepared register template")

	// Login
	loginTmp := views.Must(
		views.ParseFS(
			templates.FS,
			filepath.Join(views.BaseTemplate),
			filepath.Join(views.LoginTemplate),
		),
	)
	loginTmp.TemplateData = &views.TemplateData{
		Title: views.LoginTitle,
		Scripts: []string{
			views.LoginJS,
		},
		Styles: []string{
			views.TemplateCSS,
			views.LoginCSS,
		},
	}
	log.Println("prepared login template")
	return &TemplateHandler{
		TemplateMap: map[string]*views.Template{
			RegisterPath: regTmp,
			LoginPath:    loginTmp,
		},
	}
}
