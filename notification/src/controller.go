package src

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type NotificationController struct{}

const secretKey = "123hn312r9fxh28739dn182gdahs987tgd56afs"

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
			return []byte(secretKey), nil
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

func NewNotificationController(srvr *gin.RouterGroup) NotificationController {
	ret := NotificationController{}
	srvr.GET("/", RequrieAuth(ret.getNotifications))
	srvr.POST("/", ret.createNewNotification)
	srvr.DELETE("/:id", RequrieAuth(ret.deletNotifications))
	return ret
}

type NotificationInput struct {
	UID   uint
	Actor uint
	Type  string
}

func (ic *NotificationController) createNewNotification(ctx *gin.Context) {
	input := NotificationInput{}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actor, err := GetUserInfo(input.Actor)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := Notification{UID: input.UID, Type: input.Type}
	switch input.Type {
	case "follow":
		notification.Message = fmt.Sprintf("%s starts following on you!", actor.Nickname)
	case "comment":
		notification.Message = fmt.Sprintf("%s comment on your post!", actor.Nickname)
	}

	if notification, err := CreateMsg(notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, notification)
	}
}

func (ic *NotificationController) getNotifications(uid uint, ctx *gin.Context) {
	if notList, err := GetMsgList(uid); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": notList})
	}
}

func (ic *NotificationController) deletNotifications(uid uint, ctx *gin.Context) {
	id := ctx.Param("id")
	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The id of image must be string"})
	} else if err := RemoveMsgList(uid, uint(idNum)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": "you have succcessfully removed notifications"})
	}
}
