package controller

import (
	"net/http"
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
	srvr.DELETE(":id", RequrieAuth(res.delete))
	srvr.POST(":id/comment", RequrieAuth(res.comment))
	srvr.POST(":id/like", RequrieAuth(res.like))
	srvr.POST(":id/dislike", RequrieAuth(res.dislike))

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
func (p *PictureController) getAll(ctx *gin.Context) {
	if data, err := p.pictureManager.All(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": data})
	}
}

func (p *PictureController) insert(uid uint, ctx *gin.Context) {
	input := models.PictureInput{Uid: uid}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if pic, err := p.pictureManager.Insert(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, pic)
	}
}

func (p *PictureController) comment(uid uint, ctx *gin.Context) {
	id := ctx.Param("id")
	input := models.Comment{UID: uid}

	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The id of image must be string"})
	} else if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if comment, err := p.pictureManager.Comment(uint(idNum), input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, comment)
	}

}

func (p *PictureController) like(uid uint, ctx *gin.Context) {
	id := ctx.Param("id")

	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The id of image must be string"})
	} else if err := p.pictureManager.Like(uid, uint(idNum)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": "You have successfully like the picture"})
	}
}

func (p *PictureController) dislike(uid uint, ctx *gin.Context) {
	id := ctx.Param("id")

	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The id of image must be string"})
	} else if err := p.pictureManager.Dislike(uid, uint(idNum)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": "You have successfully like the picture"})
	}
}

func (p *PictureController) delete(uid uint, ctx *gin.Context) {
	id := ctx.Param("id")
	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The id of image must be string"})
	} else if err := p.pictureManager.Delete(uid, uint(idNum)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": "You have successfully delete the picture"})
	}
}
