package websocket

import (
	"encoding/json"
	"friday/config"
	"friday/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	conn   *websocket.Conn
	server *Server
	rooms  map[*Room]bool
	send   chan []byte
}

func newClient(conn *websocket.Conn, server *Server, id string) *Client {
	userId, err := uuid.Parse(id)
	if err != nil {
		log.Fatal(err)
	}

	user := server.userRepository.FindClientById(id)

	return &Client{
		ID:     userId,
		Name:   user.GetName(),
		conn:   conn,
		server: server,
		rooms:  make(map[*Room]bool),
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
	for room := range client.rooms {
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

	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
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

	joinedRoom := client.joinRoom(roomName, target)

	if joinedRoom != nil {
		client.inviteTargetUser(target, joinedRoom)
	}
}

func (client *Client) joinRoom(roomName string, sender models.ChatClient) *Room {

	room := client.server.findRoomByName(roomName)
	if room == nil {
		room = client.server.createRoom(roomName, sender != nil)
	}

	// Don't allow to join private rooms through public room message
	if sender == nil && room.Private {
		return nil
	}

	if !client.isInRoom(room) {
		client.rooms[room] = true
		room.register <- client
		client.notifyRoomJoined(room, sender)
	}

	return room
}

func (client *Client) inviteTargetUser(target models.ChatClient, room *Room) {
	inviteMessage := &Message{
		Action:  JoinRoomPrivateAction,
		Message: target.GetId(),
		Target:  room,
		Sender:  client,
	}

	if err := config.PubSubRedis.Publish(ctx, PubSubGeneralChannel, inviteMessage.encode()).Err(); err != nil {
		log.Println(err)
	}
}

func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}

	return false
}

func (client *Client) notifyRoomJoined(room *Room, sender models.ChatClient) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}

	client.send <- message.encode()
}
