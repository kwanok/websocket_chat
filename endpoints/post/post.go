package post

import (
	"github.com/gin-gonic/gin"
	"github.com/kwanok/friday/config/utils"
	"github.com/kwanok/friday/models"
	"net/http"
)

func GetPosts(c *gin.Context) {
	posts, err := models.GetAllPosts()
	utils.FatalError{Error: err}.Handle()

	c.JSON(http.StatusOK, posts)
}
