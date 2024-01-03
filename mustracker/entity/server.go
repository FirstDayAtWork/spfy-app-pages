package entity

import "encoding/json"

type SimpleJson struct {
	Username string
	Email    string
	Password string
}

type RegistrationData struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	HashedPassword string
}

func (rd *RegistrationData) Unmarshal(bts []byte) error {
	return json.Unmarshal(bts, rd)
}
