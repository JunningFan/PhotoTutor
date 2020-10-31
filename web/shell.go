package main

import (
	"fmt"
	"phototutor/backend/models"
)

func main() {
	models.Setup()
	// var pictures []models.Picture
	// models.Conn.Debug().Joins("Location").Preload("Tags").Find(&pictures)
	pm := models.NewPictureManager()
	pm.Like(5, 1)
	picture, _ := pm.One(1)
	fmt.Println(picture.NLike)
	fmt.Println(picture.NDislike)

	// for _, p := range pictures {
	// 	fmt.Println(p)
	// }
}
