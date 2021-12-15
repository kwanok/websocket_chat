package main

import (
	"database/sql"
	"friday/routes"
	"friday/server"
	"friday/server/utils"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var mainLogger *log.Logger

func main() {
	r := gin.Default()
	mainLogger = log.New(os.Stdout, "MAIN: ", log.LstdFlags)
	mainLogger.Println(server.InitDB())

	defer func(DBCon *sql.DB) {
		err := DBCon.Close()
		utils.FatalError{Error: err}.Handle()

	}(server.DBCon)

	defer func(Sqlite3 *sql.DB) {
		err := Sqlite3.Close()
		utils.FatalError{Error: err}.Handle()

	}(server.Sqlite3)

	routes.Routes(r)

	err := r.Run()
	utils.FatalError{Error: err}.Handle()
}
