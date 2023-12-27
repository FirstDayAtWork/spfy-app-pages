package basic_server

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

const (
	pagesDir     = "pages"
	registerPage = "register.html"
)

func RenderRegister(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(filepath.Join(pagesDir, registerPage))
	if err != nil {
		fmt.Printf("Error parsing %s: %v\n", registerPage, err)
		return
	}
	err = t.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		fmt.Printf("Error executing template %s: %v\n", registerPage, err)
		return
	}
}
