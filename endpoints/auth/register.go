package auth

import (
	"friday/config"
	"friday/models"
	"friday/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func Register(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	roomRepository := repository.UserRepository{Db: config.Sqlite3}
	roomRepository.AddClient(repository.User{
		Email:    json.Email,
		Level:    models.LevelUser,
		Password: json.Password,
		Name:     json.Name,
	})

	c.JSON(http.StatusOK, json)
}
