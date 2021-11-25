package utils

import "golang.org/x/crypto/bcrypt"

func Hash(str string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)

	return string(hash)
}
