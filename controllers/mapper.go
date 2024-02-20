package controllers

import (
	"encoding/json"
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
