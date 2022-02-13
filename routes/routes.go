package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	authConfig "github.com/kwanok/friday/config/auth"
	"github.com/kwanok/friday/endpoints/admin"
	"github.com/kwanok/friday/endpoints/auth"
	"github.com/kwanok/friday/endpoints/websocket"
	"github.com/kwanok/friday/middlewares"
	"github.com/kwanok/friday/repository"
	"log"
	"net/http"
	"os"
)

func Routes(r *gin.Engine, db *sql.DB) {
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}

	wsServer := websocket.NewServer(&repository.RoomRepository{Db: db}, &repository.UserRepository{Db: db})
	go wsServer.Run()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "We got Gin")
	})

	r.GET("/users", admin.GetUsers)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
		authGroup.POST("/logout", middlewares.IsAuthorized, auth.Logout)
		authGroup.POST("/refresh", auth.Refresh)
	}

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
