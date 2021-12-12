package websocket

import (
	"encoding/json"
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
	}

	// Attach the client object as the sender of the messsage.
	message.Sender = client

	switch message.Action {
	case SendMessageAction:
		// The send-message action, this will send messages to a specific room now.
		// Which room wil depend on the message Target
		roomName := message.Target.GetName()
		// Use the ChatServer method to find the room, and if found, broadcastToClient!
		if room := client.server.findRoomByName(roomName); room != nil {
			room.broadcast <- &message
		}
		// We delegate the join and leave actions.
	case JoinRoomAction:
		client.handleJoinRoomMessage(message)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)
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

func (client *Client) joinRoom(roomName string, sender *Client) {

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

func (client *Client) notifyRoomJoined(room *Room, sender *Client) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}

	client.send <- message.encode()
}
