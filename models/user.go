package models

const (
	LevelAdmin = 0
	LevelUser  = 1
)

type ChatClient interface {
	GetId() string
	GetName() string
}

type Client interface {
	GetId() string
	GetName() string
	GetLevel() int
	GetPassword() string
	GetEmail() string
}

type UserRepository interface {
	AddUser(user ChatClient)
	RemoveUser(user ChatClient)
	FindChatClientById(ID string) ChatClient
	FindClientById(Id string) Client
	GetAllUsers() []ChatClient
}
