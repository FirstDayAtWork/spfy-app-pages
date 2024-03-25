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

func (th TemplateHandler) Render(
	w http.ResponseWriter,
	r *http.Request,
	extra interface{},
) {
	tmpl, err := th.getTmpl(r.URL.Path)
	if err != nil {
		log.Printf("error getting template for %s: %v\n", r.URL.Path, err)
		fmt.Fprintf(w, "Error rendering page, we are so sorry!\n")
		return
	}
	log.Println("about to execute template...")
	tmpl.Execute(w, extra)
}

// New parses templates needed for the application and constructs the handler.
func NewTemplateHandler() *TemplateHandler {
	// Register
	regTmp := views.Must(
		views.ParseFS(
			templates.FS,
			filepath.Join(views.BaseTemplate),
			filepath.Join(views.RegisterTemplate),
			filepath.Join(views.GuestNavbarTemplate),
			filepath.Join(views.UserNavbarTemplate),
		),
	)
	regTmp.TemplateData = &views.TemplateData{
		Title: views.RegisterTitle,
		Scripts: []string{
			views.LoaderJS,
			views.TemplateJS,
			views.RegisterJS,
		},
		Styles: []string{
			views.TemplateCSS,
			views.BaseCSS,
		},
	}
	log.Println("prepared register template")

	// Login
	loginTmp := views.Must(
		views.ParseFS(
			templates.FS,
			filepath.Join(views.BaseTemplate),
			filepath.Join(views.LoginTemplate),
			filepath.Join(views.GuestNavbarTemplate),
			filepath.Join(views.UserNavbarTemplate),
		),
	)
	loginTmp.TemplateData = &views.TemplateData{
		Title: views.LoginTitle,
		Scripts: []string{
			views.LoaderJS,
			views.TemplateJS,
			views.LoginJS,
		},
		Styles: []string{
			views.TemplateCSS,
			views.BaseCSS,
		},
	}
	log.Println("prepared login template")

	// Index
	indexTmp := views.Must(
		views.ParseFS(
			templates.FS,
			filepath.Join(views.BaseTemplate),
			filepath.Join(views.IndexTemplate),
			filepath.Join(views.GuestNavbarTemplate),
			filepath.Join(views.UserNavbarTemplate),
		),
	)
	indexTmp.TemplateData = &views.TemplateData{
		Title: views.IndexTitle,
		Scripts: []string{
			views.LoaderJS,
			views.TemplateJS,
		},
		Styles: []string{
			views.TemplateCSS,
			views.BaseCSS,
		},
	}
	log.Println("prepared index template")

	aboutTmp := views.Must(
		views.ParseFS(
			templates.FS,
			filepath.Join(views.BaseTemplate),
			filepath.Join(views.AboutTemplate),
			filepath.Join(views.GuestNavbarTemplate),
			filepath.Join(views.UserNavbarTemplate),
		),
	)
	aboutTmp.TemplateData = &views.TemplateData{
		Title: views.AboutTitle,
		Scripts: []string{
			views.LoaderJS,
			views.TemplateJS,
			views.AboutJS,
		},
		Styles: []string{
			views.TemplateCSS,
			views.BaseCSS,
		},
	}
	log.Println("prepared about template")

	donateTmp := views.Must(
		views.ParseFS(
			templates.FS,
			filepath.Join(views.BaseTemplate),
			filepath.Join(views.DonateTemplate),
			filepath.Join(views.GuestNavbarTemplate),
			filepath.Join(views.UserNavbarTemplate),
		),
	)
	donateTmp.TemplateData = &views.TemplateData{
		Title: views.DonateTitle,
		Scripts: []string{
			views.LoaderJS,
			views.TemplateJS,
			views.DonateJS,
		},
		Styles: []string{
			views.TemplateCSS,
			views.BaseCSS,
		},
	}
	log.Println("prepared donate template")

	return &TemplateHandler{
		TemplateMap: map[string]*views.Template{
			RegisterPath: regTmp,
			LoginPath:    loginTmp,
			HomePath:     indexTmp,
			AboutPath:    aboutTmp,
			DonatePath:   donateTmp,
		},
	}

	// Donate

}
