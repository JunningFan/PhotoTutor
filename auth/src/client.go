package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ImgInfo of the image
type ImgInfo struct {
	UID    uint
	Height uint
	Width  uint
	Big    string
	Small  string
}

// ErrorResp json type of remote error
type ErrorResp struct {
	Error string
}

var imgSer string

//NewClient for remote img server
func NewClient(server string) {
	if server == "" {
		imgSer = "http://localhost:8083/"
	} else {
		imgSer = server
	}
}

func getImgInfo(id, uid uint) (ImgInfo, error) {
	img := ImgInfo{}
	errResp := ErrorResp{}
	resp, err := http.Get(fmt.Sprintf("%s%d", imgSer, id))
	if err != nil {
		return ImgInfo{}, err
	} else if resp.StatusCode != 200 {
		if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
			return ImgInfo{}, err
		} else if err := json.Unmarshal(bodyBytes, &errResp); err != nil {
			return ImgInfo{}, err
		} else {
			return ImgInfo{}, fmt.Errorf(errResp.Error)
		}
	} else if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
		return ImgInfo{}, err
	} else if err := json.Unmarshal(bodyBytes, &img); err != nil {
		return ImgInfo{}, err
	} else if err := resp.Body.Close(); err != nil {
		return ImgInfo{}, err
	} else if img.UID != uid {
		return ImgInfo{}, fmt.Errorf("could not access this photo")
	} else {
		return img, nil
	}
}

// GetImgInfo from remote server
func GetImgInfo(id, uid uint) (ImgInfo, error) {
	ImgInfo, err := getImgInfo(id, uid)
	if err != nil {
		err = fmt.Errorf("uploader: %v", err.Error())
	}
	return ImgInfo, err
}
