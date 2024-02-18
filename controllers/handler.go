package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/FirstDayAtWork/mustracker/models"
	"github.com/FirstDayAtWork/mustracker/views"
	"gorm.io/gorm"
)

// AppHandler is a generic handler.
type AppHandler struct {
	Tpl        *views.Template
	Repository *models.Repository
}

/*
RegisterPOST handles POST request to /register endpoint. Algo:
1. Unmarshall request body to models.AccountData.
2. Perform validations for provided username, email and password.
3. Check if username is already taken.
4. Hash password.
5. Respond to client.
*/
func (rh *AppHandler) RegisterPOST(w http.ResponseWriter, r *http.Request) {
	var sr models.ServerResponse
	defer func() {
		jsonResp, err := sr.Marshall()
		if err != nil {
			// TODO LOG this
			fmt.Printf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
	}()

	regData, err := RequestBodyToAccountData(r)
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
	if !CheckPassword(regData.Password, hashedPassword) {
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

// RegisterGET handles a GET request to /register by rendering a login page.
func (rh *AppHandler) RegisterGET(w http.ResponseWriter, r *http.Request) {
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

func (rh *AppHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rh.RegisterGET(w, r)
	case http.MethodPost:
		rh.RegisterPOST(w, r)
	default:
		fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
	}
}

// TODO
func (rh *AppHandler) LoginPOST(w http.ResponseWriter, r *http.Request) {
	var sr models.ServerResponse
	defer func() {
		jsonResp, err := sr.Marshall()
		if err != nil {
			// TODO LOG this
			fmt.Printf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
	}()
	// Body to account data
	accData, err := RequestBodyToAccountData(r)
	if err != nil {
		fmt.Println("Error converting request data to account data", err)
		sr.StatusCode = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = fmt.Sprintf("Registration data parsing failed. Details: %v", err)
		return
	}
	// lookup user
	user, err := rh.Repository.GetAccountDataByUserame(accData.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("user does not exist", err)
			sr.StatusCode = http.StatusNotFound
			w.WriteHeader(http.StatusNotFound)
			sr.Message = fmt.Sprintf(
				"User with %s username does not exist. Details: %v",
				accData.Username,
				err,
			)
			return
		}
		fmt.Println("DB error when fetching user data", err)
		sr.StatusCode = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = fmt.Sprintf("Internal server error. Details: %v", err)
	}
	// check password
	if !CheckPassword(accData.Password, user.Password) {
		fmt.Println("Hashed password does not match with raw password", err)
		sr.StatusCode = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = models.PasswordHashAndPasswordMismatch
		return
	}
	// Happy path
	sr.StatusCode = http.StatusOK
	w.WriteHeader(http.StatusOK)
	sr.Message = models.SuccessMessage
	// TODO cookie & session
	// TODO redirect to login
}

// LoginGET handles a GET request to /login by rendering a login.
func (rh *AppHandler) LoginGET(w http.ResponseWriter, r *http.Request) {
	rh.Tpl.Execute(
		w,
		views.TemplateData{
			Title: views.LoginTitle,
			Styles: []string{
				views.TemplateCSS,
				views.LoginCSS,
			},
			Scripts: []string{
				views.TemplateJS,
				views.LoginJS,
			},
		},
	)
}

func (rh *AppHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rh.LoginGET(w, r)
	case http.MethodPost:
		rh.LoginPOST(w, r)
	default:
		fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
	}
}

// / Index
func (rh *AppHandler) indexGET(w http.ResponseWriter, r *http.Request) {
	rh.Tpl.Execute(
		w,
		views.TemplateData{
			Title: views.IndexTitle,
			Styles: []string{
				views.TemplateCSS,
				views.LoginCSS,
			},
			Scripts: []string{
				views.TemplateJS,
			},
		},
	)
}

func (rh *AppHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rh.indexGET(w, r)
	default:
		fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
	}
}

// About
func (rh *AppHandler) aboutGET(w http.ResponseWriter, r *http.Request) {
	rh.Tpl.Execute(
		w,
		views.TemplateData{
			Title: views.AboutTitle,
			Styles: []string{
				views.TemplateCSS,
				views.LoginCSS,
			},
			Scripts: []string{
				views.TemplateJS,
				views.AboutJS,
			},
		},
	)
}

func (rh *AppHandler) HandleAbout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rh.indexGET(w, r)
	default:
		fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
	}
}

// Donate
func (rh *AppHandler) donateGET(w http.ResponseWriter, r *http.Request) {
	rh.Tpl.Execute(
		w,
		views.TemplateData{
			Title: views.DonateTitle,
			Styles: []string{
				views.TemplateCSS,
				views.LoginCSS,
			},
			Scripts: []string{
				views.TemplateJS,
				views.DonateJS,
			},
		},
	)
}

func (rh *AppHandler) HandleDonate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rh.donateGET(w, r)
	default:
		fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
	}
}

// ServeHTTP implements Handle interface.
func (rh *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/register":
		rh.HandleRegister(w, r)
	case "/login":
		rh.HandleLogin(w, r)
	case "/index":
		rh.HandleIndex(w, r)
	case "/about":
		rh.HandleAbout(w, r)
	case "/donate":
		rh.HandleDonate(w, r)
	default:
		fmt.Fprintf(w, "ERROR! %s path is not supported!", r.URL.Path)
	}
}
