package main

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
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


	r.Use(TlsHandler())

	r.Static("img/", util.ImgStaticPrefix)
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
	//err := r.Run()
	//running in tls
	err:= r.RunTLS(":8080", "ssl/rootCA.pem", "ssl/rootCA.key")
	if err != nil {
		println(err.Error())
	}
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "127.0.0.1:8080",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}