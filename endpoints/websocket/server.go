package websocket

import (
	"log"
	"net/http"
)

type Server struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func ServeWs(wsServer *Server, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer)

	go client.writePump()
	go client.readPump()

	wsServer.Register <- client
}

func NewWebsocketServer() *Server {
	return &Server{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (server *Server) Run() {
	for {
		select {

		case client := <-server.Register:
			server.registerClient(client)

		case client := <-server.Unregister:
			server.unregisterClient(client)
		}
	}
}

func (server *Server) broadcastToClients(message []byte) {
	for client := range server.Clients {
		client.Send <- message
	}
}

func (server *Server) registerClient(client *Client) {
	server.Clients[client] = true
}

func (server *Server) unregisterClient(client *Client) {
	if _, ok := server.Clients[client]; ok {
		delete(server.Clients, client)
	}
}

func (server *Server) notifyClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}

	server.broadcastToClients(message.encode())
}

func (server *Server) notifyClientLeft(client *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}

	server.broadcastToClients(message.encode())
}

func (server *Server) listOnlineClients(client *Client) {
	for existingClient := range server.Clients {
		message := &Message{
			Action: UserJoinedAction,
			Sender: existingClient,
		}
		client.Send <- message.encode()
	}
}
