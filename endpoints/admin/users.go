package admin

import (
	"friday/server/models"
	"friday/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsers(c *gin.Context) {
	users, err := models.GetAllUser()
	utils.FatalError{Error: err}.Handle()

	c.JSON(http.StatusOK, users)
}
