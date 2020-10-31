package src

import (
	"fmt"
	"image"
	"os"
	"path"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Img struct {
	Id     uint
	Uid    uint
	Suffix string
}

var (
	conn *gorm.DB
)

func Setup(DB_DSN string) {
	var err error
	if DB_DSN == "" {
		conn, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	} else {
		fmt.Println(DB_DSN)
		//	only support postgres connection
		conn, err = gorm.Open(postgres.Open(DB_DSN), &gorm.Config{})
	}
	if err != nil {
		panic(fmt.Sprintf("Fail to connect to database %v", err.Error()))
	}
	//conn = conn.LogMode(true).Set("gorm:auto_preload", true)

	//register objects
	err = conn.AutoMigrate(&Img{})
	if err != nil {
		panic(fmt.Sprintf("Fail to migrate database %v", err.Error()))
	}
	//println("Finish set up databse conn dsn ", DB_DSN)
}

func AllocImgId(uid uint, suffix string) (uint, error) {
	img := Img{Uid: uid, Suffix: suffix}
	res := conn.Create(&img)
	return img.Id, res.Error
}

func (i *Img) GetImgFileName(uid uint) (string, error) {
	if uid != i.Uid {
		return "", fmt.Errorf("not premit to use this picture")
	}
	return i.picFileName(), nil
}

func (i *Img) picFileName() string {
	return fmt.Sprintf("%d.%s", i.Id, i.Suffix)
}

func (i *Img) getResloution() (uint, uint, error) {
	if reader, err := os.Open(path.Join(ImgBigPath, i.picFileName())); err != nil {
		return 0, 0, err
	} else if im, _, err := image.DecodeConfig(reader); err != nil {
		return 0, 0, err
	} else {
		return uint(im.Height), uint(im.Width), nil
	}
}

func First(pid uint) (Img, error) {
	img := Img{}
	res := conn.First(&img, pid)
	return img, res.Error
}
