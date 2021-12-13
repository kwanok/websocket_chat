package websocket

import (
	"fmt"
	"github.com/google/uuid"
	"log"
)

const welcomeMessage = "%s joined the room"

type Room struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	Private    bool `json:"private"`
}

func NewRoom(name string, private bool) *Room {
	return &Room{
		Id:         uuid.New(),
		Name:       name,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		Private:    private,
	}
}

func (room *Room) Start() {
	for {
		select {

		case client := <-room.register:
			room.registerClient(client)

		case client := <-room.unregister:
			room.unregisterClient(client)

		case message := <-room.broadcast:
			room.broadcastToClients(message.encode())
		}
	}
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetId() string {
	return room.Id.String()
}

//----------------------------------------------------------------------------------------------------------------------

//registerClient 클라이언트를 채팅방에 등록
func (room *Room) registerClient(client *Client) {
	log.Println(room.GetId())
	room.notifyClientJoined(client)
	room.clients[client] = true
}

//unregisterClient 클라이언트를 채팅방에 등록 해제
func (room *Room) unregisterClient(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

//broadcastToClients 채팅방의 클라이언트들에게 브로드캐스팅
func (room *Room) broadcastToClients(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}

	room.broadcastToClients(message.encode())
}
