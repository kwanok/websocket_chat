package websocket

import (
	"encoding/json"
	"friday/server/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	conn   *websocket.Conn
	pool   map[*Room]bool
	send   chan []byte
	server *Server
}

func newClient(conn *websocket.Conn, server *Server, name string) *Client {
	return &Client{
		pool:   make(map[*Room]bool),
		server: server,
		ID:     uuid.New(),
		Name:   name,
		conn:   conn,
		send:   make(chan []byte, 256),
	}
}

func (client *Client) GetId() string {
	return client.ID.String()
}

func (client *Client) GetName() string {
	return client.Name
}

//----------------------------------------------------------------------------------------------------------------------

func (client *Client) disconnect() {
	client.server.unregister <- client
	for room := range client.pool {
		room.unregister <- client
	}
}

func (client *Client) handleNewMessage(jsonMessage []byte) {

	var message Message

	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("response: %q", jsonMessage)
	}

	// Attach the client object as the sender of the messsage.
	message.Sender = client

	switch message.Action {
	case SendMessageAction:
		roomID := message.Target.GetId()
		if room := client.server.findRoomByID(roomID); room != nil {
			room.broadcast <- &message
		}

	case JoinRoomAction:
		client.handleJoinRoomMessage(message)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)

	case JoinRoomPrivateAction:
		client.handleJoinRoomPrivateMessage(message)
	}
}

func (client *Client) handleJoinRoomMessage(message Message) {
	roomName := message.Message

	client.joinRoom(roomName, nil)
}

func (client *Client) handleLeaveRoomMessage(message Message) {
	room := client.server.findRoomByID(message.Message)
	if room == nil {
		return
	}

	if _, ok := client.pool[room]; ok {
		delete(client.pool, room)
	}

	room.unregister <- client
}

func (client *Client) handleJoinRoomPrivateMessage(message Message) {

	target := client.server.findClientByID(message.Message)

	if target == nil {
		return
	}

	// create unique room name combined to the two IDs
	roomName := message.Message + client.ID.String()

	client.joinRoom(roomName, target)
	target.joinRoom(roomName, client)

}

func (client *Client) joinRoom(roomName string, sender models.User) {

	room := client.server.findRoomByName(roomName)
	if room == nil {
		room = client.server.createRoom(roomName, sender != nil)
	}

	// Don't allow to join private rooms through public room message
	if sender == nil && room.Private {
		return
	}

	if !client.isInRoom(room) {

		client.pool[room] = true
		room.register <- client

		client.notifyRoomJoined(room, sender)
	}

}

func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.pool[room]; ok {
		return true
	}

	return false
}

func (client *Client) notifyRoomJoined(room *Room, sender models.User) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}

	client.send <- message.encode()
}
