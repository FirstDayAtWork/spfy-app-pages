package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/FirstDayAtWork/mustracker/controllers"
	"github.com/FirstDayAtWork/mustracker/templates"
	"github.com/FirstDayAtWork/mustracker/views"
)

const myport int = 3222

func main() {
	tpl := views.Must(
		views.ParseFS(
			templates.FS,
			filepath.Join("base.html"),
			filepath.Join("login.html"),
		),
	)

	pageData := &views.TemplateData{
		// Change this for a different title
		Title: "Login",
		// Change here to add more styles
		Styles: []string{
			"../static/styles/template.css",
			"../static/styles/login-style.css",
		},
		// Change here to add more scripts
		Scripts: []string{
			"../static/scripts/app-login.js",
		},
	}

	// Handler struct bc it implements ServerHTTP method
	st := controllers.Static{
		Template: tpl,
		Data:     pageData,
	}
	mux := http.NewServeMux()
	mux.Handle("/", st)
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", myport), mux))
}
