package auth

import (
	"friday/config"
	"friday/config/auth"
	"friday/config/utils"
	"friday/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Login(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	userRepository := repository.UserRepository{Db: config.Sqlite3}
	user := userRepository.FindClientByEmail(json.Email)
	log.Println(user)
	if user == nil {
		c.JSON(http.StatusNotFound, "Not found")
		c.Abort()
		return
	}

	if user.GetEmail() != json.Email || !auth.CompareHash(user.GetPassword(), json.Password) {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	token, err := auth.CreateToken(user.GetId(), user.GetLevel())
	utils.HttpError{Error: err, Context: c, Status: http.StatusUnprocessableEntity}.Handle()

	saveErr := auth.CreateAuth(user.GetId(), token)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
		return
	}

	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
}
