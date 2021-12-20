package websocket

import (
	"encoding/json"
	"friday/models"
	"log"
)

//SendMessageAction 메시지 보낼때
const SendMessageAction = "send-message"

//JoinRoomAction 방 참가할때
const JoinRoomAction = "join-room"

//LeaveRoomAction 방 나갈때
const LeaveRoomAction = "leave-room"

//UserJoinedAction 유저가 참가할 때
const UserJoinedAction = "user-join"

//UserLeftAction 유저 나갈때
const UserLeftAction = "user-left"

const JoinRoomPrivateAction = "join-room-private"
const RoomJoinedAction = "room-joined"

type Message struct {
	Action  string            `json:"action"`
	Message string            `json:"message"`
	Target  *Room             `json:"target"`
	Sender  models.ChatClient `json:"sender"`
}

//encode 제이슨을 평문으로 마샬링
func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}

func (message *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	msg := &struct {
		Sender Client `json:"sender"`
		*Alias
	}{
		Alias: (*Alias)(message),
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	message.Sender = &msg.Sender
	return nil
}
