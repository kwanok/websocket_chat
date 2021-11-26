package models

import (
	"friday/server"
	"friday/server/auth"
	"friday/server/utils"
)

const (
	LevelAdmin = 0
	LevelUser  = 1
)

type User struct {
	Id        uint64 `json:"id"`
	Level     int    `json:"level"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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

func GetUserById(id uint64) (User, error) {
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
