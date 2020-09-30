package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Username  string     `gorm:"unique"`
	Password  string     `json:"-"`
	Nickname  string
	Signature string
	Img       uint   `json:"-"`
	ImgLoc    string `gorm:"-" json:"img"`
}

type UserRegisterInput struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
	Nickname string `binding:"required"`
}

type UserLoginInput struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

type UserChangePasswordInput struct {
	Password    string `binding:"required"`
	NewPassword string `binding:"required"`
}

type UserUpdateInput struct {
	Nickname  string `binding:"required"`
	Signature string `binding:"required"`
	Img       uint   `binding:"required"`
}

type UserManager struct{}

func NewUserManager() UserManager {
	return UserManager{}
}

func (u *User) SetPassword(password string) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(encrypted)
	return nil
}

func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("wrong password")
	}
	return nil
}

func (um *UserManager) Create(input *UserRegisterInput) (User, error) {
	user := User{Username: input.Username, Nickname: input.Nickname}
	if err := user.SetPassword(input.Password); err != nil {
		return User{}, err
	}
	res := conn.Create(&user)
	if err := res.Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (um *UserManager) Login(input *UserLoginInput) (User, error) {
	var user User
	res := conn.Where("Username = ?", input.Username).Find(&user)
	if err := res.Error; err != nil {
		return User{}, fmt.Errorf("User doesn't exist!")
	}
	if err := user.CheckPassword(input.Password); err != nil {
		return User{}, err
	}
	return user, nil
}

func (um *UserManager) GetCurUser(uid uint) (User, error) {
	var ret User
	res := conn.First(&ret, uid)
	return ret, res.Error
}

//func (um *UserManager) Update(user *User, input UserUpdateInput) {
//
//}
