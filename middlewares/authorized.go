package middlewares

import (
	"fmt"
	"friday/config/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsAuthorized(c *gin.Context) {
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"result": "there is invalid",
		})
		c.Abort()
		return
	}

	tokenAuth, err := auth.ExtractTokenMetadata(c.Request)
	fmt.Println("tokenAuth: ", tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	userId, err := auth.FetchAuth(tokenAuth)
	fmt.Println("userId: ", userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	c.Next()
}
