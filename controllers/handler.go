package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/FirstDayAtWork/mustracker/models"
	"github.com/FirstDayAtWork/mustracker/views"
	"gorm.io/gorm"
)

type RegisterHandler struct {
	Tpl        *views.Template
	Repository *models.Repository
}

func (rh *RegisterHandler) RegisterPOST(w http.ResponseWriter, r *http.Request) {
	var sr models.ServerResponse
	defer func() {
		jsonResp, err := sr.Marshall()
		if err != nil {
			// TODO LOG this
			fmt.Printf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
	}()

	regData, err := RegistrationRequestToAccountData(r)
	if err != nil {
		fmt.Println("Error converting request data to account data", err)
		sr.StatusCode = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = fmt.Sprintf("Registration data parsing failed. Details: %v", err)
		return
	}
	// Data validations, doing 1 by 1 to have a more informative message in response
	// TODO Make it DRY
	if !regData.IsValidUsername() {
		fmt.Printf("Username %s did not pass validation\n", regData.Username)
		sr.StatusCode = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = fmt.Sprintf(models.InvalidUsernameInput, regData.Username)
		return
	}
	if !regData.IsValidEmail() {
		fmt.Printf("Email %s did not pass validation\n", regData.Email)
		sr.StatusCode = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = fmt.Sprintf(models.InvalidEmailInput, regData.Email)
		return
	}
	if !regData.IsValidPassword() {
		fmt.Print("Password did not pass validation\n")
		sr.StatusCode = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = models.PasswordIsTooLongOrEmpty
		return
	}

	// TODO migrate this to cache
	userData, err := rh.Repository.GetAccountDataByUserame(regData.Username)
	if err != nil {
		fmt.Println("DB error when fetching user data", err)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			sr.StatusCode = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			sr.Message = fmt.Sprintf("Internal server error. Details: %v", err)
			return
		}
	} else if userData != nil {
		sr.StatusCode = http.StatusConflict
		w.WriteHeader(http.StatusConflict)
		sr.Message = fmt.Sprintf(models.UsernameAlreadyTakenMessage, userData.Username)
		return
	}

	hashedPassword, err := PasswordToHashedPassword(regData.Password)
	if err != nil {
		fmt.Println("Error hashing password", err)
		sr.StatusCode = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = models.InvalidPasswordInputMessage
		return
	}
	passMatch := CheckPassword(regData.Password, hashedPassword)
	if !passMatch {
		fmt.Println("Hashed password does not match with raw password", err)
		sr.StatusCode = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = models.PasswordHashAndPasswordMismatch
		return
	}

	regData.Password = hashedPassword
	if err := rh.Repository.CreateAccountData(regData); err != nil {
		fmt.Println("Error Writing User Data to DB", err)
		sr.StatusCode = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = fmt.Sprintf("Internal server error. Details: %v", err)
	}
	// Happy path
	sr.StatusCode = http.StatusOK
	w.WriteHeader(http.StatusOK)
	sr.Message = models.SuccessMessage
}

func (rh *RegisterHandler) RegisterGET(w http.ResponseWriter, r *http.Request) {
	rh.Tpl.Execute(
		w,
		views.TemplateData{
			Title: views.RegisterTitle,
			Styles: []string{
				views.TemplateCSS,
				views.LoginCSS,
			},
			Scripts: []string{
				views.RegisterJS,
			},
		},
	)
}

func (rh *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rh.RegisterGET(w, r)
	case http.MethodPost:
		rh.RegisterPOST(w, r)
	default:
		fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
	}
}
