package server

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Type      string `json:"type"`
	User      User   `json:"user"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

type DatabaseInfo struct {
	Name     string
	Host     string
	Password string
	Root     string
}

func initDatabase() *sql.DB {

	databaseInfo := DatabaseInfo{
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ROOT"),
	}

	db, err := sql.Open("mysql", getSourceName(databaseInfo))
	errorHandler(err)

	return db
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getSourceName(db DatabaseInfo) string {
	return db.Root + ":" + db.Password + "@tcp(" + db.Host + ":3306)/" + db.Name
}

func Routes(r *gin.Engine) {
	err := godotenv.Load(".env")
	errorHandler(err)

	db := initDatabase()
	defer db.Close()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "We got Gin")
	})

	admin := r.Group("/admin")
	admin.Use()
	{
		admin.GET("/users", func(c *gin.Context) {
			var users []User

			rows, err := db.Query("SELECT * FROM users")
			errorHandler(err)

			defer rows.Close()

			for rows.Next() {
				var user User
				err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
				if err != nil {
					log.Fatal(err)
				}
				users = append(users, user)
			}

			c.JSON(200, users)
		})
	}

	r.GET("/ws", func(c *gin.Context) {
		socketHandler(c.Writer, c.Request)
	})
}
