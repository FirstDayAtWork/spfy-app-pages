package controllers

import (
	"io"
	"log"
	"net/http"

	"github.com/FirstDayAtWork/mustracker/models"
)

// RequestBodyToAccountData unmarshalls request body to models.AccountData struct.
func RequestBodyToAccountData(r *http.Request) (*models.AccountData, error) {
	bts, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO log this
		return nil, err
	}
	rd := &models.AccountData{}
	if err = rd.Unmarshal(bts); err != nil {
		return nil, err
	}
	if rd.Role == models.EmptyString {
		log.Println("Role not provided, defaulting to:", models.AdminRole)
		rd.Role = models.AdminRole
	}
	return rd, nil
}
