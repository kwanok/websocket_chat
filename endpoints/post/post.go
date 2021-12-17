package post

import (
	"friday/config/models"
	"friday/config/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPosts(c *gin.Context) {
	posts, err := models.GetAllPosts()
	utils.FatalError{Error: err}.Handle()

	c.JSON(http.StatusOK, posts)
}
