package websocket

import (
	"encoding/json"
	"friday/config"
	models "friday/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const PubSubGeneralChannel = "general"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	clients        map[*Client]bool
	rooms          map[*Room]bool
	register       chan *Client
	unregister     chan *Client
	broadcast      chan []byte
	users          []models.ChatClient
	roomRepository models.RoomRepository
	userRepository models.UserRepository
}

func NewServer(
	roomRepository models.RoomRepository,
	userRepository models.UserRepository,
) *Server {
	server := &Server{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte),
		rooms:          make(map[*Room]bool),
		roomRepository: roomRepository,
		userRepository: userRepository,
	}

	server.users = userRepository.GetAllUsers()

	return server
}

func Handler(server *Server, w http.ResponseWriter, r *http.Request) {
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

	client := newClient(conn, server, name[0])

	go client.writePump()
	go client.readPump()

	server.register <- client
}

func (server *Server) Run() {
	go server.listenPubSubChannel()
	for {
		select {
		case client := <-server.register:
			server.registerClient(client)
		case client := <-server.unregister:
			server.unregisterClient(client)
		case message := <-server.broadcast:
			server.broadcastToClients(message)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------

//registerClient 서버에 클라이언트를 등록
func (server *Server) registerClient(client *Client) {
	// 유저를 db에 저장
	server.userRepository.AddUser(client)

	server.publishClientJoined(client)

	server.listOnlineClients(client)
	server.clients[client] = true
}

//unregisterClient 서버에서 클라이언트를 등록 해제
func (server *Server) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)

		// Remove user from repo
		server.userRepository.RemoveUser(client)

		// Publish user left in PubSub
		server.publishClientLeft(client)
	}
}

//broadcastToClients 서버에 있는 클라이언트들에게 브로드캐스팅
func (server *Server) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
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
	for _, user := range server.users {
		message := &Message{
			Action: UserJoinedAction,
			Sender: user,
		}
		client.send <- message.encode()
	}
}

func (server *Server) runRoomFromRepository(name string) *Room {
	var room *Room
	dbRoom := server.roomRepository.FindRoomByName(name)
	if dbRoom != nil {
		room = NewRoom(dbRoom.GetName(), dbRoom.GetPrivate())
		room.Id, _ = uuid.Parse(dbRoom.GetId())

		go room.Start()
		server.rooms[room] = true
	}

	return room
}

//findRoomByName 이름으로 채팅방 찾기
func (server *Server) findRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetName() == name {
			foundRoom = room
			break
		}
	}

	if foundRoom == nil {
		foundRoom = server.runRoomFromRepository(name)
	}

	return foundRoom
}

//findRoomByID ID로 채팅방 찾기
func (server *Server) findRoomByID(ID string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetId() == ID {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

//findClientByID ID로 클라이언트 찾기
func (server *Server) findClientByID(ID string) *Client {
	var foundClient *Client
	for client := range server.clients {
		if client.ID.String() == ID {
			foundClient = client
			break
		}
	}

	return foundClient
}

//findUserByID ID로 유저 찾기
func (server *Server) findUserByID(ID string) models.ChatClient {
	var foundUser models.ChatClient
	for _, client := range server.users {
		if client.GetId() == ID {
			foundUser = client
			break
		}
	}

	return foundUser
}

//createRoom 채팅방 생성
func (server *Server) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private)
	server.roomRepository.AddRoom(room)

	go room.Start()
	server.rooms[room] = true

	return room
}

//publishClientJoined Redis 유저 합류 PUB
func (server *Server) publishClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}

	if err := config.PubSubRedis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println(err)
	}
}

//publishClientLeft Redis 유저 이탈 PUB
func (server *Server) publishClientLeft(client *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}

	if err := config.PubSubRedis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println(err)
	}
}

func (server *Server) listenPubSubChannel() {

	pubsub := config.PubSubRedis.Subscribe(ctx, PubSubGeneralChannel)
	ch := pubsub.Channel()
	for msg := range ch {

		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("Error on unmarshal JSON message %s", err)
			return
		}

		switch message.Action {
		case UserJoinedAction:
			server.handleUserJoined(message)
		case UserLeftAction:
			server.handleUserLeft(message)
		case JoinRoomPrivateAction:
			server.handleUserJoinPrivate(message)
		}
	}
}

func (server *Server) handleUserJoined(message Message) {
	// Add the user to the slice
	server.users = append(server.users, message.Sender)
	server.broadcastToClients(message.encode())
}

func (server *Server) handleUserLeft(message Message) {
	// Remove the user from the slice
	for i, user := range server.users {
		if user.GetId() == message.Sender.GetId() {
			server.users[i] = server.users[len(server.users)-1]
			server.users = server.users[:len(server.users)-1]
		}
	}
	server.broadcastToClients(message.encode())
}

func (server *Server) handleUserJoinPrivate(message Message) {
	// Find client for given user, if found add the user to the room.
	targetClient := server.findClientByID(message.Message)
	if targetClient != nil {
		targetClient.joinRoom(message.Target.GetName(), message.Sender)
	}
}
