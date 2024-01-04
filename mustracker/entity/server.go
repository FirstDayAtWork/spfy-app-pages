package entity

import "encoding/json"

const (
	InvalidPasswordInputMessage     string = "Invalid password provided."
	SuccessMessage                  string = "Success."
	UsernameAlreadyTakenMessage     string = "Username %s is already taken."
	PasswordHashAndPasswordMismatch string = "Hashed password does not match the original one. Try Again."
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
