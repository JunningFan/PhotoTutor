package util

import "os"

const SecretKey = "123hn312r9fxh28739dn182gdahs987tgd56afs"
const ImgStaticPrefix = "img/"
const ImgSmallPath = ImgStaticPrefix + "small/"
const ImgBigPath = ImgStaticPrefix + "big/"

var (
	DB_DSN    string
	ELS_BASE  string
	NOTIF_SER string
	IMG_SER   string
)

func SetUp() {
	DB_DSN = os.Getenv("DB_DSN")
	ELS_BASE = os.Getenv("ELS_BASE")
	if ELS_BASE == "" {
		ELS_BASE = "http://localhost:9200/"
	}
	IMG_SER = os.Getenv("IMG_SER")
	if IMG_SER == "" {
		IMG_SER = "http://localhost:8083/"
	}

	NOTIF_SER = os.Getenv("NOTIF_SER")
	if IMG_SER == "" {
		IMG_SER = "http://localhost:8084/"
	}
}
