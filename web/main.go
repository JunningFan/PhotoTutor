package main

import (
	"net/http"
	"os"
	"phototutor/backend/controller"
	"phototutor/backend/models"
	"phototutor/backend/util"

	"github.com/gin-gonic/gin"
)

func main() {
	os.MkdirAll(util.ImgSmallPath, os.ModePerm)
	os.MkdirAll(util.ImgBigPath, os.ModePerm)

	server := gin.Default()
	// server.Use(TlsHandler())

	server.Static("img/", util.ImgStaticPrefix)
	models.Setup()
	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	picRoute := server.Group("/pictures/")
	{
		controller.NewPictureController(picRoute)
		controller.NewImgController(picRoute)
	}
	controller.NewUserController(server.Group("/users/"))
	err := server.Run()
	//running in tls
	// err := server.RunTLS(":8080", "ssl/rootCA.pem", "ssl/rootCA.key")
	if err != nil {
		println(err.Error())
	}
}

// func TlsHandler() gin.HandlerFunc {
// 	secureMiddleware := secure.New(secure.Options{
// 		SSLRedirect: true,
// 		SSLHost:     "127.0.0.1:8080",
// 	})
// 	return func(c *gin.Context) {
// 		err := secureMiddleware.Process(c.Writer, c.Request)

// 		// If there was an error, do not continue.
// 		if err != nil {
// 			return
// 		}

// 		c.Next()
// 	}
// }
