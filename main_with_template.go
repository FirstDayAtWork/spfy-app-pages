package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/FirstDayAtWork/mustracker/views"
)

const myport int = 3222

func main() {
	tpl, err := views.Parse(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", "register.html"),
	)
	if err != nil {
		panic(err)
	}

	pageData := views.TemplateData{
		// Change this for a different title
		Title: "Register",
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
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	var handler func(w http.ResponseWriter, r *http.Request)
	handler = func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, pageData)
	}
	mux.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", myport), mux))
}
