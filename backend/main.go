package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"phototutor/backend/controller"
	"phototutor/backend/models"
	"phototutor/backend/util"
)

func main() {
	os.MkdirAll(util.ImgSmallPath, os.ModePerm)
	os.MkdirAll(util.ImgBigPath, os.ModePerm)

	r := gin.Default()
	models.Setup()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	picRoute := r.Group("/pictures/")
	{
		controller.NewPictureController(picRoute)
		controller.NewImgController(picRoute)
	}
	controller.NewUserController(r.Group("/users/"))
	err := r.Run()
	if err != nil {
		println(err.Error())
	}
}
