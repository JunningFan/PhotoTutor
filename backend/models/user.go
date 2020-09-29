package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique"`
	password  string
	Nickname  string
	Signature string
	img       uint
	Img       string
}

type UserRegisterInput struct {
	Username string
	Password string
	Nickname string
}

type UserChangePasswordInput struct {
	Password    string
	NewPassword string
}

type UserUpdateInput struct {
	Nickname  string
	Signature string
	img       uint
}

type UserManager struct{}

func NewUserManager() UserManager {
	return UserManager{}
}

func (um *UserManager) Create(input *UserRegisterInput) (User, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return User{}, err
	}
	user := User{Username: input.Username, password: string(encrypted), Nickname: input.Nickname}
	res := conn.Create(&user)
	if err := res.Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (um *UserManager) Update(user *User, input UserUpdateInput) {

}
