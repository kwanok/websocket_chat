package auth

import (
	"friday/config/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	accessDetail, err := auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleted, delErr := auth.DeleteAuth(accessDetail.AccessUuid)
	if delErr != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}
