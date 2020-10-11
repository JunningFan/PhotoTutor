package util

import "os"

var SecretKey = "123hn312r9fxh28739dn182gdahs987tgd56afs"
var ImgStaticPrefix = "img/"
var ImgSmallPath = ImgStaticPrefix + "small/"
var ImgBigPath = ImgStaticPrefix + "big/"

var DB_DSN string

func SetUp() {
	DB_DSN = os.Getenv("DB_DSN")
}
