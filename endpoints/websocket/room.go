package websocket

import "fmt"

const welcomeMessage = "%s joined the room"

type Room struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func New() *Room {
	return &Room{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (room *Room) Start() {
	for {
		select {

		case client := <-room.Register:
			room.register(client)

		case client := <-room.Unregister:
			room.unregister(client)

		case message := <-room.Broadcast:
			room.broadcast(message.encode())
		}
	}
}

func (room *Room) register(client *Client) {
	room.notifyClientJoined(client)
	room.Clients[client] = true
}

func (room *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}

	room.broadcastToClientsInRoom(message.encode())
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.Clients {
		client.Send <- message
	}
}

func (room *Room) unregister(client *Client) {
	if _, ok := room.Clients[client]; ok {
		delete(room.Clients, client)
	}
}

func (room *Room) broadcast(message []byte) {
	for client := range room.Clients {
		client.Send <- message
	}
}
