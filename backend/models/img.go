package models

import "fmt"

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
	return fmt.Sprintf("%d.%s", i.Id, i.Suffix), nil
}
