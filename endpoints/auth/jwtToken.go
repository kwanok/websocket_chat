package auth

import (
	"friday/utils"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

func CreateToken(id uint64) (string, error) {
	err := os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	utils.FatalError{Error: err}.Handle()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = id
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	utils.FatalError{Error: err}.Handle()

	return token, nil
}
