package main

import (
	"friday/server"
	"friday/tools"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	server.InitDB()
	server.Routes(r)

	err := r.Run()
	tools.ErrorHandler(err)
}
