package admin

import (
	"friday/config/repository"
	"friday/config/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsers(c *gin.Context) {
	users, err := repository.GetAllUser()
	utils.FatalError{Error: err}.Handle()

	c.JSON(http.StatusOK, users)
}
