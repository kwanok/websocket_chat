package repository

import (
	"database/sql"
	"friday/config/auth"
	"friday/models"
	"github.com/google/uuid"
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

func (user *User) GetPassword() string {
	return user.Password
}

type UserRepository struct {
	Db *sql.DB
}

func (repo *UserRepository) AddUser(user models.ChatClient) {
	stmt, err := repo.Db.Prepare("INSERT INTO users(id, name) values(?,?)")
	checkErr(err)

	_, err = stmt.Exec(user.GetId(), user.GetName())
	checkErr(err)
}

func (repo *UserRepository) RemoveUser(user models.ChatClient) {
	stmt, err := repo.Db.Prepare("DELETE FROM users WHERE id = ?")
	checkErr(err)

	_, err = stmt.Exec(user.GetId())
	checkErr(err)
}

func (repo *UserRepository) FindChatClientById(ID string) models.ChatClient {

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

func (repo *UserRepository) FindClientById(Id string) models.Client {
	row := repo.Db.QueryRow("SELECT id, level, email, password, name FROM users where id = ? LIMIT 1", Id)

	var user User

	if err := row.Scan(
		&user.Id,
		&user.Level,
		&user.Email,
		&user.Password,
		&user.Name,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return &user
}

func (repo *UserRepository) FindClientByEmail(email string) models.Client {
	row := repo.Db.QueryRow("SELECT id, level, email, password, name FROM users where email = ? LIMIT 1", email)
	var user User

	if err := row.Scan(
		&user.Id,
		&user.Level,
		&user.Email,
		&user.Password,
		&user.Name,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return &user
}

func (repo *UserRepository) GetAllUsers() []models.ChatClient {

	rows, err := repo.Db.Query("SELECT id, name FROM users")

	if err != nil {
		log.Fatal(err)
	}
	var users []models.ChatClient
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name)
		users = append(users, &user)
	}

	return users
}

func (repo *UserRepository) AddClient(user User) {
	stmt, err := repo.Db.Prepare("INSERT INTO users(id, name, level, email, password) values(?,?,?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(uuid.New(), user.GetName(), user.GetLevel(), user.GetEmail(), auth.Hash(user.GetPassword()))
	checkErr(err)
}

func (repo *UserRepository) GetAllClients() []models.Client {
	rows, err := repo.Db.Query("SELECT id, level, email, password, name FROM users")

	if err != nil {
		log.Fatal(err)
	}
	var users []models.Client
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(
			&user.Id,
			&user.Level,
			&user.Email,
			&user.Password,
			&user.Name,
		)
		users = append(users, &user)
	}

	return users
}
