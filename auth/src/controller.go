package src

import (
	"fmt"
	"net/http"
	"strconv"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const secretKey = "123hn312r9fxh28739dn182gdahs987tgd56afs"

// UserController for holding controller neccessary functions
type UserController struct {
	userManager UserManager
}

type jwtClaims struct {
	ID     uint
	Access bool
	Expire int64
}

type jwtClaimAccess jwtClaims
type jwtClaimRefresh jwtClaims

func (c jwtClaimAccess) Valid() error {
	if c.Access != true {
		return fmt.Errorf("not an access token")
	}
	if time.Now().Unix() > c.Expire || c.Expire == 0 {
		return fmt.Errorf("auth token is expred")
	}
	return nil
}
func (c jwtClaimRefresh) Valid() error {
	if c.Access == true {
		return fmt.Errorf("not a refresh token")
	}
	if time.Now().Unix() > c.Expire || c.Expire == 0 {
		return fmt.Errorf("auth token is expred")
	}
	return nil
}

// NewUserController for creating user controller
func NewUserController(srvr *gin.RouterGroup) UserController {
	res := UserController{NewUserManager()}

	//srvr.GET("", RequrieAuth(res.getCurrUser))
	srvr.POST("", res.register)
	srvr.PUT("", RequrieAuth(res.update))
	srvr.POST("login/", res.login)
	srvr.POST("refresh/", res.refresh)
	srvr.GET(":id", res.getOne)
	srvr.GET("", RequrieAuth(res.getCurrUser))
	return res
}

// RequrieAuth Passing a handler that User is the first variable
// Then is auth is success, the handler will be called
func RequrieAuth(handler func(uint, *gin.Context)) gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenStr := context.GetHeader("Authorization")
		if len(tokenStr) == 0 {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "login required"})
			return
		}
		token, err := jwt.ParseWithClaims(tokenStr, &jwtClaimAccess{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if !token.Valid {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "expired Token"})
		} else if claims, ok := token.Claims.(*jwtClaimAccess); token.Valid && ok {
			handler(claims.ID, context)
		} else {
			//authorized and pass the context
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad Token"})
		}
	}
}

func (uc *UserController) getOne(ctx *gin.Context) {
	id := ctx.Param("id")
	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The id of image must be string"})
	} else if user, err := uc.userManager.GetUser(uint(idNum)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, user)
	}

}

func (uc *UserController) register(ctx *gin.Context) {
	var input UserRegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if user, err := uc.userManager.Create(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if accessToken, refreshToken, err := getJwtString(user.ID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"user": user, "access": accessToken, "refresh": refreshToken})
	}
}

func (uc *UserController) update(uid uint, ctx *gin.Context) {
	var input UserUpdateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if user, err := uc.userManager.Update(uid, input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, user)
	}
}

func getAccessToken(user uint) (string, error) {
	short := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtClaimAccess{
		ID:     user,
		Access: true,
		Expire: time.Now().Add(time.Minute * 20).Unix()})
	return short.SignedString([]byte(secretKey))
}

func getJwtString(user uint) (string, string, error) {
	long := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtClaimRefresh{
		ID:     user,
		Access: false,
		Expire: time.Now().Add(time.Hour * 24 * 180).Unix()})

	if shortBarr, err := getAccessToken(user); err != nil {
		return "", "", err
	} else if longBarr, err := long.SignedString([]byte(secretKey)); err != nil {
		return "", "", err
	} else {
		return shortBarr, longBarr, nil
	}
}

func (uc *UserController) login(ctx *gin.Context) {
	var input UserLoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//get user info
	user, err := uc.userManager.Login(&input)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if accessToken, refreshToken, err := getJwtString(user.ID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"user": user, "access": accessToken, "refresh": refreshToken})
	}
}

type refreshInput struct {
	Refresh string `binding:"required"`
}

func (uc *UserController) refresh(ctx *gin.Context) {
	var input refreshInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.ParseWithClaims(input.Refresh, &jwtClaimRefresh{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if !token.Valid {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "expired Token"})
	} else if claims, ok := token.Claims.(*jwtClaimRefresh); token.Valid && ok {
		if access, err := getAccessToken(claims.ID); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"access": access})
		}
	} else {
		//authorized and pass the ctx
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad Token"})
	}
}

func (uc *UserController) getCurrUser(uid uint, ctx *gin.Context) {
	if user, err := uc.userManager.GetUser(uid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, user)
	}
}
