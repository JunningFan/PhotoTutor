package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"phototutor/backend/util"
)

// PutElsObj create a remote elastic object
func PutElsObj(url string, v interface{}) {
	if jbytes, err := json.Marshal(v); err != nil {
		log.Printf("Els Marshal Err: %s\n", err.Error())
	} else if resp, err := http.Post(util.ELS_BASE+url, "application/json", bytes.NewReader(jbytes)); err != nil {
		log.Printf("Els Sync Err: %s\n", err.Error())
	} else if err := resp.Body.Close(); err != nil {
		log.Printf("Els Close Fp Err: %s\n", err.Error())
	}
}

// DelElsObj Delete a remote elastic object
func DelElsObj(url string) {
	if req, err := http.NewRequest("DELETE", util.ELS_BASE+url, nil); err != nil {
		log.Printf("Els Create Request Err: %s\n", err.Error())
	} else if resp, err := http.DefaultClient.Do(req); err != nil {
		log.Printf("Els Sync Err: %s\n", err.Error())
	} else if err := resp.Body.Close(); err != nil {
		log.Printf("Els Close Fp Err: %s\n", err.Error())
	}

}

// NotificationInput of notification
type NotificationInput struct {
	UID   uint
	Actor uint
	Type  string
}

// CreateNotification post a notification to remote
func CreateNotification(v interface{}) {
	if jbytes, err := json.Marshal(v); err != nil {
		log.Printf("Notification Marshal Err: %s\n", err.Error())
	} else if resp, err := http.Post(util.NOTIF_SER, "application/json", bytes.NewReader(jbytes)); err != nil {
		log.Printf("Notification Sync Err: %s\n", err.Error())
	} else if resp.StatusCode != http.StatusOK {
		log.Printf("Notification Server Err: %s \n", resp.Body)
	} else if err := resp.Body.Close(); err != nil {
		log.Printf("Notification Close Fp Err: %s\n", err.Error())
	}
}

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

func getImgInfo(id, uid uint) (ImgInfo, error) {
	img := ImgInfo{}
	errResp := ErrorResp{}
	resp, err := http.Get(fmt.Sprintf("%s%d", util.IMG_SER, id))
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
