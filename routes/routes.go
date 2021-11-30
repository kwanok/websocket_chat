package routes

import (
	"friday/endpoints/admin"
	"friday/endpoints/auth"
	"friday/endpoints/ws"
	"friday/middlewares"
	"friday/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Routes(r *gin.Engine) {
	err := godotenv.Load(".env")
	utils.FatalError{Error: err}.Handle()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "We got Gin")
	})

	adminGroup := r.Group("/admin")
	adminGroup.Use(middlewares.IsAuthorized)
	{
		adminGroup.GET("/users", admin.GetUsers)
	}

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
		authGroup.POST("/logout", middlewares.IsAuthorized, auth.Logout)
		authGroup.POST("/refresh", auth.Refresh)
	}

	r.GET("/ws", func(c *gin.Context) {
		ws.SocketHandler(c.Writer, c.Request)
	})
}
