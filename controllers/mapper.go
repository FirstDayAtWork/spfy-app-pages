package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/FirstDayAtWork/mustracker/models"
)

// RequestBodyToAccountData unmarshalls request body to models.AccountData struct.
func RequestBodyToAccountData(r *http.Request) (*models.AccountData, error) {
	rd := &models.AccountData{}
	err := json.NewDecoder(r.Body).Decode(rd)
	if err != nil {
		// TODO log this
		return nil, err
	}
	if rd.Role == 0 {
		log.Println("Role not provided, defaulting to:", models.UserRole)
		rd.Role = models.UserRole
	}
	return rd, nil
}

func AuthStateToExtra(authState *models.AuthCheckResult) interface{} {
	// TODO type it statically at some point
	return struct {
		IsGuest  bool
		UserName string
	}{
		IsGuest:  authState.IsGuest(),
		UserName: authState.GetUserName(),
	}
}

func RegistrationDataToErrorMessage(ad *models.AccountData) error {
	if !ad.IsValidUsername() {
		return fmt.Errorf(models.InvalidUsernameInput, ad.Username)
	}
	if !ad.IsValidEmail() {
		return fmt.Errorf(models.InvalidEmailInput, ad.Email)
	}
	if !ad.IsValidPassword() {
		return fmt.Errorf(models.PasswordIsTooLongOrEmpty)
	}
	return nil
}

func AccountDataToErrorMessage(ad *models.AccountData) error {
	if !ad.IsValidUsername() {
		return fmt.Errorf(models.InvalidUsernameInput, ad.Username)
	}
	if !ad.IsValidPassword() {
		return fmt.Errorf(models.PasswordIsTooLongOrEmpty)
	}
	return nil
}
