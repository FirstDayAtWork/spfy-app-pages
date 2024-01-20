package controllers

import (
	"io"
	"net/http"

	"github.com/FirstDayAtWork/mustracker/models"
)

// RegistrationRequestToAccountData unmarshalls request body to models.AccountData struct.
func RegistrationRequestToAccountData(r *http.Request) (*models.AccountData, error) {
	bts, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO log this
		return nil, err
	}
	rd := &models.AccountData{}
	if err = rd.Unmarshal(bts); err != nil {
		return nil, err
	}

	return rd, nil
}
