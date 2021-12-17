package routes

import (
	"database/sql"
	"friday/config/repository"
	"friday/config/utils"
	"friday/endpoints/admin"
	"friday/endpoints/auth"
	"friday/endpoints/post"
	"friday/endpoints/websocket"
	"friday/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Routes(r *gin.Engine, sqlite *sql.DB) {
	err := godotenv.Load(".env")
	utils.FatalError{Error: err}.Handle()

	wsServer := websocket.NewServer(&repository.RoomRepository{Db: sqlite}, &repository.UserRepository{Db: sqlite})
	go wsServer.Run()

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

	r.GET("posts", post.GetPosts)

	r.GET("/websocket", func(c *gin.Context) {
		websocket.Handler(wsServer, c.Writer, c.Request)
	})
}
