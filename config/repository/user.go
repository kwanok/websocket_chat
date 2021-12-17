package repository

import (
	"database/sql"
	"friday/config"
	"friday/config/auth"
	"friday/config/models"
	"friday/config/utils"
	"log"
)

type User struct {
	Id        string `json:"id"`
	Level     int    `json:"level"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (user *User) GetId() string {
	return user.Id
}

func (user *User) GetName() string {
	return user.Name
}

func (user *User) GetLevel() int {
	return user.Level
}

func (user *User) GetEmail() string {
	return user.Email
}

type UserRepository struct {
	Db *sql.DB
}

func (repo *UserRepository) AddUser(user models.User) {
	stmt, err := repo.Db.Prepare("INSERT INTO users(id, name) values(?,?)")
	checkErr(err)

	_, err = stmt.Exec(user.GetId(), user.GetName())
	checkErr(err)
}

func (repo *UserRepository) RemoveUser(user models.User) {
	stmt, err := repo.Db.Prepare("DELETE FROM users WHERE id = ?")
	checkErr(err)

	_, err = stmt.Exec(user.GetId())
	checkErr(err)
}

func (repo *UserRepository) FindUserById(ID string) models.User {

	row := repo.Db.QueryRow("SELECT id, name FROM users where id = ? LIMIT 1", ID)

	var user User

	if err := row.Scan(&user.Id, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return &user

}

func (repo *UserRepository) GetAllUsers() []models.User {

	rows, err := repo.Db.Query("SELECT id, name FROM users")

	if err != nil {
		log.Fatal(err)
	}
	var users []models.User
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name)
		users = append(users, &user)
	}

	return users
}

func GetAllUser() ([]User, error) {
	users := make([]User, 0)

	rows, err := config.DBCon.Query("SELECT id, level, name, email, created_at, updated_at FROM users")
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

	err := config.DBCon.QueryRow(query, email).Scan(&user.Id, &user.Level, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	utils.FatalError{Error: err}.Handle()

	return user, nil
}

func GetUserById(id string) (User, error) {
	var user User

	query := "SELECT id, level, name, email, password, created_at, updated_at FROM users WHERE id = ? LIMIT 1"

	err := config.DBCon.QueryRow(query, id).Scan(&user.Id, &user.Level, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	utils.FatalError{Error: err}.Handle()

	return user, nil
}

func CreateUser(email string, level int, password string, name string) error {
	tx, err := config.DBCon.Begin()
	utils.PanicError{Error: err}.Handle()
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO users (name, level, email, password) VALUES (?, ?, ?, ?)", name, level, email, auth.Hash(password))
	utils.PanicError{Error: err}.Handle()

	err = tx.Commit()
	utils.PanicError{Error: err}.Handle()

	return nil
}
