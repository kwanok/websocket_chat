package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/kwanok/friday/config"
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
		if err != nil {
			log.Fatal(err)
		}

	}(config.DBCon)

	routes.Routes(r, config.DBCon)

	err := r.Run(":80")
	if err != nil {
		log.Fatal(err)
	}
}
