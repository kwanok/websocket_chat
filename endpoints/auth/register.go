package auth

import (
	"friday/config/utils"
	"friday/models"
	"friday/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := repository.CreateUser(json.Email, models.LevelUser, json.Password, "노과농")
	utils.FatalError{Error: err}.Handle()

	c.JSON(http.StatusOK, json)
}
