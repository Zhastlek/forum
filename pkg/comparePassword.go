package pkg

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func ComparePassword(userPasswordInBD []byte, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPasswordInBD), []byte(str))
	if err != nil {
		log.Println(("Error user handler sign-in password is wrong"))
		log.Printf("%v\n", err)
		return false
	}
	return true
}
