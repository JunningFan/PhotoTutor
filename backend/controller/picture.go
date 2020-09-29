package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"phototutor/backend/models"
)

type PictureController struct {
	pictureManager models.PictureManager
}

func NewPictureController(srvr *gin.RouterGroup) PictureController {
	res := PictureController{models.NewPictureManager()}

	srvr.GET("", res.getAll)
	srvr.POST("", res.insert)
	return res
}

func (p PictureController) getAll(ctx *gin.Context) {
	if data, err := p.pictureManager.All(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": data})
	}
}

func (p PictureController) insert(ctx *gin.Context) {
	var input models.PictureInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if pic, err := p.pictureManager.Insert(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, pic)
	}
}
