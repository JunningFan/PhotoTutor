package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"phototutor/backend/util"
)

var conn *gorm.DB

func Setup() {
	var err error
	if util.DB_DSN == "" {
		conn, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	} else {
		//	only support postgres connection
		conn, err = gorm.Open(postgres.Open(util.DB_DSN), &gorm.Config{})
	}
	if err != nil {
		panic(fmt.Sprintf("Fail to connect to database %v", err.Error()))
	}
	//conn = conn.LogMode(true).Set("gorm:auto_preload", true)

	//register objects
	err = conn.AutoMigrate(&User{}, &Img{}, &Location{}, &Picture{})
	if err != nil {
		panic(fmt.Sprintf("Fail to migrate database %v", err.Error()))
	}
	//println("Finish set up databse conn dsn ", util.DB_DSN)
}
