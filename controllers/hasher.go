package controllers

import "golang.org/x/crypto/bcrypt"

const hashCost int = 16

func PasswordToHashedPassword(password string) (string, error) {
	hashedBts, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return "", err
	}
	return string(hashedBts), nil
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
