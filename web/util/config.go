package util

import "os"

var SecretKey = "123hn312r9fxh28739dn182gdahs987tgd56afs"
var ImgStaticPrefix = "img/"
var ImgSmallPath = ImgStaticPrefix + "small/"
var ImgBigPath = ImgStaticPrefix + "big/"

var DB_DSN string
var ELS_BASE string

func SetUp() {
	DB_DSN = os.Getenv("DB_DSN")
	ELS_BASE = os.Getenv("ELS_BASE")
	if ELS_BASE == "" {
		ELS_BASE = "http://localhost:9200/"
	}
}
