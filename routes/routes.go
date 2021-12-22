package routes

import (
	"database/sql"
	authConfig "friday/config/auth"
	"friday/config/utils"
	"friday/endpoints/admin"
	"friday/endpoints/auth"
	"friday/endpoints/post"
	"friday/endpoints/websocket"
	"friday/middlewares"
	"friday/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
)

func Routes(r *gin.Engine, sqlite *sql.DB) {
	err := godotenv.Load(".env")
	utils.FatalError{Error: err}.Handle()

	wsServer := websocket.NewServer(&repository.RoomRepository{Db: sqlite}, &repository.UserRepository{Db: sqlite})
	go wsServer.Run()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "We got Gin")
	})

	r.GET("/users", admin.GetUsers)

	adminGroup := r.Group("/admin")
	adminGroup.Use(middlewares.IsAuthorized)
	{
		adminGroup.GET("/users", admin.GetUsers)
		adminGroup.GET("/websockets", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"ws": wsServer,
			})
		})
	}

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
		authGroup.POST("/logout", middlewares.IsAuthorized, auth.Logout)
		authGroup.POST("/refresh", auth.Refresh)
	}

	r.GET("posts", post.GetPosts)

	r.GET("/websocket", middlewares.IsAuthorized, func(c *gin.Context) {
		accessDetail, err := authConfig.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}

		userId, err := authConfig.FetchAuth(accessDetail)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}

		websocket.Handler(wsServer, c.Writer, c.Request, userId)
	})

}
