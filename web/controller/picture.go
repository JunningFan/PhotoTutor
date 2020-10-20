package controller

import (
	"fmt"
	"net/http"
	"phototutor/backend/elsClient"
	"phototutor/backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PictureController struct {
	pictureManager models.PictureManager
}

func NewPictureController(srvr *gin.RouterGroup) PictureController {
	res := PictureController{models.NewPictureManager()}

	srvr.GET("", res.getAll)
	srvr.POST("", RequrieAuth(res.insert))
	srvr.GET(":id", res.getOne)
	return res
}
func (p *PictureController) getOne(ctx *gin.Context) {
	id := ctx.Param("id")
	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The id of image must be string"})
	} else if pic, err := p.pictureManager.One(uint(idNum)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, pic)
	}

}
func (p PictureController) getAll(ctx *gin.Context) {
	if data, err := p.pictureManager.All(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": data})
	}
}

func (p PictureController) insert(uid uint, ctx *gin.Context) {
	// user, err := models.GetUserByID(uid)
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user not exist"})
	// 	return
	// }
	input := models.PictureInput{Uid: uid}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if pic, err := p.pictureManager.Insert(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		go elsClient.PutElsObj(fmt.Sprintf("picture/_doc/%d", pic.ID), pic)
		ctx.JSON(http.StatusOK, pic)
	}
}
