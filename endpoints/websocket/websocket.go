package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SocketHandler(pool *Pool, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	client := &Client{
		ID:   "fejfjeffef",
		Conn: conn,
		Pool: pool,
	}

	fmt.Println("client pool:", client.ID)
	fmt.Println("client ID:", &client.ID)

	pool.Register <- client
	client.Read()
}
