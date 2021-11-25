package auth

import (
	"fmt"
	"friday/models"
	"friday/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := models.GetUserByEmail(json.Email)
	utils.FatalError{Error: err}.Handle()

	fmt.Println(user.Email, user.Password)
	fmt.Println(json.Email, utils.Hash(json.Password), json.Password)

	if user.Email != json.Email || user.Password != utils.Hash(json.Password) {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	token, err := CreateToken(user.Id)
	utils.HttpError{Error: err, Context: c, Status: http.StatusUnprocessableEntity}.Handle()

	c.JSON(http.StatusOK, token)
}

func Register(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := models.CreateUser(json.Email, json.Password, "노과농")
	utils.FatalError{Error: err}.Handle()

	c.JSON(http.StatusOK, json)
}
