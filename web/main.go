package main

import (
	"phototutor/backend/controller"
	"phototutor/backend/models"
	"phototutor/backend/util"

	"github.com/gin-gonic/gin"
)

func setUpEnv() {

	util.SetUp()
	models.Setup()
}

func main() {
	setUpEnv()

	server := gin.Default()
	picRoute := server.Group("/")
	{
		controller.NewPictureController(picRoute)
	}
	// controller.NewUserController(server.Group("/users/"))
	err := server.Run("0.0.0.0:8081")
	if err != nil {
		println(err.Error())
	}
}
