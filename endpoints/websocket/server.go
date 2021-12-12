package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Server struct {
	Clients    map[*Client]bool
	Rooms      map[*Room]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Handler(wsServer *Server, w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]
	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer, name[0])

	go client.writePump()
	go client.readPump()

	wsServer.register <- client
}

func NewServer() *Server {
	return &Server{
		Clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (server *Server) Run() {
	for {
		select {
		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------

//registerClient 서버에 클라이언트를 등록
func (server *Server) registerClient(client *Client) {
	server.Clients[client] = true
}

//unregisterClient 서버에서 클라이언트를 등록 해제
func (server *Server) unregisterClient(client *Client) {
	if _, ok := server.Clients[client]; ok {
		delete(server.Clients, client)
	}
}

//broadcastToClient 서버에 있는 클라이언트들에게 브로드캐스팅
func (server *Server) broadcastToClient(message []byte) {
	for client := range server.Clients {
		client.send <- message
	}
}

func (server *Server) notifyClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}

	server.broadcastToClient(message.encode())
}

func (server *Server) notifyClientLeft(client *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}

	server.broadcastToClient(message.encode())
}

func (server *Server) listOnlineClients(client *Client) {
	for existingClient := range server.Clients {
		message := &Message{
			Action: UserJoinedAction,
			Sender: existingClient,
		}
		client.send <- message.encode()
	}
}

//findRoomByName 이름으로 채팅방 찾기
func (server *Server) findRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range server.Rooms {
		if room.GetName() == name {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

//findRoomByID ID로 채팅방 찾기
func (server *Server) findRoomByID(ID string) *Room {
	var foundRoom *Room
	for room := range server.Rooms {
		if room.GetId() == ID {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

//createRoom 채팅방을 생성
func (server *Server) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private)
	go room.Start()
	server.Rooms[room] = true

	return room
}
