package models

import "github.com/jinzhu/gorm"

type Picture struct {
	gorm.Model
	Title    string
	Uid      uint
	img      uint
	ImgSmall string `gorm:"-"`
	ImgBig   string `gorm:"-"`
}

type PictureInput struct {
	Title string `json:"title"`
	Uid   uint   `json:"uid"`
	Img   uint   `json:"img"`
}

func (p *Picture) AfterFind(_ *gorm.DB) (err error) {
	p.ImgBig = "/Some/prefix/" + string(p.img)
	p.ImgSmall = "/Some/prefix/small/" + string(p.img)
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
	pic := Picture{Title: input.Title, Uid: input.Uid}
	res := conn.Create(&pic).Find(&input)
	return pic, res.Error
}
