package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	Name   string
	Conn   *websocket.Conn
	Pool   map[*Room]bool
	Send   chan []byte
	Server *Server
}

func newClient(conn *websocket.Conn, server *Server) *Client {
	return &Client{
		Conn:   conn,
		Pool:   make(map[*Room]bool),
		Server: server,
	}
}

func (client *Client) disconnect() {
	client.Server.Unregister <- client
	for room := range client.Pool {
		room.Unregister <- client
	}
}

func (client *Client) GetName() string {
	return client.Name
}
