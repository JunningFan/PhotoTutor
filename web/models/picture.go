package models

import (
	"fmt"
	"phototutor/backend/client"
	"time"

	"github.com/mmcloughlin/geohash"
	"gorm.io/gorm"
)

type Tag struct {
	Name string `gorm:"primaryKey"`
}

type Picture struct {
	gorm.Model
	Title  string
	UserID uint
	NView  uint // how many views

	// `json:"-"`
	// User   User
	//Uid uint
	// Img        string `json:"-"`
	Lng        float64
	Lat        float64
	LocationID uint `json:"-"`
	Location   Location
	GeoHash    string // for elastic search

	Iso          uint
	FocalLength  uint
	Aperture     float64
	ShutterSpeed float64
	Timestamp    time.Time
	Orientation  float64
	Elevation    float64

	// fill while creating the picture obj
	Height uint
	Width  uint

	// fill by creation
	ImgSmall string
	ImgBig   string

	Tags []Tag `gorm:"many2many:picture_tag;"`
}

type PictureInput struct {
	Title string `binding:"required"`
	Uid   uint
	// `json:"-"` // inject after login
	// User         User
	Lng          float64 `binding:"required"`
	Lat          float64 `binding:"required"`
	Iso          uint
	FocalLength  uint
	Aperture     float64
	ShutterSpeed float64
	Timestamp    uint
	Orientation  float64
	Elevation    float64

	Location Location
	Img      uint `binding:"required"`
	Tags     []string
}

type PictureManager struct{}

func NewPictureManager() PictureManager {
	return PictureManager{}
}

func (p *PictureManager) All() ([]Picture, error) {
	var pictures []Picture
	res := conn.Debug().Joins("Location").Preload("Tags").Find(&pictures)
	// print(conn.Debug().Association("Tag"))
	return pictures, res.Error
}

func (p *PictureManager) Insert(input *PictureInput) (Picture, error) {
	//RPC to img server to get img info
	imgInfo, err := client.GetImgInfo(input.Img)
	if err != nil {
		return Picture{}, err
	}
	if err := GetLocation(&input.Location); err != nil {
		return Picture{}, err
	}
	tags := make([]Tag, len(input.Tags))
	for i, tag_name := range input.Tags {
		tags[i] = Tag{Name: tag_name}
	}

	pic := Picture{
		Title:    input.Title,
		UserID:   input.Uid,
		NView:    0,
		ImgSmall: imgInfo.Small,
		ImgBig:   imgInfo.Big,
		Height:   imgInfo.Height,
		Width:    imgInfo.Width,
		Lng:      input.Lng,
		Lat:      input.Lat,
		Location: input.Location,
		GeoHash:  geohash.Encode(input.Lat, input.Lng),

		// specified keywords
		Iso:          input.Iso,
		FocalLength:  input.FocalLength,
		Aperture:     input.Aperture,
		ShutterSpeed: input.ShutterSpeed,
		Timestamp:    time.Unix(int64(input.Timestamp), 0),
		// Timestamp:   input.Timestamp,
		Orientation: input.Orientation,
		Elevation:   input.Elevation,
		Tags:        tags,
	}
	res := conn.Create(&pic).Find(&pic)
	go syncElsObj(pic)
	return pic, res.Error
}

// One Find the one picture
func (p *PictureManager) One(pid uint) (Picture, error) {
	var picture Picture
	res := conn.Debug().Joins("Location").Preload("Tags").First(&picture, pid)
	go incPicNView(picture)
	return picture, res.Error
}

/* Coroutine function, all expected the caller has wrap these function in a coroutine */

func syncElsObj(p Picture) {
	client.PutElsObj(fmt.Sprintf("picture/_doc/%d", p.ID), p)
}

func incPicNView(p Picture) {
	p.NView++
	conn.Model(&p).Update("NView", p.NView)
	syncElsObj(p)
}
