package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"phototutor/backend/models"
	"phototutor/backend/util"
	"time"
)

type UserController struct {
	userManager models.UserManager
}

// for jwt token
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
func NewUserController(srvr *gin.RouterGroup) UserController {
	res := UserController{models.NewUserManager()}

	//srvr.GET("", RequrieAuth(res.getCurrUser))
	srvr.POST("", res.register)
	srvr.PUT("", RequrieAuth(res.update))
	srvr.POST("login/", res.login)
	srvr.GET("", RequrieAuth(res.getCurrUser))
	return res
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
			return []byte(util.SecretKey), nil
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

func (p *UserController) register(ctx *gin.Context) {
	var input models.UserRegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if user, err := p.userManager.Create(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if tokenStr, err := getJwtString(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"user": user, "token": tokenStr})
	}
}

func (p *UserController) update(uid uint, ctx *gin.Context) {
	var input models.UserUpdateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if user, err := p.userManager.Update(uid, input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, user)
	}
}

func getJwtString(user models.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, JwtClaims{
		ID:     user.ID,
		Expire: time.Now().Add(time.Minute * 20).Unix()})
	return token.SignedString([]byte(util.SecretKey))
}

func (uc *UserController) login(ctx *gin.Context) {
	var input models.UserLoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}
	//get user info
	user, err := uc.userManager.Login(&input)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tokenStr, err := getJwtString(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"user": user, "token": tokenStr})
	}
}

func (uc *UserController) getCurrUser(uid uint, ctx *gin.Context) {
	if user, err := uc.userManager.GetCurUser(uid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, user)
	}
}
