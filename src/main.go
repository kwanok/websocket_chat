package main

import (
	"Friday/server"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()

	server.Routes(r)

	r.Run()
}
