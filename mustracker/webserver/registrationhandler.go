package webserver

import (
	"errors"
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

func (dh *DataHandler) RecordRegistration(w http.ResponseWriter, r *http.Request) {
	var sr entity.ServerResponse
	// Defer to avoid repetitive code
	defer func() {
		jsonResp, err := sr.Marshall()
		if err != nil {
			// TODO LOG this
			fmt.Printf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
	}()

	regData, err := mapper.RegistrationRequestToAccountData(r)
	if err != nil {
		fmt.Println("Error converting request data to account data", err)
		sr.StatusCode = http.StatusBadRequest
		sr.Message = fmt.Sprintf("Registration data parsing failed. Details: %v", err)
		return
	}

	// TODO migrate this to cache
	userData, err := dh.getAccountDataByUserame(regData.Username)
	if err != nil {
		fmt.Println("DB error when fetching user data", err)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			sr.StatusCode = http.StatusInternalServerError
			sr.Message = fmt.Sprintf("Internal server error. Details: %v", err)
			return
		}
	} else if userData != nil {
		sr.StatusCode = http.StatusConflict
		sr.Message = fmt.Sprintf(entity.UsernameAlreadyTakenMessage, userData.Username)
		return
	}

	hashedPassword, err := mapper.PasswordToHashedPassword(regData.Password)
	if err != nil {
		fmt.Println("Error hashing password", err)
		sr.StatusCode = http.StatusBadRequest
		sr.Message = entity.InvalidPasswordInputMessage
		return
	}
	passMatch := mapper.CheckPassword(regData.Password, hashedPassword)
	if !passMatch {
		fmt.Println("Hashed password does not match with raw password", err)
		sr.StatusCode = http.StatusInternalServerError
		sr.Message = entity.PasswordHashAndPasswordMismatch
		return
	}

	regData.Password = hashedPassword
	if err := dh.insertRegistration(regData); err != nil {
		fmt.Println("Error Writing User Data to DB", err)
		sr.StatusCode = http.StatusInternalServerError
		sr.Message = fmt.Sprintf("Internal server error. Details: %v", err)
	}
	// Happy path
	sr.StatusCode = http.StatusOK
	sr.Message = entity.SuccessMessage

}

func (dh *DataHandler) insertRegistration(ad *entity.AccountData) error {
	res := dh.DB.Create(ad)
	if res.Error != nil {
		// Logging?
		return res.Error
	}
	return nil
}

func (dh *DataHandler) getAccountDataByUserame(username string) (*entity.AccountData, error) {
	// Create a dummy struct for query filters
	resultData := &entity.AccountData{}
	res := dh.DB.Where(
		&entity.AccountData{Username: username},
	).First(resultData)

	if res.Error != nil {
		return nil, res.Error
	}

	return resultData, nil

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
