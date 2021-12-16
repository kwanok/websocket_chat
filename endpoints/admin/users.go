package admin

import (
	"friday/server/repository"
	"friday/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsers(c *gin.Context) {
	users, err := repository.GetAllUser()
	utils.FatalError{Error: err}.Handle()

	c.JSON(http.StatusOK, users)
}
