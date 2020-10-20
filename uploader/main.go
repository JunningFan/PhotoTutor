package main

/*
 * This file fill contain all the neccessary functionality bout auth
 * Author Tecty
 */

import (
	"github.com/JunningFan/PhotoTutor/uploader/src"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	os.MkdirAll("img/small/", os.ModePerm)
	os.MkdirAll("img/big/", os.ModePerm)

	server := gin.Default()
	src.NewImgController(server.Group("/"))
	src.Setup(os.Getenv("DB_DSN"))
	err := server.Run("0.0.0.0:8083")
	if err != nil {
		println(err.Error())
	}
}
