package utils

import "golang.org/x/crypto/bcrypt"

const DefaultPassword = "p"

func GenerateDefaultHashedPassword() string {
	password, err := HashPassword(DefaultPassword)
	if err != nil {
		return ""
	}
	return password
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
