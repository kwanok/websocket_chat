package main

import (
	"database/sql"
	"friday/config"
	"friday/config/utils"
	"friday/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var mainLogger *log.Logger

func main() {
	r := gin.Default()
	mainLogger = log.New(os.Stdout, "MAIN: ", log.LstdFlags)
	mainLogger.Println(config.InitDB())

	defer func(DBCon *sql.DB) {
		err := DBCon.Close()
		utils.FatalError{Error: err}.Handle()

	}(config.DBCon)

	sqlite := config.InitSqlite3()
	defer sqlite.Close()

	routes.Routes(r, sqlite)

	err := r.Run()
	utils.FatalError{Error: err}.Handle()
}
