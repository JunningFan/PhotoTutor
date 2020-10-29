package main

/*
 * This file fill contain all the neccessary functionality bout auth
 * Author Tecty
 */

import (
	"os"

	"github.com/JunningFan/PhotoTutor/notification/src"
	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()
	src.NewNotificationController(server.Group("/"))
	src.Setup(os.Getenv("DB_DSN"))
	err := server.Run("0.0.0.0:8084")
	if err != nil {
		println(err.Error())
	}
}
