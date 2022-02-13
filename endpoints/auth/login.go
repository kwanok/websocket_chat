package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/kwanok/friday/config"
	"github.com/kwanok/friday/config/auth"
	"github.com/kwanok/friday/repository"
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

	userRepository := repository.UserRepository{Db: config.DBCon}
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
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "StatusUnprocessableEntity")
		return
	}

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
