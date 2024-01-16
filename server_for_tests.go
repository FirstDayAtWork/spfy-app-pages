package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FirstDayAtWork/mustracker/views"
)

const port int = 3222

func main() {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))
	mux.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func handler(w http.ResponseWriter, r *http.Request) {

	// Create a page struct with all fields populated
	regPage := views.Page{
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
		Content: "login.html", // change this for new pages
		Base:    "base.html",  // don't change this!
	}
	regPage.Render(w)

}
