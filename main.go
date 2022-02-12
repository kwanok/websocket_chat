package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/kwanok/friday/config"
	"github.com/kwanok/friday/config/utils"
	"github.com/kwanok/friday/routes"
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
