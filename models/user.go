package models

import (
	"encoding/json"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type AccountData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ad *AccountData) Unmarshal(bts []byte) error {
	return json.Unmarshal(bts, ad)
}

func (ad *AccountData) IsValidUsername() bool {
	return ad.Username != emptyString
}

func (ad *AccountData) IsValidEmail() bool {
	return regexp.MustCompile(emailRegex).MatchString(strings.ToLower(ad.Email))
}

func (ad *AccountData) IsValidPassword() bool {
	if len([]byte(ad.Password)) > maxPasswordLen {
		return false
	}
	return ad.Password != emptyString
}

func MigrateAccountData(db *gorm.DB) error {
	return db.AutoMigrate(&AccountData{})
}
