package routes

import (
	"fmt"
	"friday/endpoints/admin"
	"friday/endpoints/auth"
	"friday/middlewares"
	"friday/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		utils.FatalError{Error: err}.Handle()

		err = conn.WriteMessage(t, msg)
		utils.FatalError{Error: err}.Handle()
	}
}

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
		socketHandler(c.Writer, c.Request)
	})
}
