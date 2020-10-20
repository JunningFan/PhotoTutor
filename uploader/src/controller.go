package src

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
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"

	"strings"
)

type ImgController struct{}
const SecretKey = "123hn312r9fxh28739dn182gdahs987tgd56afs"
const ImgStaticPrefix = "img/"
const ImgSmallPath = ImgStaticPrefix + "small/"
const ImgBigPath = ImgStaticPrefix + "big/"



type JwtClaims struct {
	ID     uint
	Expire int64
}

func (c JwtClaims) Valid() error {
	if time.Now().Unix() > c.Expire || c.Expire == 0 {
		return fmt.Errorf("auth token is expred")
	}
	return nil
}

// Passing a handler that User is the first variable
// Then is auth is success, the handler will be called
func RequrieAuth(handler func(uint, *gin.Context)) gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenStr := context.GetHeader("Authorization")
		if len(tokenStr) == 0 {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "login required"})
			return
		}
		token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !token.Valid {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "expired Token"})
			return
		}
		if claims, ok := token.Claims.(*JwtClaims); token.Valid && ok {
			handler(claims.ID, context)

		} else {
			//authorized and pass the context
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad Token"})
		}
	}
}



func NewImgController(srvr *gin.RouterGroup) ImgController {
	res := ImgController{}
	srvr.GET("/:id",res.getPath)
	srvr.POST("/", RequrieAuth(res.upload))
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
	out, err := os.Create(path.Join(ImgSmallPath, imgName))
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
	file, err := os.Open(path.Join(ImgBigPath, imgName))
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

	imgId, err = AllocImgId(uid, suffix)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	imgName := fmt.Sprintf("%d.%s", imgId, suffix)
	out, err = os.Create(path.Join(ImgBigPath, imgName))
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

func (c *ImgController) getPath(ctx *gin.Context) {
	if idNum, err := strconv.ParseUint(ctx.Param("id"), 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "the id of image must be string"})
	} else if i, err := First(uint(idNum));  err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if height, width, err:= i.getResloution(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"uid": i.Uid,
			"height": height,
			"width": width,
			"big" : fmt.Sprintf("%s%d", ImgBigPath, i.Id),
			"small":  fmt.Sprintf("%s%d", ImgSmallPath, i.Id),
		})
	}
}
