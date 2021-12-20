package admin

import (
	"friday/config"
	"friday/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetUsers(c *gin.Context) {
	userRepository := repository.UserRepository{Db: config.Sqlite3}
	users := userRepository.GetAllClients()

	log.Println("유저: ", users)

	c.JSON(http.StatusOK, users)
}
