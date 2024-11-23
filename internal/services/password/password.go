package password

import (
	"github.com/albakov/go-weather-viewer/internal/util"
	"golang.org/x/crypto/bcrypt"
)

const f = "auth"

func CreateHashedPassword(password string) (string, error) {
	const op = "CreateHashedPassword"

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		util.LogError(f, op, err)

		return "", err
	}

	return string(hashed), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}
