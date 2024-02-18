package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FirstDayAtWork/mustracker/models"
	"github.com/FirstDayAtWork/mustracker/utils"
	"github.com/FirstDayAtWork/mustracker/views"
	"gorm.io/gorm"
)

// App represents web application.
type App struct {
	Th         *TemplateHandler
	Repository *models.Repository
	Auth       *Authorizer
}

/*
RegisterPOST handles POST request to /register endpoint. Algo:
1. Unmarshall request body to models.AccountData.
2. Perform validations for provided username, email and password.
3. Check if username is already taken.
4. Hash password.
5. Respond to client.
*/
func (ah *App) RegisterPOST(w http.ResponseWriter, r *http.Request) {
	var sr models.ServerResponse
	defer func() {
		jsonResp, err := sr.Marshall()
		if err != nil {
			// TODO LOG this
			fmt.Printf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		// TODO redirect to login???
	}()

	regData, err := RequestBodyToAccountData(r)
	if err != nil {
		fmt.Println("Error converting request data to account data", err)
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = fmt.Sprintf("Registration data parsing failed. Details: %v", err)
		return
	}
	// Data validations, doing 1 by 1 to have a more informative message in response
	// TODO Make it DRY, 1 func that returns different error messages!
	if !regData.IsValidUsername() {
		fmt.Printf("Username %s did not pass validation\n", regData.Username)
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = fmt.Sprintf(models.InvalidUsernameInput, regData.Username)
		return
	}
	if !regData.IsValidEmail() {
		fmt.Printf("Email %s did not pass validation\n", regData.Email)
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = fmt.Sprintf(models.InvalidEmailInput, regData.Email)
		return
	}
	if !regData.IsValidPassword() {
		fmt.Print("Password did not pass validation\n")
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = models.PasswordIsTooLongOrEmpty
		return
	}

	// TODO migrate this to cache
	userData, err := ah.Repository.GetAccountDataByUserame(regData.Username)
	if err != nil {
		fmt.Println("DB error when fetching user data", err)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			sr.Message = fmt.Sprintf("Internal server error. Details: %v", err)
			return
		}
	} else if userData != nil {
		w.WriteHeader(http.StatusConflict)
		sr.Message = fmt.Sprintf(models.UsernameAlreadyTakenMessage, userData.Username)
		return
	}

	hashedPassword, err := PasswordToHashedPassword(regData.Password)
	if err != nil {
		fmt.Println("Error hashing password", err)
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = models.InvalidPasswordInputMessage
		return
	}
	if !CheckPassword(regData.Password, hashedPassword) {
		fmt.Println("Hashed password does not match with raw password", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = models.PasswordHashAndPasswordMismatch
		return
	}

	regData.Password = hashedPassword
	regData.UUID = utils.NewUUID()
	if err := ah.Repository.CreateAccountData(regData); err != nil {
		fmt.Println("Error Writing User Data to DB", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = fmt.Sprintf("Internal server error. Details: %v", err)
	}
	// Happy path
	w.WriteHeader(http.StatusOK)
	sr.Message = models.SuccessMessage
}

func (ah *App) HandleRegister(w http.ResponseWriter, r *http.Request, acr *models.AuthCheckResult) {
	if acr.ValidAccess && acr.ValidRefresh {
		ah.AccountGET(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		ah.Th.Render(w, r)
	case http.MethodPost:
		ah.RegisterPOST(w, r)
	default:
		fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
	}
}

// TODO
func (ah *App) LoginPOST(w http.ResponseWriter, r *http.Request) {
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
		log.Println("Error converting request data to account data", err)
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = fmt.Sprintf("Registration data parsing failed. Details: %v", err)
		return
	}
	// lookup user
	user, err := ah.Repository.GetAccountDataByUserame(accData.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("user does not exist", err)
			w.WriteHeader(http.StatusNotFound)
			sr.Message = fmt.Sprintf(
				"User with %s username does not exist. Details: %v",
				accData.Username,
				err,
			)
			return
		}
		log.Println("DB error when fetching user data", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = fmt.Sprintf("Internal server error. Details: %v", err)
	}
	// check password
	if !CheckPassword(accData.Password, user.Password) {
		fmt.Println("Hashed password does not match with raw password", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = models.PasswordHashAndPasswordMismatch
		return
	}
	// Auth tokens creation and user shareout
	now := time.Now()
	refreshClms, err := ah.Auth.NewTokenClaims(
		user.UUID, user.Role, models.RefreshTokenType, &now,
	)
	if err != nil {
		log.Printf("refresh token generation error: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = "Authorization error: refresh"
		return
	}
	refreshTkn, err := ah.Auth.ClaimsToSignedString(refreshClms)
	if err != nil {
		log.Printf("refresh token generation error: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = "Authorization error: refresh"
		return
	}
	if err := ah.Repository.InvalidateUserTokens(user.UUID); err != nil {
		log.Printf("error revoking refresh token for %s\n", user.UUID)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = "Authorization error: refresh"
		return
	}
	err = ah.Repository.StoreToken(refreshTkn, refreshClms)
	if err != nil {
		log.Printf("database error: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = fmt.Sprintf("Internal server error. Details: %v", err)
		return
	}

	accessClms, err := ah.Auth.NewTokenClaims(
		user.UUID, user.Role, models.AccessTokenType, &now,
	)
	if err != nil {
		log.Printf("access token generation error: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = "Authorization error: access"
		return
	}
	accessTkn, err := ah.Auth.ClaimsToSignedString(accessClms)
	if err != nil {
		log.Printf("access token generation error: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = "Authorization error: access"
		return
	}
	// Cookie creation
	refreshCookie, err := ah.Auth.JWTTokenToCookie(refreshTkn, models.RefreshTokenType)
	if err != nil {
		log.Printf("refresh cookie generation error: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = "Authorization error: refresh"
		return
	}
	accessCookie, err := ah.Auth.JWTTokenToCookie(accessTkn, models.AccessTokenType)
	if err != nil {
		log.Printf("access cookie generation error: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = "Authorization error: access"
		return
	}
	http.SetCookie(w, refreshCookie)
	http.SetCookie(w, accessCookie)
	log.Println("Wrote both cookies to response")
	// Happy path
	w.WriteHeader(http.StatusOK)
	sr.Message = models.SuccessMessage
}

func (ah *App) HandleLogin(w http.ResponseWriter, r *http.Request, acr *models.AuthCheckResult) {
	if acr.ValidAccess && acr.ValidRefresh {
		ah.AccountGET(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		ah.Th.Render(w, r)
	case http.MethodPost:
		ah.LoginPOST(w, r)
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

func (ah *App) handleAuthorizationResult(
	w http.ResponseWriter,
	r *http.Request,
	authState *models.AuthCheckResult,
) *models.AuthHandlingResult {
	res := &models.AuthHandlingResult{}
	if !authState.NeedsHandling {
		res.IsAuthorized = true
		return res
	}
	if authState.Err != nil {
		return ah.Auth.handleAuthStateError(authState, res)
	}
	switch {
	case authState.ValidAccess && authState.ValidRefresh:
		if !authState.ValidRole {
			return res
		}
		res.IsAuthorized = true
		log.Printf("got a fully authorized request to %s\n", r.URL.Path)
		return res
	case authState.ValidAccess && !authState.ValidRefresh:
		return ah.Auth.handleMissingRefreshResult(authState, res)
	case !authState.ValidAccess && authState.ValidRefresh:
		newAccessCookie := ah.Auth.handleMissingAccessResult(authState, res, r.URL.Path)
		if newAccessCookie != nil {
			http.SetCookie(w, newAccessCookie)
			log.Println("Updated access cookie")
		}
		return res
	default:
		return ah.Auth.handleInvalidAuthResult(authState, res)
	}
}

func (ah *App) AccountGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Ciao Ragazzi, this is account section <h1>")
}

// ServeHTTP implements Handle interface.
func (ah *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authState := ah.Auth.CheckAuthorization(r)
	log.Printf("auth results check: %+v\n", authState)
	authHandleRes := ah.handleAuthorizationResult(w, r, authState)
	switch {
	case authHandleRes.Redirect != nil:
		http.Redirect(w, r, authHandleRes.Redirect.Path, authHandleRes.Redirect.Status)
		return
	case !authHandleRes.IsAuthorized:
		fmt.Fprintf(w, "Not authorized!\n")
		return
	}
	log.Printf("got an inboud to %s\n", r.URL.Path)
	// Only authorized requests get here
	switch r.URL.Path {
	case RegisterPath:
		ah.HandleRegister(w, r, authState)
	case LoginPath:
		ah.HandleLogin(w, r, authState)
	case AccountPath:
		ah.AccountGET(w, r)
	default:
		fmt.Fprintf(w, "ERROR! %s path is not supported!", r.URL.Path)
	}
}
