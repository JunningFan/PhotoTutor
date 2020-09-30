package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"phototutor/backend/util"
)

type Picture struct {
	gorm.Model
	Title    string
	User     User
	Img      uint `json:"-"`
	Lng      float64
	Lat      float64
	ImgSmall string `gorm:"-"`
	ImgBig   string `gorm:"-"`
}

type PictureInput struct {
	Title string  `binding:"required"`
	Uid   uint    // inject after login
	Lng   float64 `binding:"required"`
	Lat   float64 `binding:"required"`
	Img   uint    `binding:"required"`
}

func (p *Picture) AfterFind(_ *gorm.DB) (err error) {
	p.ImgBig = fmt.Sprintf("%s%v", util.ImgBigPath, p.Img)
	p.ImgSmall = fmt.Sprintf("%s%v", util.ImgSmallPath, p.Img)
	return
}

type PictureManager struct{}

func NewPictureManager() PictureManager {
	return PictureManager{}
}

func (p PictureManager) All() ([]Picture, error) {
	var pictures []Picture
	res := conn.Find(&pictures)
	return pictures, res.Error
}

func (p PictureManager) Insert(input *PictureInput) (Picture, error) {
	pic := Picture{Title: input.Title, Img: input.Img, Lng: input.Lng, Lat: input.Lat, User: User{ID: input.Uid}}
	res := conn.Create(&pic).Find(&pic)
	return pic, res.Error
}
