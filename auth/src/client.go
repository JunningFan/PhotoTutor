package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

type NotificationInput struct {
	UID   uint
	Actor uint
	Type  string
}

// ErrorResp json type of remote error
type ErrorResp struct {
	Error string
}

var imgSer string
var notifySer string

//NewClient for remote img server
func NewClient(img  , notify string ) {
	if img == "" {
		img = "http://localhost:8083/"
	}
	if notify == "" {
		notify = "http://localhost:8084/"
	}
	imgSer = img
	notifySer = notify
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

// CreateNotification post a notification to remote
func CreateNotification(v interface{}) {
	
	fmt.Println(notifySer)
	if jbytes, err := json.Marshal(v); err != nil {
		log.Printf("Notification Marshal Err: %s\n", err.Error())
	} else if resp, err := http.Post(notifySer, "application/json", bytes.NewReader(jbytes)); err != nil {
		log.Printf("Notification Sync Err: %s\n", err.Error())
	} else if err := resp.Body.Close(); err != nil {
		log.Printf("Notification Close Fp Err: %s\n", err.Error())
	}
}