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

	srvr.GET("/", res.getAll)
	srvr.POST("", RequrieAuth(res.insert))
	return res
}

func (p PictureController) getAll(ctx *gin.Context) {
	if data, err := p.pictureManager.All(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": data})
	}
}

func (p PictureController) insert(uid uint, ctx *gin.Context) {
	user, err := models.GetUserByID(uid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user not exist"})
		return
	}
	input := models.PictureInput{Uid: uid, User: user}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if pic, err := p.pictureManager.Insert(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, pic)
	}
}
