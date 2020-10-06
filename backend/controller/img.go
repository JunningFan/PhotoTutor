package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"phototutor/backend/models"
	"phototutor/backend/util"
	"strings"
)

type ImgController struct{}

func NewImgController(srvr *gin.RouterGroup) ImgController {
	res := ImgController{}

	srvr.POST("upload/", RequrieAuth(res.upload))
	return res
}
func checkSuffix(suffix *string) bool {
	*suffix = strings.ToLower(*suffix)
	// TODO: more on the types
	return *suffix == "jpg" || *suffix == "png"
}

func (i *ImgController) mkThumbnail(path string, imgType string) {
	img, err := getBigImgByName(path, imgType)
	if err != nil {
		return
	}
	thImg := resize.Thumbnail(512, 512, img, resize.Bilinear)
	putSmallImgByImg(thImg, path, imgType)
}

func putSmallImgByImg(img image.Image, imgName string, imgType string) {
	out, err := os.Create(path.Join(util.ImgSmallPath, imgName))
	if err != nil {
		log.Println(err.Error())
	}
	defer out.Close()

	// write new image to file
	if imgType == "jpg" || imgType == "jpeg" {
		err = jpeg.Encode(out, img, nil)
	} else if imgType == "png" {
		err = png.Encode(out, img)
	}
	if err != nil {
		log.Println(err.Error())
	}
}

func getBigImgByName(imgName string, imgType string) (image.Image, error) {
	var empty image.Image
	file, err := os.Open(path.Join(util.ImgBigPath, imgName))
	if err != nil {
		return empty, err
	}
	if imgType == "jpg" {
		if ret, err := jpeg.Decode(file); err != nil {
			return empty, err
		} else {
			return ret, nil
		}

	} else if imgType == "png" {
		if ret, err := png.Decode(file); err != nil {
			return empty, err
		} else {
			return ret, err
		}
	}
	return empty, nil
}

func (c *ImgController) upload(uid uint, ctx *gin.Context) {
	fileReader, header, err := ctx.Request.FormFile("upload")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	originalName := header.Filename
	originNameArr := strings.Split(originalName, ".")

	suffix := originNameArr[len(originNameArr)-1]
	if len(originNameArr) < 2 || !checkSuffix(&suffix) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("filetype is not suported for '%v'", originalName)})
		return
	}
	var imgId uint
	var out *os.File

	imgId, err = models.AllocImgId(uid, suffix)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	imgName := fmt.Sprintf("%d.%s", imgId, suffix)
	out, err = os.Create(path.Join(util.ImgBigPath, imgName))
	defer out.Close()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("creat fileReader failed %s", err.Error())})
		return
	}
	_, err = io.Copy(out, fileReader)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("io filed %s", err.Error())})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"img": imgId})
	// generate thumbnail import "github.com/nfnt/resize"
	go c.mkThumbnail(imgName, suffix)
}
