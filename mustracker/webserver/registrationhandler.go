package webserver

import (
	"mustracker/entity"
	"mustracker/mapper"
	"net/http"

	"fmt"
	"html/template"
	"path/filepath"

	"gorm.io/gorm"
)

const (
	pagesDir     = "pages"
	registerPage = "register.html"
)

type DataHandler struct {
	DB *gorm.DB
}

func (dh *DataHandler) RecordRegistration(w http.ResponseWriter, r *http.Request) error {
	regData, err := mapper.RegistrationRequestToRegistrationData(r)
	if err != nil {
		// TODO figure out what status to respond with here
		return err
	}
	if err = dh.insertRegistration(regData); err != nil {
		return err
	}
	return nil
}

func (dh *DataHandler) insertRegistration(rd *entity.RegistrationData) error {
	res := dh.DB.Create(rd)
	if res.Error != nil {
		// Logging?
		return res.Error
	}
	return nil
}

func (dh *DataHandler) RenderRegister(w http.ResponseWriter, r *http.Request) {
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

func (dh *DataHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		dh.RenderRegister(w, r)
	case http.MethodPost:
		dh.RecordRegistration(w, r)
	default:
		fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
	}
}
