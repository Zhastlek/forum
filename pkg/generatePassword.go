package pkg

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string) ([]byte, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// log.Println("Error user handler sign-up generate password")
		return nil, err
	}
	return pass, nil
}
