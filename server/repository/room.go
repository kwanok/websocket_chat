package repository

import (
	"database/sql"
	"friday/server"
	"friday/server/auth"
	"friday/server/models"
	"friday/server/utils"
)

type Room struct {
	Id      string
	Name    string
	Private bool
}

func (room *Room) GetId() string {
	return room.Id
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
	return room.Private
}

type RoomRepository struct {
	Db *sql.DB
}

func (repo *RoomRepository) AddRoom(room models.Room) {
	stmt, err := repo.Db.Prepare("INSERT INTO room(id, name, private) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(room.GetId(), room.GetName(), room.GetPrivate())
	checkErr(err)
}

func (repo *RoomRepository) FindRoomByName(name string) models.Room {

	row := repo.Db.QueryRow("SELECT id, name, private FROM room where name = ? LIMIT 1", name)

	var room Room

	if err := row.Scan(&room.Id, &room.Name, &room.Private); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return &room

}

func GetAllUser() ([]User, error) {
	users := make([]User, 0)

	rows, err := server.DBCon.Query("SELECT id, level, name, email, created_at, updated_at FROM users")
	utils.FatalError{Error: err}.Handle()
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Level, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		utils.FatalError{Error: err}.Handle()
		users = append(users, user)
	}

	return users, nil
}

func GetUserByEmail(email string) (User, error) {
	var user User

	query := "SELECT id, level, name, email, password, created_at, updated_at FROM users WHERE email = ? LIMIT 1"

	err := server.DBCon.QueryRow(query, email).Scan(&user.Id, &user.Level, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	utils.FatalError{Error: err}.Handle()

	return user, nil
}

func GetUserById(id string) (User, error) {
	var user User

	query := "SELECT id, level, name, email, password, created_at, updated_at FROM users WHERE id = ? LIMIT 1"

	err := server.DBCon.QueryRow(query, id).Scan(&user.Id, &user.Level, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	utils.FatalError{Error: err}.Handle()

	return user, nil
}

func CreateUser(email string, level int, password string, name string) error {
	tx, err := server.DBCon.Begin()
	utils.PanicError{Error: err}.Handle()
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO users (name, level, email, password) VALUES (?, ?, ?, ?)", name, level, email, auth.Hash(password))
	utils.PanicError{Error: err}.Handle()

	err = tx.Commit()
	utils.PanicError{Error: err}.Handle()

	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
