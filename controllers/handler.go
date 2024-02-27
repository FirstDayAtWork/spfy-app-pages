package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FirstDayAtWork/mustracker/models"
	"github.com/FirstDayAtWork/mustracker/utils"
	"gorm.io/gorm"
)

// App represents web application.
type App struct {
	Th         *TemplateHandler
	Repository *models.Repository
	Auth       *Authorizer
}

/*
registerPOST handles POST request to /register endpoint. Algo:
1. Unmarshall request body to models.AccountData.
2. Perform validations for provided username, email and password.
3. Check if username is already taken.
4. Hash password.
5. Respond to client.
*/
func (ah *App) registerPOST(w http.ResponseWriter, r *http.Request) {
	var sr models.ServerResponse
	defer func() {
		sr.SendToClient(w)
	}()
	regData, err := RequestBodyToAccountData(r)
	if err != nil {
		log.Printf("error converting request data to account data: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = fmt.Sprintf("Registration data parsing failed. Details: %v", err)
		return
	}
	// Data validations, doing 1 by 1 to have a more informative message in response
	// TODO Make it DRY, 1 func that returns different error messages!
	validationErr := AccountDataToErrorMessage(regData)
	if validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		sr.Message = validationErr.Error()
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
	sr.Message = models.SuccessMessage
}

// handleRegister handles GET anb POST requests to /register.
func (ah *App) handleRegister(w http.ResponseWriter, r *http.Request, authState *models.AuthCheckResult) {
	if authState.HasValidTokens() {
		http.Redirect(w, r, AccountPath, http.StatusSeeOther)
		return
	}
	switch r.Method {
	case http.MethodGet:
		ah.Th.Render(w, r, struct{ IsGuest bool }{true})
	case http.MethodPost:
		ah.registerPOST(w, r)
	default:
		fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
	}
}

/*
loginPOST handles POST request to /login endpoint. Algo:
1. Unmarshall request body to models.AccountData.
2. Get user data from DB by the provided username.
3. Validate password.
4. Generate a new refresh token.
5. Invalidate existing refresh tokens of a user present in DB.
6. Record the refresh token to DB.
5. Generate a new access token.
6. Write both tokens to client's cookies.
7. Respond to client.
*/
func (ah *App) loginPOST(
	w http.ResponseWriter,
	r *http.Request,
) {
	var sr models.ServerResponse
	defer func() {
		sr.SendToClient(w)
	}()
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
			w.WriteHeader(http.StatusUnauthorized)
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
		return
	}
	// check password
	if !CheckPassword(accData.Password, user.Password) {
		fmt.Println("Hashed password does not match with raw password", err)
		w.WriteHeader(http.StatusUnauthorized)
		sr.Message = models.PasswordHashAndPasswordMismatch
		return
	}
	// Auth tokens creation and user shareout
	now := time.Now()
	refreshCookie, err := ah.Auth.GrantNewToken(user, models.RefreshTokenType, &now)
	if err != nil {
		log.Printf("refresh token generation error: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = err.Error()
	}
	accessCookie, err := ah.Auth.GrantNewToken(user, models.AccessTokenType, &now)
	if err != nil {
		log.Printf("access token generation error: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		sr.Message = err.Error()
	}
	http.SetCookie(w, refreshCookie)
	http.SetCookie(w, accessCookie)
	log.Println("Wrote both cookies to response")
	// Happy path
	sr.Message = models.SuccessMessage
}

/*
handleRegister handles GET anb POST requests to /register. Algo:
1. If AuthCheckResult indicates client having access & refresh, redirect to account with a 303.
2. Depending on request method, either render a page or execute loginPost().
*/
func (ah *App) handleLogin(w http.ResponseWriter, r *http.Request, authState *models.AuthCheckResult) {
	if authState.ValidAccess && authState.ValidRefresh {
		http.Redirect(w, r, AccountPath, http.StatusSeeOther)
		return
	}
	switch r.Method {
	case http.MethodGet:
		ah.Th.Render(w, r, struct{ IsGuest bool }{true})
	case http.MethodPost:
		ah.loginPOST(w, r)
	default:
		fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
	}
}

func (ah *App) handleIndex(w http.ResponseWriter, r *http.Request, authState *models.AuthCheckResult) {
	switch r.Method {
	case http.MethodGet:
		ah.Th.Render(w, r, AuthStateToExtra(authState))
	default:
		fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
	}
}

func (ah *App) handleAbout(w http.ResponseWriter, r *http.Request, authState *models.AuthCheckResult) {
	switch r.Method {
	case http.MethodGet:
		ah.Th.Render(w, r, AuthStateToExtra(authState))
	default:
		fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
	}
}

func (ah *App) handleDonate(w http.ResponseWriter, r *http.Request, authState *models.AuthCheckResult) {
	switch r.Method {
	case http.MethodGet:
		ah.Th.Render(w, r, AuthStateToExtra(authState))
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

func (ah *App) accountGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Ciao Ragazzi, this is account section <h1>")
}

// ServeHTTP implements Handle interface.
func (ah *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("got an inboud to %s\n", r.URL.Path)
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
	// Only authorized requests get here
	switch r.URL.Path {
	case RegisterPath:
		ah.handleRegister(w, r, authState)
	case LoginPath:
		ah.handleLogin(w, r, authState)
	case AccountPath:
		ah.accountGET(w, r)
	case HomePath:
		ah.handleIndex(w, r, authState)
	case AboutPath:
		ah.handleAbout(w, r, authState)
	case DonatePath:
		ah.handleDonate(w, r, authState)
	default:
		fmt.Fprintf(w, "ERROR! %s path is not supported!", r.URL.Path)
	}
}
