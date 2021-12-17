package auth

import (
	"friday/config/auth"
	"friday/config/repository"
	"friday/config/utils"
	"github.com/gin-gonic/gin"
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

	user, err := repository.GetUserByEmail(json.Email)
	utils.FatalError{Error: err}.Handle()

	if user.Email != json.Email || !auth.CompareHash(user.Password, json.Password) {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	token, err := auth.CreateToken(user.Id, user.Level)
	utils.HttpError{Error: err, Context: c, Status: http.StatusUnprocessableEntity}.Handle()

	saveErr := auth.CreateAuth(user.Id, token)
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
