package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"phototutor/backend/controller"
	"phototutor/backend/models"
)

func main() {
	r := gin.Default()
	models.Setup()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	controller.NewPictureController(r.Group("/pictures/"))
	err := r.Run()
	println(err.Error())
}
