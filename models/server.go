package models

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	InvalidPasswordInputMessage     string = "invalid password provided."
	SuccessMessage                  string = "success."
	InvalidEmailInput               string = "email %s is not valid"
	InvalidUsernameInput            string = "username %s is not valid"
	UsernameAlreadyTakenMessage     string = "username %s is already taken."
	PasswordHashAndPasswordMismatch string = "hashed password does not match the original one. Try Again."
	PasswordIsTooLongOrEmpty        string = "password is too long or empty"
)

// ServerResponse models a generic server response
type ServerResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func (sr *ServerResponse) Unmarshal(bts []byte) error {
	return json.Unmarshal(bts, sr)
}

func (sr *ServerResponse) Marshall() ([]byte, error) {
	return json.Marshal(*sr)
}

func (sr *ServerResponse) SendToClient(w http.ResponseWriter) {
	jsonResp, err := sr.Marshall()
	if err != nil {
		log.Printf("error marshalling server response: %v\n", err)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Printf("error writing response: %v\n", err)
	}
}
