package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/kwanok/friday/config"
	"github.com/kwanok/friday/repository"
	"log"
	"net/http"
)

func GetUsers(c *gin.Context) {
	userRepository := repository.UserRepository{Db: config.DBCon}
	users := userRepository.GetAllClients()

	log.Println("유저: ", users)

	c.JSON(http.StatusOK, users)
}
