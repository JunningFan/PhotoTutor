package models

import (
	"fmt"
	"image"
	"os"
	"path"
	"phototutor/backend/util"
)

type Img struct {
	Id     uint
	Uid    uint
	Suffix string
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
	if reader, err:= os.Open(path.Join(util.ImgBigPath,i.picFileName())); err!= nil {
		return 0,0, err
	} else if im, _ , err:= image.DecodeConfig(reader); err != nil {
		return 0,0,err
	} else {
		return uint(im.Height), uint(im.Width), nil
	}
}