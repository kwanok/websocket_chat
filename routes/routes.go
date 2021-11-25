package routes

import (
	"fmt"
	"friday/endpoints/admin"
	"friday/endpoints/auth"
	"friday/utils"
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
	adminGroup.Use()
	{
		adminGroup.GET("/users", func(c *gin.Context) { admin.GetUsers(c) })
	}

	authGroup := r.Group("/auth")
	{
		authGroup.GET("/login", func(c *gin.Context) { auth.Login(c) })
		authGroup.POST("/register", func(c *gin.Context) { auth.Register(c) })
	}

	r.GET("/ws", func(c *gin.Context) {
		socketHandler(c.Writer, c.Request)
	})
}
