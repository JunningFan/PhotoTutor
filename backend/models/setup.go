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
	conn = database.LogMode(true).Set("gorm:auto_preload", true)

	//register objects
	database.AutoMigrate(&User{})
	database.AutoMigrate(&Img{})
	database.AutoMigrate(&Location{})
	database.AutoMigrate(&Picture{})
}
