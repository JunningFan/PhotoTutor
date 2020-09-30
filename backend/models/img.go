package models

type Img struct {
	Id  uint
	Uid uint
}

func AllocImgId(uid uint) uint {
	img := Img{Uid: uid}
	conn.Create(&img)
	return img.Id
}
