package auth

import (
	"crypto/rand"
	"encoding/base64"
)

type TokenClaims struct {
	// jwt.StandardClaims
	Role string `json:"role"`
	CSRF string `json:"csrf"`
}

func GenerateCSFRSecret() (string, error) {
	return generateRandomString(32)
}

func generateRandomString(n int) (string, error) {
	b, err := generateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}
