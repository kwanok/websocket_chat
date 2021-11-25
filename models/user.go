package models

import (
	"friday/utils"
)

type User struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func GetAllUser() ([]User, error) {
	users := make([]User, 0)

	rows, err := DBCon.Query("SELECT id, name, email, created_at, updated_at FROM users")
	utils.FatalError{Error: err}.Handle()
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		utils.FatalError{Error: err}.Handle()
		users = append(users, user)
	}

	return users, nil
}

func GetUserByEmail(email string) (User, error) {
	var user User

	query := "SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ?"

	err := DBCon.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	utils.FatalError{Error: err}.Handle()

	return user, nil
}

func CreateUser(email string, password string, name string) error {
	tx, err := DBCon.Begin()
	utils.PanicError{Error: err}.Handle()
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, utils.Hash(password))
	utils.PanicError{Error: err}.Handle()

	err = tx.Commit()
	utils.PanicError{Error: err}.Handle()

	return nil
}
