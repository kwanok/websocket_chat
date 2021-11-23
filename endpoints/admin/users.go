package admin

import (
	"database/sql"
	"friday/server"
	"friday/tools"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var db *sql.DB

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func GetUsers(c *gin.Context) {
	db = server.DBCon

	rows, err := db.Query("SELECT * FROM users")
	tools.ErrorHandler(err)
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}
