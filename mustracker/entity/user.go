package entity

import (
	"encoding/json"
)

type AccountData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ad *AccountData) Unmarshal(bts []byte) error {
	return json.Unmarshal(bts, ad)
}
