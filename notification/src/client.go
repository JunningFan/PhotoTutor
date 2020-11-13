package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type UserInfo struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Username  string     `gorm:"unique"`
	Password  string     `json:"-"`
	Nickname  string
	Signature string
	ImgLoc    string `gorm:"-" json:"img"`
}

type ErrorResp struct {
	Error string
}

var authSerAddr string

func NewClient(authSer string) {
	if authSer == "" {
		authSer = "http://localhost:8080/"
	}
	authSerAddr = authSer + "detail/"
}

func getUserInfo(id uint) (UserInfo, error) {
	img := UserInfo{}
	errResp := ErrorResp{}
	resp, err := http.Get(fmt.Sprintf("%s%d", authSerAddr, id))
	if err != nil {
		// fmt.Println("===>>>> 1")
		return UserInfo{}, err
	} else if resp.StatusCode != 200 {
		if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
			// fmt.Println("===>>>> 2")
			return UserInfo{}, err
		} else if err := json.Unmarshal(bodyBytes, &errResp); err != nil {
			// fmt.Println("===>>>> 3")
			return UserInfo{}, err
		} else {
			// fmt.Println("===>>>> 4")
			return UserInfo{}, fmt.Errorf(errResp.Error)
		}
	} else if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
		return UserInfo{}, err
	} else if err := json.Unmarshal(bodyBytes, &img); err != nil {
		return UserInfo{}, err
	} else if err := resp.Body.Close(); err != nil {
		return UserInfo{}, err
	} else {
		return img, nil
	}
}

// GetImgInfo from remote server
func GetUserInfo(id uint) (UserInfo, error) {
	ret, err := getUserInfo(id)
	if err != nil {
		err = fmt.Errorf("auth: %s", err.Error())
	}
	return ret, err
}
