package main

/*
 * This file fill contain all the neccessary functionality bout auth
 * Author Tecty
 */

import (
	"os"

	"github.com/JunningFan/PhotoTutor/auth/src"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	src.NewUserController(server.Group("/"))
	src.Setup(os.Getenv("DB_DSN"))

	err := server.Run()
	if err != nil {
		println(err.Error())
	}
}
