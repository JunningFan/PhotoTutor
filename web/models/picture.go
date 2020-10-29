package models

import (
	"fmt"
	"phototutor/backend/client"
	"time"

	"github.com/mmcloughlin/geohash"
	"gorm.io/gorm"
)

type Tag struct {
	Name string `gorm:"primaryKey";binding:"required"`
}

type Comment struct {
	gorm.Model
	PictureID uint
	Message   string `binding:"required"`
	UID       uint
}

type Vote struct {
	PictureID uint `gorm:"primaryKey"`
	UID       uint `gorm:"primaryKey"`
	Like      bool
}

type Picture struct {
	gorm.Model
	Title  string
	UserID uint
	NView  uint // how many views

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

	// for calculated filed
	NLike    int64  `gorm:"-"`
	NDislike int64  `gorm:"-"`
	Votes    []Vote `json:"-"`

	Tags []Tag `gorm:"many2many:picture_tag;"`
	// has many tags
	Comments []Comment `json:"-"`
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

func (p *Picture) AfterFind(tx *gorm.DB) (err error) {
	res := tx.Find(&Vote{PictureID: p.ID}).Where("like", true).Count(&p.NLike)
	if res.Error != nil {
		return res.Error
	}
	res = tx.Find(&Vote{PictureID: p.ID}).Where("like", false).Count(&p.NDislike)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

type PictureManager struct{}

func NewPictureManager() PictureManager {
	return PictureManager{}
}

func (p *PictureManager) All() ([]Picture, error) {
	var pictures []Picture
	res := conn.Joins("Location").Preload("Tags").Find(&pictures)
	// print(conn.Association("Tag"))
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
	go syncElsPicture(pic)
	return pic, res.Error
}

// One Find the one picture
func (p *PictureManager) One(pid uint) (Picture, error) {
	var picture Picture
	res := conn.Joins("Location").Preload("Tags").First(&picture, pid)
	go incPicNView(picture)
	return picture, res.Error
}

// Delete a photo that belong to the user
func (p *PictureManager) Delete(uid, pid uint) error {
	var picture Picture
	res := conn.First(&picture, pid)
	if res.Error != nil {
		return res.Error
	}
	if picture.UserID != uid {
		return fmt.Errorf("you can't delete the post not belongs to you")
	}
	res = conn.Delete(&picture)
	if res.Error != nil {
		return res.Error
	}
	// delete the entry of elastic and back to user
	go delElsPicture(picture.ID)
	return nil
}

// Comment Make a comment to a picture
func (p *PictureManager) Comment(pid uint, comment Comment) (Comment, error) {
	pic, err := p.One(pid)
	if err != nil {
		return Comment{}, err
	}

	comment.PictureID = pic.ID
	res := conn.Save(&comment)
	if res.Error != nil {
		return Comment{}, res.Error
	}
	go syncElsComment(comment)
	return comment, nil
}

// Like a picture post
func (p *PictureManager) Like(uid, pid uint) error {
	res := conn.Save(&Vote{PictureID: pid, UID: uid, Like: true})
	if res.Error != nil {
		go p.syncElsVote(pid)
	}
	return res.Error
}

// Dislike a picture post
func (p *PictureManager) Dislike(uid, pid uint) error {
	res := conn.Save(&Vote{PictureID: pid, UID: uid, Like: false})
	if res.Error != nil {
		go p.syncElsVote(pid)
	}
	return res.Error
}

/* Coroutine function, all expected the caller has wrap these function in a coroutine */
func (p *PictureManager) syncElsVote(pid uint) {
	pic, err := p.One(pid)
	if err != nil {
		fmt.Println("Sync Els Vote Error: ", err)
	}
	syncElsPicture(pic)
}

func delElsPicture(pid uint) {
	client.DelElsObj(fmt.Sprintf("picture/_doc/%d", pid))
}

func syncElsPicture(p Picture) {
	client.PutElsObj(fmt.Sprintf("picture/_doc/%d", p.ID), p)
}

func syncElsComment(c Comment) {
	client.PutElsObj(fmt.Sprintf("comment/_doc/%d", c.ID), c)
}

func incPicNView(p Picture) {
	p.NView++
	conn.Model(&p).Update("NView", p.NView)
	syncElsPicture(p)
}
