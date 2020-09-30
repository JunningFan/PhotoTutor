package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var conn *gorm.DB

func Setup() {
	database, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("Failed to connect to database!")
	}
	conn = database

	//register objects
	database.AutoMigrate(&Picture{})
	database.AutoMigrate(&User{})
	database.AutoMigrate(&Img{})
}
