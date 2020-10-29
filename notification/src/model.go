package src

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Notification of user
type Notification struct {
	ID      uint `gorm:"primaryKey"`
	UID     uint
	Type    string
	Message string
}

var (
	conn *gorm.DB
)

// Setup the database connection
func Setup(dbDsn string) {
	var err error
	if dbDsn == "" {
		conn, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	} else {
		fmt.Println(dbDsn)
		//	only support postgres connection
		conn, err = gorm.Open(postgres.Open(dbDsn), &gorm.Config{})
	}
	if err != nil {
		panic(fmt.Sprintf("Fail to connect to database %v", err.Error()))
	}
	//conn = conn.LogMode(true).Set("gorm:auto_preload", true)

	//register objects
	err = conn.AutoMigrate(&Notification{})
	if err != nil {
		panic(fmt.Sprintf("Fail to migrate database %v", err.Error()))
	}
	//println("Finish set up databse conn dsn ", dbDsn)
}

// CreateMsg creat a new message in databse
func CreateMsg(notification Notification) (Notification, error) {
	res := conn.Create(&notification)
	return notification, res.Error
}

// GetMsgList get the message list for a user by uid
func GetMsgList(uid uint) ([]Notification, error) {
	var ret []Notification
	res := conn.Where("uid", uid).Find(&ret)
	return ret, res.Error
}

// RemoveMsgList from lastID Before for that user, the record will be deleted
func RemoveMsgList(uid, lastID uint) error {
	res := conn.Where("id <= ? AND uid = ?", lastID, uid).Delete(Notification{})
	return res.Error
}
