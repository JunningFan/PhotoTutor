package models

import (
	"fmt"
	"phototutor/backend/util"

	"github.com/mmcloughlin/geohash"
	"gorm.io/gorm"
)

type Picture struct {
	gorm.Model
	Title  string
	UserID uint
	// `json:"-"`
	// User   User
	//Uid uint
	Img        string `json:"-"`
	Lng        float64
	Lat        float64
	LocationID uint `json:"-"`
	Location   Location
	GeoHash    string // for elastic search

	Iso          uint
	FocalLength  uint
	Aperture     float64
	ShutterSpeed float64
	Timestamp    uint
	Orientation  float64
	Elevation    float64

	// fill while creating the picture obj
	Height uint
	Width  uint

	//Lid uint Fill by system
	ImgSmall string `gorm:"-"`
	ImgBig   string `gorm:"-"`
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
}

func (p *Picture) AfterFind(_ *gorm.DB) (err error) {
	p.ImgBig = fmt.Sprintf("%s%s", util.ImgBigPath, p.Img)
	p.ImgSmall = fmt.Sprintf("%s%s", util.ImgSmallPath, p.Img)
	return
}

type PictureManager struct{}

func NewPictureManager() PictureManager {
	return PictureManager{}
}

func (p *PictureManager) All() ([]Picture, error) {
	var pictures []Picture
	res := conn.Debug().Preload("Location").Find(&pictures)

	return pictures, res.Error
}

func (p *PictureManager) Insert(input *PictureInput) (Picture, error) {
	var img Img
	imgDb := conn.First(&img, input.Img)
	if imgDb.Error != nil {
		return Picture{}, fmt.Errorf("img %v is not exist", input.Img)
	}
	picName, err := img.GetImgFileName(input.Uid)
	if err != nil {
		return Picture{}, err
	}
	picHeight, picWidth, err := img.getResloution()
	if err != nil {
		return Picture{}, err
	}

	if err := GetLocation(&input.Location); err != nil {
		return Picture{}, err
	}
	pic := Picture{
		Title:    input.Title,
		UserID:   input.Uid,
		Img:      picName,
		Height:   picHeight,
		Width:    picWidth,
		Lng:      input.Lng,
		Lat:      input.Lat,
		Location: input.Location,
		GeoHash:  geohash.Encode(input.Lat, input.Lng),

		// specified keywords
		Iso:          input.Iso,
		FocalLength:  input.FocalLength,
		Aperture:     input.Aperture,
		ShutterSpeed: input.ShutterSpeed,
		Timestamp:    input.Timestamp,
		Orientation:  input.Orientation,
		Elevation:    input.Elevation,
	}
	res := conn.Create(&pic).Find(&pic)

	return pic, res.Error
}

// One Find the one picture
func (p *PictureManager) One(pid uint) (Picture, error) {
	var picture Picture
	res := conn.Debug().Preload("Location").First(&picture, pid)
	return picture, res.Error
}
