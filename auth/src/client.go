package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"bytes"
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
var notifSer string

//NewClient for remote img server
func NewClient(server string) {
	if server == "" {
		imgSer = "http://localhost:8083/"
	} else {
		imgSer = server
	}
}

func NotifServ(server string) {
	if server == "" {
		notifSer = "http://localhost:8084/"
	} else {
		notifSer = server
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

// CreateNotification post a notification to remote
func CreateNotification(v interface{}) {
	
	fmt.Println(notifSer)
	if jbytes, err := json.Marshal(v); err != nil {
		log.Printf("Notification Marshal Err: %s\n", err.Error())
	} else if resp, err := http.Post(notifSer, "application/json", bytes.NewReader(jbytes)); err != nil {
		log.Printf("Notification Sync Err: %s\n", err.Error())
	} else if err := resp.Body.Close(); err != nil {
		log.Printf("Notification Close Fp Err: %s\n", err.Error())
	}
}