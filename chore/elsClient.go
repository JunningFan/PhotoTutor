package main

import (
	"fmt"
	"log"
)

// Info of the image
type ImgInfo struct {
	height uint
	width  uint
	big    string
	small  string
}

func GetImgInfo(id uint) (ImgInfo, error) {
	img := ImgInfo{}
	if resp, err := http.get(fmt.Sprintf("%s"), "application/json"); err != nil {
		log.Printf("Els Sync Err: %s\n", err.Error())
	} else if err := resp.Body.Close(); err != nil {
		log.Printf("Els Close Fp Err: %s\n", err.Error())
	}
	return img, nil
}

func main() {
	img, err := GetImgInfo(1)
	if err != nil {
		fmd.Printf(err.Error())
	}
}
