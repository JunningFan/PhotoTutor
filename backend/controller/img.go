package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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
func checkSuffix(suffix string) bool {
	suffix = strings.ToLower(suffix)
	// TODO: more on the types
	return suffix == "jpg" || suffix == "png"
}

func (c *ImgController) upload(uid uint, ctx *gin.Context) {
	fileReader, header, err := ctx.Request.FormFile("upload")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	originalName := header.Filename
	originNameArr := strings.Split(originalName, ".")
	imgId := models.AllocImgId(uid)
	suffix := originNameArr[len(originNameArr)-1]
	if len(originNameArr) < 2 || !checkSuffix(suffix) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("filetype is not suported for '%v'", originalName)})
		return
	}
	imgName := fmt.Sprintf("%d.%s", imgId, suffix)
	var out *os.File
	out, err = os.Create(path.Join(util.ImgBigPath, imgName))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("creat fileReader failed %s", err.Error())})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, fileReader)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("io filed %s", err.Error())})
		return
	}

	//	TODO: generate thumbnail import "github.com/nfnt/resize"
	ctx.JSON(http.StatusOK, gin.H{"img": imgId})

}
