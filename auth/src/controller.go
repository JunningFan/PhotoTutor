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
	srvr.POST("nicknames/", res.getNicknames)
	srvr.GET("detail/:id", res.getOne)
	srvr.GET("", RequrieAuth(res.getCurrUser))
	srvr.POST("follow/:id", RequrieAuth(res.follow))
	srvr.DELETE("follow/:id", RequrieAuth(res.unfollow))
	srvr.GET("follow/", RequrieAuth(res.getMyFollowing))
	srvr.GET("follow/ami/:id", RequrieAuth(res.amIFollowing))
	srvr.GET("follow/ing/:id", res.getFollowing)
	srvr.GET("follow/ers/:id", res.getFollowers)
	return res
}

/* Redis helper */
func setUserLoginRedis(uid uint, access, refresh string) {
	err := rdb.HSet(redisCtx,
		fmt.Sprintf("%v", uid),
		map[string]interface{}{"access": access, "refresh": refresh}).Err()
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func checkRefreshToken(uid uint, refresh string) (string, error) {
	val, err := rdb.HGet(redisCtx, fmt.Sprintf("%v", uid), "refresh").Result()
	if err != nil {
		return "", err
	}
	if val != refresh {
		return "", fmt.Errorf("invalid refresh token")
	}
	return getAccessToken(uid)
}

// updateAccessToken coroutine to update the access token
func updateAccessToken(uid uint, access string) {
	err := rdb.HSet(redisCtx,
		fmt.Sprintf("%v", uid), "access", access).Err()
	if err != nil {
		fmt.Printf("redis update access: %v", err.Error())
	}
}

func checkAccessToken(uid uint, access string) bool {
	val, err := rdb.HGet(redisCtx, fmt.Sprintf("%v", uid), "access").Result()
	//fmt.Printf("val: %s\nAccess:%s\n", val, access)
	return err == nil && val == access
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
		} else if claims, ok := token.Claims.(*jwtClaimAccess); !ok || !checkAccessToken(claims.ID, tokenStr) {
			//authorized and pass the context
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad Token"})
		} else {
			handler(claims.ID, context)
		}
	}
}

func (uc *UserController) getOne(ctx *gin.Context) {
	id := ctx.Param("id")
	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "the id of user must be string"})
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
		go setUserLoginRedis(user.ID, accessToken, refreshToken)
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
		Expire: time.Now().Add(time.Hour * 24).Unix()})
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
		go setUserLoginRedis(user.ID, accessToken, refreshToken)
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
	} else if claims, ok := token.Claims.(*jwtClaimRefresh); !ok {
		//authorized and pass the ctx
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad Token"})
	} else if access, err := checkRefreshToken(claims.ID, input.Refresh); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		go updateAccessToken(claims.ID, access)
		ctx.JSON(http.StatusOK, gin.H{"access": access})
	}
}

func (uc *UserController) getCurrUser(uid uint, ctx *gin.Context) {
	if user, err := uc.userManager.GetUser(uid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, user)
	}
}

type idsInput struct {
	IDs []uint `binding:"required"`
}

func (uc *UserController) getNicknames(ctx *gin.Context) {
	var input idsInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if ret, err := uc.userManager.ResolveNicknameByIds(input.IDs); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": ret})
	}
}

func (uc *UserController) follow(uid uint, ctx *gin.Context) {
	id := ctx.Param("id")
	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "the id of user must be string"})
	} else if err := uc.userManager.Follow(uid, uint(idNum)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": "Followed"})
	}
}

func (uc *UserController) unfollow(uid uint, ctx *gin.Context) {
	id := ctx.Param("id")
	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "the id of user must be string"})
	} else if err := uc.userManager.Unfollow(uid, uint(idNum)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": "Unfollowed"})
	}
}

func (uc *UserController) getFollowing(ctx *gin.Context) {
	id := ctx.Param("id")
	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "the id of user must be string"})
	} else {
		uc.getMyFollowing(uint(idNum), ctx)
	}
}

func (uc *UserController) getFollowers(ctx *gin.Context) {
	id := ctx.Param("id")
	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "the id of user must be string"})
	} else if users, err := uc.userManager.FollowerList(uint(idNum)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": users})
	}
}

func (uc *UserController) getMyFollowing(uid uint, ctx *gin.Context) {
	if users, err := uc.userManager.FollowingList(uid); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": users})
	}
}

func (uc *UserController) amIFollowing(uid uint, ctx *gin.Context) {
	id := ctx.Param("id")
	if idNum, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "the id of user must be string"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": uc.userManager.IsFolloing(uid, uint(idNum))})
	}
}
