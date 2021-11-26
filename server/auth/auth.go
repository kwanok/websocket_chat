package auth

import (
	"golang.org/x/crypto/bcrypt"
)

//Hash 문자열을 해싱해서 리턴
func Hash(str string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)

	return string(hash)
}

//CompareHash 문자열을 검증
func CompareHash(pwHash string, pwString string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwHash), []byte(pwString))
	if err != nil {
		return false
	} else {
		return true
	}
}
