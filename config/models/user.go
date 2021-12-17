package models

const (
	LevelAdmin = 0
	LevelUser  = 1
)

type User interface {
	GetId() string
	GetName() string
}

type UserRepository interface {
	AddUser(user User)
	RemoveUser(user User)
	FindUserById(ID string) User
	GetAllUsers() []User
}
